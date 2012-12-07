package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
)

type command struct {
	Action     string
	Parameters map[string]string
}

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

func listen() {
	log.Printf("[hopwatch] open http://localhost:23456/hopwatch.html ...\n")
	if err := http.ListenAndServe(":23456", nil); err != nil {
		log.Printf("[hopwatch] failed to start listener:%v", err.Error())
	}
}

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

// Break stops the execution of the program and passes values to the Hopwatch page to show.
// The execution of the program is resumed after receiving the proceed command. 
func Break(location string, key string, value interface{}) {
	if !hopwatchEnabled {
		log.Printf("%v: %v = %v", location, key, value)
		return
	}
	cmd := command{Action: "watch"}
	cmd.addParam("go.location", location)
	cmd.addParam(key, fmt.Sprint(value))
	toBrowserChannel <- cmd
	cmd = <-fromBrowserChannel
}
