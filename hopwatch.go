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

// Watchpoint is a helper to provide A fluent style api
type Watchpoint struct{}

func Display(nameValuePairs ...interface{}) *Watchpoint {
	_, file, line, ok := runtime.Caller(1)
	cmd := command{Action: "display"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	for i := 0; i < len(nameValuePairs); i += 2 {
		k := nameValuePairs[i]
		v := nameValuePairs[i+1]
		cmd.addParam(fmt.Sprint(k), fmt.Sprint(v))
	}
	channelExchangeCommands(cmd)
	return new(Watchpoint)
}

func (self Watchpoint) Break(conditions ...bool) {
	Break(conditions...)
}

func Break(conditions ...bool) {
	for _, condition := range conditions {
		if !condition {
			return
		}
	}
	_, file, line, ok := runtime.Caller(1)
	cmd := command{Action: "break"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	channelExchangeCommands(cmd)
}

func channelExchangeCommands(cmd command) {
	if !hopwatchEnabled {
		log.Printf("[hopwatch] %v", cmd)
		return
	}
	toBrowserChannel <- cmd
	cmd = <-fromBrowserChannel
}
