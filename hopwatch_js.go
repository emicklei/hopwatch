// Copyright 2012-2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"io"
	"net/http"
)

func js(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	io.WriteString(w, `
	var wsUri = "ws://" + window.location.hostname + ":" + window.location.port + "/hopwatch";
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
		writeToScreen("<-> connected","info mono");		
		sendConnected();
	}
	function onClose(evt) {
		handleDisconnected();
	}
	function onMessage(evt) {
 		try {
            var cmd = JSON.parse(evt.data);
        } catch (e) {
            console.log('[hopwatch] failed to read valid JSON: ', message.data);
            return;
        }		
        // console.log("[hopwatch] received: " + evt.data);
        if (cmd.Action == "display") {
        	var tr = document.createElement("tr");
        	addTime(tr);
        	addGoline(tr,cmd);
        	addMessage(tr,watchParametersToHtml(cmd.Parameters),"watch mono");
        	output.appendChild(tr);
        	tr.scrollIntoView();
        	sendResume();
        	return;
        }
        if (cmd.Action == "print") {
        	var tr = document.createElement("tr");
        	addTime(tr);
        	addGoline(tr,cmd);
        	addMessage(tr,cmd.Parameters["line"],"watch mono");
        	output.appendChild(tr);
        	sendResume();
        	return;
        }        
        if (cmd.Action == "break") {
        	handleSuspended(cmd);
        	return;
        }				        				
	}
	function onError(evt) {
		writeToScreen(evt,"err mono");
	}
	function handleSuspended(cmd) {
        suspended = true;
        document.getElementById("resume").className = "buttonEnabled";
        var tr = document.createElement("tr");
       	addTime(tr);
       	addGoline(tr,cmd);
       	var td = addMessage(tr,"--> program suspended", "suspend mono");
       	         addStack(td,cmd);       	
       	output.appendChild(tr); 
		tr.scrollIntoView();       	
       	loadSource(cmd.Parameters["go.file"], cmd.Parameters["go.line"]);       
	}	
	function writeToScreen(text,cls) {
		var tr = document.createElement("tr");
		addTime(tr);
		addEmptiness(tr);
		addMessage(tr,text,cls)
		output.appendChild(tr);
	}	
	function addTime(tr) {
		var stamp = document.createElement("td");
		stamp.innerHTML = timeHHMMSS();
		stamp.className = "time mono"
		tr.appendChild(stamp);			
	}	
	function addMessage(tr,msg,msgcls) {
		var txt = document.createElement("td");
		txt.className = msgcls		
		txt.innerHTML = msg;
		tr.appendChild(txt);
		return txt;
	}
	function addEmptiness(tr) {
		var empt = document.createElement("td");
		empt.className = "goline"		
		empt.innerHTML = "&nbsp;";
		tr.appendChild(empt);
	}
	function addGoline(tr,cmd) {
		var where = document.createElement("td");		
		var link = document.createElement("a");
		link.href = "#";
		link.className = "goline mono";
		link.onclick = function() { 
			loadSource(cmd.Parameters["go.file"], cmd.Parameters["go.line"]); 
		};
		link.innerHTML = goline(cmd.Parameters);
		where.appendChild(link);
		tr.appendChild(where);
	}
	function loadSource(fileName, nr) {
		$("#gofile").html(shortenFileName(fileName));
		$("#gosource-pane").show();
		$("#gosource").load("/gosource?file="+fileName, function(responseText,status,xhr) {
			handleSourceLoaded(responseText,nr);
			}
		);
	}
	function handleSourceLoaded(responseText,line) {
		var nrs = $("#nrs");		
		nrs.empty();
		// Insert line numbers		
		var arr = responseText.split('\n');
		var breakElm;		
        for (var i = 0; i < arr.length; i++) {
        	var nr = i+1;        	
        	var elm = document.createElement("div");
        	elm.innerHTML = nr;
        	if (line == nr) {
        		elm.className = "break";
        		breakElm = elm;
        	} 
        	nrs.append(elm)
        }
		$("#gosource").text(responseText);
		breakElm.scrollIntoView();            
	}
	function shortenFileName(fileName) {
		return fileName.length > 48 ? "..." + fileName.substring(fileName.length - 48) : fileName;
	}
	function addStack(td,cmd) {
		var stack = cmd.Parameters["go.stack"];
		if (stack != null && stack.length > 0) {
			addNonEmptyStackTo(stack,td);
		}	
	}	
	function addNonEmptyStackTo(stack,td) {
		var toggle = document.createElement("a");
		toggle.href = "#";
		toggle.className = "toggle";
		toggle.onclick = function() { toggleStack(toggle); };
		toggle.innerHTML="stack &#x25B6;";
		td.appendChild(toggle);
		
		var stk = document.createElement("div");
		stk.style.display = "none";
		var lines = document.createElement("pre");
		lines.innerHTML = stack	
		lines.className = "stack mono"			
		stk.appendChild(lines)		
		td.appendChild(stk)	
	}
	function toggleStack(link) {
		var stack = link.nextSibling;
		if (stack.style.display == "none") {	
			link.innerHTML = "stack &#x25BC;";	
			stack.style.display = "block"
			stack.scrollIntoView();
		} else {		
			link.innerHTML = "stack &#x25B6;";
			stack.style.display = "none";
		}
	}	
	// http://www.quirksmode.org/js/keys.html
	function handleKeyDown(event) {
		if (event.keyCode == 119) {
			actionResume();			
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
	function actionResume() {
		if (!connected) return;
		if (!suspended) return;
		suspended = false;
		document.getElementById("resume").className = "buttonDisabled";
		// writeToScreen("<-- resume program","info mono");
		sendResume();
	}
	function actionDisconnect() {
		if (!connected) return;
		connected = false;
		document.getElementById("disconnect").className = "buttonDisabled";
		sendQuit();
		writeToScreen("<-- disconnect requested","info mono");
		websocket.close();  // seems not to trigger close on Go-side ; so handleDisconnected cannot be used here
	}
	function handleDisconnected() {
		connected = false;
		document.getElementById("resume").className = "buttonDisabled";
		document.getElementById("disconnect").className = "buttonDisabled";
		writeToScreen(">-< disconnected","info mono");	
	}
	function timeHHMMSS()    { return new Date().toTimeString().replace(/.*(\d{2}:\d{2}:\d{2}).*/, "$1"); }
	function sendConnected() { doSend('{"Action":"connected"}'); }
	function sendResume()    { doSend('{"Action":"resume"}'); }
	function sendQuit()      { doSend('{"Action":"quit"}'); }	
	function doSend(message) {
		// ^console.log("[hopwatch] send: " + message);
		websocket.send(message);
	}
	window.addEventListener("load", init, false);
	window.addEventListener("keydown", handleKeyDown, true); `)
	return
}
