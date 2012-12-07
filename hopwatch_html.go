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
		writeToScreen(evt.data,"watch");		
	}
	function onError(evt) {
		writeToScreen(evt,"err");
	}	
	function writeToScreen(message,cls) {
		var tr = document.createElement("tr");
		var stamp = document.createElement("td");
		stamp.innerHTML = new Date().toTimeString().replace(/.*(\d{2}:\d{2}:\d{2}).*/, "$1");
		stamp.className = "time"
		tr.appendChild(stamp);		
		var td = document.createElement("td");
		td.className = cls		
		td.innerHTML = message;
		tr.appendChild(td);
		output.appendChild(tr);
	}	
	function sendConnected() { doSend('{"Action":"connected"}'); }
	function sendProceed()   { doSend('{"Action":"proceed"}'); }
	function sendQuit()      { doSend('{"Action":"quit"}'); }	
	function doSend(message) {
		console.log("[hopwatch] send: " + message);
		websocket.send(message);
	}
	window.addEventListener("load", init, false);
</script>
<head>
	<style>
	.time   {background-color:#DDD;font-family:"Lucida Console", Monaco, monospace;font-size:small;white-space:nowrap}
	.watch 	{background-color:#FFF;font-family:"Lucida Console", Monaco, monospace;font-size:small;width:100%;}
	.err 	{background-color:#FF3300;}
	.info 	{background-color:#CCFFCC;}
	body 	{background-color:#EEE;}
	</style>
</head>
<body>
	<h2>Hopwatch Debugger</h2>
	<p><a href="javascript:sendProceed();">[ Proceed ]</a><a href="javascript:sendQuit();">[ Quit ]</a></p>
	<table id="output" style="width:100%"></table>
	<h7><a href="https://github.com/emicklei/hopwatch">hopwatch on github</a></h7>
</body>
</html>
`)
	return
}
