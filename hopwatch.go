package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"flag"
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

var hopwatchParam = flag.Bool("hopwatch", true, "controls whether hopwatch agent is started")
var hopwatchEnabled = true
var currentWebsocket *websocket.Conn
var toBrowserChannel = make(chan command)
var fromBrowserChannel = make(chan command)
var connectChannel = make(chan command)

func init() {
	// see if disable is needed
	for i, arg := range os.Args {
		if arg == "-hopwatch" && i < len(os.Args) && os.Args[i+1] == "false" {
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
	if currentWebsocket != nil {
		// reloading an already connected page ; close the old		
		currentWebsocket.Close()
		log.Printf("[hopwatch] closed old connection.\n")
		// TODO break receiveLoop
	}
	log.Printf("[hopwatch] begin accepting commands ...\n")
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
	if currentWebsocket == nil {
		log.Print("[hopwatch] no browser connection, wait for it ...")
		cmd := <-connectChannel
		if cmd.Action == "quit" {
			return
		}
	}
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

// CallerOffset (default=2 and must be positive) allows you to change the file indicator in hopwatch.
// Use this method when you wrap the Display(..).CallerOffset(2+1).Break() in your own function.
func (self *Watchpoint) CallerOffset(offset int) *Watchpoint {
	if self.disabled {
		return self
	}
	if offset > 0 {
		self.offset = offset
	} else {
		log.Printf("[hopwatch] WARN: illegal caller offset:%v . Watchpoint is disabled.\n", offset)
		self.disabled = true
	}
	return self
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func (self Watchpoint) Break(conditions ...bool) {
	Break(self.offset, conditions...)
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
// callerOffset controls from which stackframe the go source file and linenumber must be read. For direct use of this function, set the offset to 1.
func Break(callerOffset int, conditions ...bool) {
	if callerOffset < 0 {
		log.Printf("[hopwatch] WARN: illegal caller offset:%v . Break is ineffective.\n", callerOffset)
		return
	}
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
func channelExchangeCommands(toCmd command) {
	if !hopwatchEnabled {
		return
	}
	toBrowserChannel <- toCmd
	_ = <-fromBrowserChannel
	//	if fromCmd == "resume" {
	//		channelExchangeCommands(Command{Action:"status"}
	//	}
}
