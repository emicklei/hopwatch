// Copyright 2012,2014 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	hopwatchServerAddressParam = flag.String("hopwatch.server", "", "HTTP host:port server running hopwatch server")
	hopwatchHostParam          = flag.String("hopwatch.host", "localhost", "HTTP host the debugger is listening on")
	hopwatchPortParam          = flag.Int("hopwatch.port", 23456, "HTTP port the debugger is listening on")
	hopwatchParam              = flag.Bool("hopwatch", true, "controls whether hopwatch agent is started")
	hopwatchOpenParam          = flag.Bool("hopwatch.open", true, "controls whether a browser page is opened on the hopwatch page")
	hopwatchBreakParam         = flag.Bool("hopwatch.break", true, "do not suspend the program if Break(..) is called")

	hopwatchEnabled             = true
	hopwatchOpenEnabled         = true
	hopwatchBreakEnabled        = true
	hopwatchHost                = "localhost"
	hopwatchPort          int64 = 23456
	hopwatchServerAddress       = ""

	currentWebsocket   *websocket.Conn
	toBrowserChannel   = make(chan command)
	fromBrowserChannel = make(chan command)
	connectChannel     = make(chan command)
	debuggerMutex      = sync.Mutex{}
)

func init() {
	// check any command line params. (needed when programs do not call flag.Parse() )
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-hopwatch=") {
			if strings.HasSuffix(arg, "false") {
				log.Printf("[hopwatch] disabled.\n")
				hopwatchEnabled = false
				return
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.open") {
			if strings.HasSuffix(arg, "false") {
				log.Printf("[hopwatch] auto open debugger disabled.\n")
				hopwatchOpenEnabled = false
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.break") {
			if strings.HasSuffix(arg, "false") {
				log.Printf("[hopwatch] suspend on Break(..) disabled.\n")
				hopwatchBreakEnabled = false
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.host") {
			if eq := strings.Index(arg, "="); eq != -1 {
				hopwatchHost = arg[eq+1:]
			} else if i < len(os.Args) {
				hopwatchHost = os.Args[i+1]
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.server") {
			if eq := strings.Index(arg, "="); eq != -1 {
				hopwatchServerAddress = arg[eq+1:]
			} else if i < len(os.Args) {
				hopwatchServerAddress = os.Args[i+1]
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.port") {
			portString := ""
			if eq := strings.Index(arg, "="); eq != -1 {
				portString = arg[eq+1:]
			} else if i < len(os.Args) {
				portString = os.Args[i+1]
			}
			port, err := strconv.ParseInt(portString, 10, 64)
			if err != nil {
				log.Panicf("[hopwatch] illegal port parameter:%v", err)
			}
			hopwatchPort = port
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

// Open calls the OS default program for uri
func open(uri string) error {
	var run string
	switch {
	case "windows" == runtime.GOOS:
		run = "start"
	case "darwin" == runtime.GOOS:
		run = "open"
	case "linux" == runtime.GOOS:
		run = "xdg-open"
	default:
		return fmt.Errorf("Unable to open uri:%v on:%v", uri, runtime.GOOS)
	}
	return exec.Command(run, uri).Start()
}

// serve a (source) file for displaying in the debugger
func gosource(w http.ResponseWriter, req *http.Request) {
	fileName := req.FormValue("file")
	// should check for permission?
	w.Header().Set("Cache-control", "no-store, no-cache, must-revalidate")
	http.ServeFile(w, req, fileName)
}

// listen starts a Http Server on a fixed port.
// listen is run in parallel to the initialization process such that it does not block.
func listen() {
	hostPort := fmt.Sprintf("%s:%d", hopwatchHost, hopwatchPort)
	if hopwatchOpenEnabled {
		log.Printf("[hopwatch] opening http://%v/hopwatch.html ...\n", hostPort)
		go open(fmt.Sprintf("http://%v/hopwatch.html", hostPort))
	} else {
		log.Printf("[hopwatch] open http://%v/hopwatch.html ...\n", hostPort)
	}
	log.Printf("[hopwatch] listening to %v\n", hostPort)
	if err := http.ListenAndServe(hostPort, nil); err != nil {
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
	// remember the connection for the sendLoop
	currentWebsocket = ws
	var cmd command
	if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
		log.Printf("[hopwatch] connectHandler.JSON.Receive failed:%v", err)
	} else {
		log.Printf("[hopwatch] connected to browser. ready to hop")
		connectChannel <- cmd
		receiveLoop()
	}
}

// receiveLoop reads commands from the websocket and puts them onto a channel.
func receiveLoop() {
	for {
		var cmd command
		if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
			log.Printf("[hopwatch] receiveLoop.JSON.Receive failed:%v", err)
			fromBrowserChannel <- command{Action: "quit"}
			break
		}
		if "quit" == cmd.Action {
			hopwatchEnabled = false
			log.Printf("[hopwatch] browser requests disconnect.\n")
			currentWebsocket.Close()
			currentWebsocket = nil
			fromBrowserChannel <- cmd
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
		log.Print("[hopwatch-exchange] no browser connection, wait for it ...")
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
			log.Print("[hopwatch-exchange] no browser connection, wait for it ...")
			cmd := <-connectChannel
			if "quit" == cmd.Action {
				break
			}
		}
		websocket.JSON.Send(currentWebsocket, &next)
	}
}

// suspend will create a new Command and send it to the browser.
// callerOffset controls from which stackframe the go source file and linenumber must be read.
// Ignore if option hopwatch.break=false
func suspend(callerOffset int, conditions ...bool) {
	if !hopwatchBreakEnabled {
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
		cmd.addParam("go.stack", trimStack(string(debug.Stack()), fmt.Sprintf("%s:%d", file, line)))
	}
	channelExchangeCommands(cmd)
}

// Peel off the part of the stack that lives in hopwatch
func trimStack(stack, fileAndLine string) string {
	lines := strings.Split(stack, "\n")
	c := 0
	for _, each := range lines {
		if strings.Index(each, fileAndLine) != -1 {
			break
		}
		c++
	}
	return strings.Join(lines[4:], "\n")
}

// Put a command on the browser channel and wait for the reply command
func channelExchangeCommands(toCmd command) {
	if !hopwatchEnabled {
		return
	}
	// synchronize command exchange ; break only one goroutine at a time
	debuggerMutex.Lock()
	toBrowserChannel <- toCmd
	<-fromBrowserChannel
	debuggerMutex.Unlock()
}
