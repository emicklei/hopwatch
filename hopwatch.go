// hopwatch is a browser-based tool for debugging Go programs.

// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license 
// that can be found in the LICENSE file.
package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
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
	http.HandleFunc("/hopwatch.html", writePage)
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

// watchpoint is the object that collects information to watch to in the debugger.
type watchpoint struct {
	file string
	line string
	vars map[string]string
}

// Watch will add a variable and its value to the list of variables to watch next in the debugger
// It will not be send to the debugger until Break is called on the watchpoint
func Watch(variableName string, value interface{}) *watchpoint {
	wp := new(watchpoint)
	wp.vars = map[string]string{}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		wp.file = file
		wp.line = fmt.Sprint(line)
	}
	wp.vars[variableName] = fmt.Sprint(value)
	return wp
}

// Watch will add a variable and its value to the list of variables to watch next in the debugger
// It will not be send to the debugger until Break is called on the watchpoint
func (self *watchpoint) Watch(variableName string, value interface{}) *watchpoint {
	self.vars[variableName] = fmt.Sprint(value)
	return self
}

// Break stops the execution of the program and passes the current watchpoint to the Hopwatch page to show.
// The execution of the program is resumed after receiving the proceed command. 
func (self *watchpoint) Break() {
	if !hopwatchEnabled {
		log.Printf("[hopwatch] %v", self)
		return
	}
	cmd := command{Action: "watch"}
	cmd.addParam("go.file", self.file)
	cmd.addParam("go.line", self.line)
	for k, v := range self.vars {
		cmd.addParam(k, v)
	}
	toBrowserChannel <- cmd
	cmd = <-fromBrowserChannel
}
