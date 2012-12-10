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
	<script type="text/javascript" src="hopwatch.js" ></script>
</head>
<body>
	<div id="topbar">
		<div class="container wide">
			<div id="heading">
				<a href="/hopwatch.html">Hopwatch - debugger</a>
			</div>		
			<div id="menu">
				<a href="javascript:sendResume();">F8 - Resume</a>
				<a href="javascript:sendQuit();">Disconnect</a>
				<a href="http://go.pkgdoc.org/github.com/emicklei/hopwatch">About</a>
			</div>
		</div>
	</div>
	<div id="page" class="wide">
		<table id="output" style="width:100%"></table>
	</div>
	<div id="footer">
		<a href="http://github.com/emicklei/hopwatch">hopwatch on github.com</a>
	</div>
</body>
</html>
`)
	return
}
