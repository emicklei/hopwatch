package hopwatch

import (
	"io"
	"net/http"
)

func js(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	io.WriteString(w, `
	var wsUri = "ws://localhost:23456/hopwatch";
	var output;
	var websocket = new WebSocket(wsUri);	
	function init() {
		output = document.getElementById("output");
		setupWebSocket();
	}
	function setupWebSocket() {		
		websocket.onopen = function(evt) { onOpen(evt) };
		websocket.onclose = function(evt) { onClose(evt) };
		websocket.onmessage = function(evt) { onMessage(evt) };
		websocket.onerror = function(evt) { onError(evt) };
	}
	function onOpen(evt) {
		writeToScreen("connected","info");
		sendConnected();
	}
	function onClose(evt) {
		writeToScreen("disconnected","info");
	}
	function onMessage(evt) {
 		try {
            var cmd = JSON.parse(evt.data);
        } catch (e) {
            console.log('[hopwatch] failed to read valid JSON: ', message.data);
            return;
        }		
        console.log("[hopwatch] received: " + evt.data);
        if (cmd.Action == "display") {
        	actionDisplay(cmd);
        	sendResume();
        }
        if (cmd.Action == "break") {
        	writeToScreen("break","info")        	
        }				        				
	}
	function onError(evt) {
		writeToScreen(evt,"err");
	}	
	function actionDisplay(cmd) {
		var tr = document.createElement("tr");
		var stamp = document.createElement("td");
		stamp.innerHTML = new Date().toTimeString().replace(/.*(\d{2}:\d{2}:\d{2}).*/, "$1");
		stamp.className = "time"
		tr.appendChild(stamp);		
		var td = document.createElement("td");
		td.className = "watch"		
		td.innerHTML = watchParametersToHtml(cmd.Parameters);
		tr.appendChild(td);
		output.appendChild(tr);
	}
	function writeToScreen(text,cls) {
		var tr = document.createElement("tr");
		var stamp = document.createElement("td");
		stamp.innerHTML = new Date().toTimeString().replace(/.*(\d{2}:\d{2}:\d{2}).*/, "$1");
		stamp.className = "time"
		tr.appendChild(stamp);		
		var td = document.createElement("td");
		td.className = cls		
		td.innerHTML = text;
		tr.appendChild(td);
		output.appendChild(tr);
	}
	function watchParametersToHtml(parameters) {
		var f = parameters["go.file"]
		f = f.substr(f.lastIndexOf("/")+1)
		var line = f + ":" + parameters["go.line"] + " ";
		for (var prop in parameters) {
			if (prop.slice(0,3) != "go.") {
				line = line + prop + "=" + parameters[prop] + ","
			}
		} 
		return line
	}
	function sendConnected() { doSend('{"Action":"connected"}'); }
	function sendResume()   { doSend('{"Action":"resume"}'); }
	function sendQuit()      { doSend('{"Action":"quit"}'); }	
	function doSend(message) {
		console.log("[hopwatch] send: " + message);
		websocket.send(message);
	}
	window.addEventListener("load", init, false);`)
	return
}
