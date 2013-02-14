// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"io"
	"net/http"
)

func html(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w,
		`<!DOCTYPE html>
<meta charset="utf-8" />
<title>Hopwatch Debugger</title>
<head>
	<link href="hopwatch.css" rel="stylesheet" type="text/css" >	
	<script src="http://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js" type="text/javascript"></script>
	<script type="text/javascript" src="hopwatch.js" ></script>
</head>
<body>
	<div id="topbar">
		<div class="container wide">
			<div id="heading">
				<a href="/hopwatch.html">Hopwatch - debugging tool</a>
			</div>		
			<div id="menu">
				<a id="resume" class="buttonDisabled" href="javascript:actionResume();">F8 - Resume</a>
				<a id="disconnect" class="buttonDisabled" href="javascript:actionDisconnect();">Disconnect</a>
				<a class="buttonEnabled" href="http://go.pkgdoc.org/github.com/emicklei/hopwatch" target="_blank">About</a>
			</div>
		</div>
	</div>
	<div id="page" class="wide">
		<div id="log-pane">
			<table id="output"></table>
		</div>		
		<div id="gosource-pane">
			<div id="gofile">somefile.go</div>
			<div id="nrs" class="mono">
				<div>1</div>
				<div>2</div>
				<div>3</div>
				<div>4</div>
			</div>			
			<pre id="gosource">
				loading go source...
			</pre>
		</div>
	</div>
	<div id="footer" style="float:left; width:100%">
		&copy; 2012-2013. <a href="http://github.com/emicklei/hopwatch" target="_blank">hopwatch on github.com</a>
	</div>
</body>
</html>
`)
	return
}
