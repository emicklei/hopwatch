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
        	row(timeHHMMSS(), goline(cmd.Parameters), "", watchParametersToHtml(cmd.Parameters), "watch")
        	sendResume();
        }
        if (cmd.Action == "break") {
        	row(timeHHMMSS(), goline(cmd.Parameters), "", " program suspended", "suspend")
        }				        				
	}
	function onError(evt) {
		writeToScreen(evt,"err");
	}
	function row(time,goline,stack,msg,msgcls) {
		var tr = document.createElement("tr");
		
		var stamp = document.createElement("td");
		stamp.innerHTML = time;
		stamp.className = "time"
		tr.appendChild(stamp);	
			
		var where = document.createElement("td");
		where.className = "goline"		
		where.innerHTML = goline;
		tr.appendChild(where);	

		var txt = document.createElement("td");
		txt.className = msgcls		
		txt.innerHTML = msg;
		tr.appendChild(txt);
		
		output.appendChild(tr);		
	}
	function writeToScreen(text,cls) {
		row(timeHHMMSS(), "", "", text ,cls)
	}
	// http://www.quirksmode.org/js/keys.html
	function handleKeyDown(event) {
		console.log(event.keyCode);
		if (event.keyCode == 199) {
			sendResume();
		}
	}
	function watchParametersToHtml(parameters) {
		var line = "";
		var multiline = false;
		for (var prop in parameters) {
			if (prop.slice(0,3) != "go.") {				
				if (multiline) { line = line + ", "; }
				line = line + prop + "=" + parameters[prop];
				multiline = true;
			}
		} 
		return line
	}
	function goline(parameters) { 
		var f = parameters["go.file"]
		f = f.substr(f.lastIndexOf("/")+1)
		return f + ":" + parameters["go.line"]
	}
	function timeHHMMSS() { return new Date().toTimeString().replace(/.*(\d{2}:\d{2}:\d{2}).*/, "$1"); }
	function sendConnected() { doSend('{"Action":"connected"}'); }
	function sendResume()   { doSend('{"Action":"resume"}'); }
	function sendQuit()      { doSend('{"Action":"quit"}'); }	
	function doSend(message) {
		console.log("[hopwatch] send: " + message);
		websocket.send(message);
	}
	window.addEventListener("load", init, false);
	window.addEventListeners("keydown", handleKeyDown, false); `)
	return
}
