package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

// command is used to transport message to and from the debugger.
type command struct {
	Action     string
	Parameters map[string]string
}

// addParam adds a key,value string pair to the command ; no check on overwrites.
func (self *command) addParam(key, value string) {
	if self.Parameters == nil {
		self.Parameters = map[string]string{}
	}
	self.Parameters[key] = value
}

var hopwatchEnabled = true
var currentWebsocket *websocket.Conn
var toBrowserChannel = make(chan command)
var fromBrowserChannel = make(chan command)
var connectChannel = make(chan command)

func init() {
	// see if disable is needed
	for _, arg := range os.Args {
		if arg == "-nobreak" {
			log.Printf("[hopwatch] disabled.\n")
			hopwatchEnabled = false
			return
		}
	}
	http.HandleFunc("/hopwatch.html", html)
	http.HandleFunc("/hopwatch.css", css)
	http.HandleFunc("/hopwatch.js", js)
	http.Handle("/hopwatch", websocket.Handler(connectHandler))
	go listen()
	go sendLoop()
}

// listen starts a Http Server on a fixed port.
// listen is run in parallel to the initialization process such that it does not block.
func listen() {
	log.Printf("[hopwatch] open http://localhost:23456/hopwatch.html ...\n")
	if err := http.ListenAndServe(":23456", nil); err != nil {
		log.Printf("[hopwatch] failed to start listener:%v", err.Error())
	}
}

// connectHandler is a Http handler and is called on loading the debugger in a browser.
// As soon as a command is received the receiveLoop is started. 
func connectHandler(ws *websocket.Conn) {
	log.Printf("[hopwatch] begin accepting commands ...\n")
	if currentWebsocket != nil {
		// reloading an already connected page ; close the old
		currentWebsocket.Close()
		// TODO break receiveLoop
	}
	// remember the connection for the sendLoop	
	currentWebsocket = ws
	var cmd command
	if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
		log.Printf("[hopwatch] JSON.Receive failed:%v", err)
	} else {
		log.Printf("[hopwatch] connected to browser. ready to hop")
		connectChannel <- cmd
		receiveLoop()
	}
	log.Printf("[hopwatch] stop accepting commands.\n")
}

// receiveLoop reads commands from the websocket and puts them onto a channel.
func receiveLoop() {
	for {
		var cmd command
		if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
			log.Printf("[hopwatch] JSON.Receive failed:%v", err)
			break
		}
		fromBrowserChannel <- cmd
	}
}

// sendLoop takes commands from a channel to send to the browser (debugger).
// If no connection is available then wait for it.
// If the command action is quit then abort the loop.
func sendLoop() {
	for {
		next := <-toBrowserChannel
		if next.Action == "quit" {
			break
		}
		if currentWebsocket == nil {
			log.Print("[hopwatch] no browser connection, wait for it ...")
			cmd := <-connectChannel
			if cmd.Action == "quit" {
				break
			}
		}
		websocket.JSON.Send(currentWebsocket, &next)
	}
}

// Watchpoint is a helper to provide a fluent style api.
// This allows for statements like hopwatch.Display("var",value).Break()
type Watchpoint struct {
	disabled bool
	offset   int
}

// Display sends variable name,value pairs to the debugger.
// The parameter nameValuePairs must be even sized.
func Display(nameValuePairs ...interface{}) *Watchpoint {
	_, file, line, ok := runtime.Caller(1)
	cmd := command{Action: "display"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	if len(nameValuePairs)%2 == 0 {
		for i := 0; i < len(nameValuePairs); i += 2 {
			k := nameValuePairs[i]
			v := nameValuePairs[i+1]
			cmd.addParam(fmt.Sprint(k), fmt.Sprint(v))
		}
	} else {
		log.Printf("[hopwatch] WARN: missing variable for Display(...) in: %v:%v\n", file, line)
		return &Watchpoint{disabled: true, offset: 2}
	}
	channelExchangeCommands(cmd)
	return &Watchpoint{offset: 2}
}

// StackOffset (default=2) allows you to change the file indicator in hopwatch.
// Use this method when you wrap the Display(..).Break() in your own function.
func (self *Watchpoint) StackOffset(offset int) *Watchpoint {
	if offset > 0 {
		self.offset = offset
	} else {
		log.Printf("[hopwatch] WARN: illegal stack (caller) offset:%v . Watchpoint is disabled.", offset)
		self.disabled = true
	}
	return self
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func (self Watchpoint) Break(conditions ...bool) {
	if self.disabled {
		return
	}
	Break(self.offset, conditions...)
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func Break(callerOffset int, conditions ...bool) {
	for _, condition := range conditions {
		if !condition {
			return
		}
	}
	_, file, line, ok := runtime.Caller(callerOffset)
	cmd := command{Action: "break"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
		cmd.addParam("go.stack", string(debug.Stack()))
	}
	channelExchangeCommands(cmd)
}

// Put a command on the browser channel and wait for the reply command
func channelExchangeCommands(cmd command) {
	if !hopwatchEnabled {
		log.Printf("[hopwatch] %v", cmd)
		return
	}
	toBrowserChannel <- cmd
	cmd = <-fromBrowserChannel
}
