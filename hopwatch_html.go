package hopwatch

import (
	"io"
	"net/http"
)

func writePage(w http.ResponseWriter, req *http.Request) {
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
