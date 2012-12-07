package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"log"
	"io"
	"fmt"
)

type Command struct {
	Action string
	Parameters map[string]string
}
func (self *Command) addParam(key,value string) {
	if self.Parameters == nil {
		self.Parameters = map[string]string{}
	}
	self.Parameters[key] = value 
}

var currentWebsocket *websocket.Conn
var toBrowserChannel = make(chan Command)
var fromBrowserChannel = make(chan Command)

func init() {
	http.HandleFunc("/hopwatch.html", watcher)
	http.Handle("/hopwatch", websocket.Handler(receiveLoop))
	go listen()
	go sendLoop()
}

func listen() {	
	log.Printf("[hopwatch] open http://localhost:23456/hopwatch.html ...\n")
	if err := http.ListenAndServe(":23456", nil) ; err != nil {
		log.Printf("[hopwatch] failed to start listener:%v",err.Error())
	} 	
}

func receiveLoop(ws *websocket.Conn) {
	log.Printf("[hopwatch] begin accepting commands...\n")
	// remember the connection for the sendLoop	
	currentWebsocket = ws
	for {
		var cmd Command 
		if err := websocket.JSON.Receive(currentWebsocket, &cmd) ; err != nil {
			log.Print("[hopwatch] JSON.Receive failed")
			break	
		}	
		fromBrowserChannel <- cmd
	}
	log.Printf("[hopwatch] stop accepting commands.\n")
}

func sendLoop() {
	for {
		next := <- toBrowserChannel
		if next.Action == "quit" {
			break
		}
		if currentWebsocket == nil {
			log.Print("[hopwatch] no browser connection, wait for it ...")	
			_ = <- fromBrowserChannel	
		} else {
			websocket.JSON.Send(currentWebsocket, &next)
		}
	}
}
// Break stops the execution of the program and passes values to the Hopwatch page to show.
// The execution of the program is resumed after receiving the proceed command. 
func Break(line string, key string, value interface{}) {
	cmd := Command{Action:"watch"}
	cmd.addParam("go.line", line)
	cmd.addParam(key, fmt.Sprint(value))
	
	toBrowserChannel <- cmd
	cmd = <- fromBrowserChannel 
}

func watcher(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w,
	`<!DOCTYPE html>
<meta charset="utf-8" />
<title>Hopwatch Debugger</title>
<script language="javascript" type="text/javascript">
	var wsUri = "ws://localhost:23456/hopwatch";
	var output;
	var websocket = new WebSocket(wsUri);	
	function init() {
		output = document.getElementById("output");
		testWebSocket();
	}
	function testWebSocket() {		
		websocket.onopen = function(evt) {
			onOpen(evt)
		};
		websocket.onclose = function(evt) {
			onClose(evt)
		};
		websocket.onmessage = function(evt) {
			onMessage(evt)
		};
		websocket.onerror = function(evt) {
			onError(evt)
		};
	}
	function onOpen(evt) {
		writeToScreen("CONNECTED");
		doSend('{"Action":"browser"}');
	}
	function onClose(evt) {
		writeToScreen("DISCONNECTED");
	}
	function onMessage(evt) {
 		try {
            var json = JSON.parse(evt.data);
        } catch (e) {
            console.log('This doesn\'t look like a valid JSON: ', message.data);
            return;
        }		
		writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data
				+ '</span>');		
	}
	function onError(evt) {
		writeToScreen('<span style="color: red;">ERROR:</span> ' + evt);
	}
	function doSend(message) {
		writeToScreen("SENT: " + message);
//		if (websocket == null) {
//			writeToScreen("ws == null");
//		} else {
			websocket.send(message);
//		}
	}	
	function writeToScreen(message) {
		var pre = document.createElement("p");
		pre.style.wordWrap = "break-word";
		pre.innerHTML = message;
		output.appendChild(pre);
	}
	
	function proceed() {
		doSend('{"Action":"proceed"}');
	}
	
	window.addEventListener("load", init, false);
</script>
<h2>Hopwatch Debugger</h2>
<a href="javascript:proceed();">Proceed</a>
<div id="output"></div>
</html>
`)
	return
}