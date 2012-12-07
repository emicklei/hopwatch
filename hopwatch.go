package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
)

type Command struct {
	Action     string
	Parameters map[string]string
}

func (self *Command) addParam(key, value string) {
	if self.Parameters == nil {
		self.Parameters = map[string]string{}
	}
	self.Parameters[key] = value
}

var currentWebsocket *websocket.Conn
var toBrowserChannel = make(chan Command)
var fromBrowserChannel = make(chan Command)

func init() {
	http.HandleFunc("/hopwatch.html", writePage)
	http.Handle("/hopwatch", websocket.Handler(receiveLoop))
	go listen()
	go sendLoop()
}

func listen() {
	log.Printf("[hopwatch] open http://localhost:23456/hopwatch.html ...\n")
	if err := http.ListenAndServe(":23456", nil); err != nil {
		log.Printf("[hopwatch] failed to start listener:%v", err.Error())
	}
}

func receiveLoop(ws *websocket.Conn) {
	log.Printf("[hopwatch] begin accepting commands...\n")
	// remember the connection for the sendLoop	
	currentWebsocket = ws
	for {
		var cmd Command
		if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
			log.Print("[hopwatch] JSON.Receive failed")
			break
		}
		fromBrowserChannel <- cmd
	}
	log.Printf("[hopwatch] stop accepting commands.\n")
}

func sendLoop() {
	for {
		next := <-toBrowserChannel
		if next.Action == "quit" {
			break
		}
		if currentWebsocket == nil {
			log.Print("[hopwatch] no browser connection, wait for it ...")
			_ = <-fromBrowserChannel
		} else {
			websocket.JSON.Send(currentWebsocket, &next)
		}
	}
}

// Break stops the execution of the program and passes values to the Hopwatch page to show.
// The execution of the program is resumed after receiving the proceed command. 
func Break(location string, key string, value interface{}) {
	cmd := Command{Action: "watch"}
	cmd.addParam("go.location", location)
	cmd.addParam(key, fmt.Sprint(value))

	toBrowserChannel <- cmd
	cmd = <-fromBrowserChannel
}
