// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
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
var debuggerMutex = sync.Mutex{}

func init() {
	// see if disable is needed. (needed when programs do not call flag.Parse() )
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
	http.HandleFunc("/gosource", gosource)
	http.Handle("/hopwatch", websocket.Handler(connectHandler))
	go listen()
	go sendLoop()
}

// serve a (source) file for displaying in the debugger
func gosource(w http.ResponseWriter, req *http.Request) {
	fileName := req.FormValue("file")
	// should check for permission?
	http.ServeFile(w, req, fileName)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	runtime.Gosched() //yield to http listener
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		cmd = exec.Command("xdg-open", url) //inaccurate assumption but satisfies many BSDs
	}
	log.Printf("[hopwatch] > %s\n", cmd.Args)
	if err := cmd.Start(); err != nil {
		log.Printf("[hopwatch] failed automatic opening of %s - %v, try manually...\n", url, err.Error())
	}
}

// listen starts a Http Server on a fixed port.
// listen is run in parallel to the initialization process such that it does not block.
func listen() {
	port := ":23456"

	go openBrowser("http://localhost" + port + "/hopwatch.html")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("[hopwatch] failed to start listener:%v", err.Error())
	}
}

// connectHandler is a Http handler and is called on loading the debugger in a browser.
// As soon as a command is received the receiveLoop is started.
func connectHandler(ws *websocket.Conn) {
	if currentWebsocket != nil {
		log.Printf("[hopwatch] already connected to a debugger; Ignore this\n")
		return
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
	log.Printf("[hopwatch] end accepting commands.\n")
}

// receiveLoop reads commands from the websocket and puts them onto a channel.
func receiveLoop() {
	for {
		var cmd command
		if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
			log.Printf("[hopwatch] JSON.Receive failed:%v", err)
			break
		}
		if "quit" == cmd.Action {
			hopwatchEnabled = false
			log.Printf("[hopwatch] browser requests disconnect.\n")
			fromBrowserChannel <- cmd
			currentWebsocket.Close() // TODO is not detected by Chrome
			currentWebsocket = nil
			break
		} else {
			fromBrowserChannel <- cmd
		}
	}
}

// sendLoop takes commands from a channel to send to the browser (debugger).
// If no connection is available then wait for it.
// If the command action is quit then abort the loop.
func sendLoop() {
	if currentWebsocket == nil {
		log.Print("[hopwatch] no browser connection, wait for it ...")
		cmd := <-connectChannel
		if "quit" == cmd.Action {
			return
		}
	}
	for {
		next := <-toBrowserChannel
		if "quit" == next.Action {
			break
		}
		if currentWebsocket == nil {
			log.Print("[hopwatch] no browser connection, wait for it ...")
			cmd := <-connectChannel
			if "quit" == cmd.Action {
				break
			}
		}
		websocket.JSON.Send(currentWebsocket, &next)
	}
}

// watchpoint is a helper to provide a fluent style api.
// This allows for statements like hopwatch.Display("var",value).Break()
type Watchpoint struct {
	disabled bool
	offset   int
}

// Printf formats according to a format specifier and writes to the debugger screen.
// It returns a new Watchpoint to send more or break.
func Printf(format string, value ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Printf(format, value...)
}

// Display sends variable name,value pairs to the debugger.
// The parameter nameValuePairs must be even sized.
func Display(nameValuePairs ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Display(nameValuePairs...)
}

// Break suspends the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func Break(conditions ...bool) {
	suspend(2, conditions...)
}

// CallerOffset (default=2) allows you to change the file indicator in hopwatch.
// Use this method when you wrap the .CallerOffset(..).Display(..).Break() in your own function.
func CallerOffset(offset int) *Watchpoint {
	wp := &Watchpoint{offset: offset}
	if offset < 0 {
		log.Panicf("[hopwatch] ERROR: illegal caller offset:%v . watchpoint is disabled.\n", offset)
		wp.disabled = true
	}
	return wp
}

// Printf formats according to a format specifier and writes to the debugger screen.
func (self *Watchpoint) Printf(format string, value ...interface{}) *Watchpoint {
	_, file, line, ok := runtime.Caller(self.offset)
	cmd := command{Action: "print"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	cmd.addParam("line", fmt.Sprintf(format, value...))
	channelExchangeCommands(cmd)
	return self
}

// Display sends variable name,value pairs to the debugger. Values are formatted using %#v.
// The parameter nameValuePairs must be even sized.
func (self *Watchpoint) Display(nameValuePairs ...interface{}) *Watchpoint {
	_, file, line, ok := runtime.Caller(self.offset)
	cmd := command{Action: "display"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	if len(nameValuePairs)%2 == 0 {
		for i := 0; i < len(nameValuePairs); i += 2 {
			k := nameValuePairs[i]
			v := nameValuePairs[i+1]
			cmd.addParam(fmt.Sprint(k), fmt.Sprintf("%#v", v))
		}
	} else {
		log.Printf("[hopwatch] WARN: missing variable for Display(...) in: %v:%v\n", file, line)
		self.disabled = true
		return self
	}
	channelExchangeCommands(cmd)
	return self
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func (self Watchpoint) Break(conditions ...bool) {
	suspend(self.offset, conditions...)
}

// suspend will create a new Command and send it to the browser.
// callerOffset controls from which stackframe the go source file and linenumber must be read.
func suspend(callerOffset int, conditions ...bool) {
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
		cmd.addParam("go.stack", trimStack(string(debug.Stack())))
	}
	channelExchangeCommands(cmd)
}

// Peel off the part of the stack that lives in hopwatch
func trimStack(stack string) string {
	lines := strings.Split(stack, "\n")
	c := 0
	for _, line := range lines {
		if strings.Index(line, "/hopwatch") == -1 { // means no function in this package
			break
		}
		c++
	}
	return strings.Join(lines[c:], "\n")
}

// Put a command on the browser channel and wait for the reply command
func channelExchangeCommands(toCmd command) {
	if !hopwatchEnabled {
		return
	}
	// synchronize command exchange ; break only one goroutine at a time
	debuggerMutex.Lock()
	toBrowserChannel <- toCmd
	_ = <-fromBrowserChannel
	debuggerMutex.Unlock()
}
