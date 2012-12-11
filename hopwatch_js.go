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
	var connected = false;
	var suspended = false;
	
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
		connected = true;
		document.getElementById("disconnect").className = "buttonEnabled";
		writeToScreen("<-> connected","info");		
		sendConnected();
	}
	function onClose(evt) {
		connected = false;
		document.getElementById("disconnect").className = "buttonDisabled";
		writeToScreen(">-< disconnected","info");
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
        	handleSuspended(cmd);
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
		
		if (stack != null && stack.length > 0) {
			addNonEmptyStackTo(stack,txt);
		}
		
		output.appendChild(tr);		
	}
	function addNonEmptyStackTo(stack,textCell) {
		var toggle = document.createElement("a");
		toggle.href = "#";
		toggle.className = "toggle";
		toggle.onclick = function() { toggleStack(toggle); };
		toggle.innerHTML="stack";
		textCell.appendChild(toggle);
		
		var stk = document.createElement("div");
		stk.style.display = "none";
		var lines = document.createElement("pre");
		lines.innerHTML = stack	
		lines.className = "stack"			
		stk.appendChild(lines)		
		textCell.appendChild(stk)	
	}
	function toggleStack(link) {
		var stack = link.nextSibling;
		stack.style.display = (stack.style.display == "none") ? "block" : "none";
	}	
	function writeToScreen(text,cls) {
		row(timeHHMMSS(), "hopwatch", "", text ,cls)
	}
	// http://www.quirksmode.org/js/keys.html
	function handleKeyDown(event) {
		if (event.keyCode == 119) {
			actionResume();
			writeToScreen("program resumed","info");  // Really should ask Go for this status
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
	function handleSuspended(cmd) {
        suspended = true;
        document.getElementById("resume").className = "buttonEnabled";
        row(timeHHMMSS(), goline(cmd.Parameters), cmd.Parameters["go.stack"], " program suspended", "suspend")	
	}
	function actionResume() {
		if (!connected) return;
		if (!suspended) return;
		suspended = false;
		document.getElementById("resume").className = "buttonDisabled";
		sendResume();
	}
	function actionDisconnect() {
		if (!connected) return;
		connected = false;
		document.getElementById("disconnect").className = "buttonDisabled";
		sendQuit();
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
	window.addEventListener("keydown", handleKeyDown, true); `)
	return
}
