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
	<script src="http://ajax.googleapis.com/ajax/libs/jquery/1.8.3/jquery.min.js" type="text/javascript"></script>
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
				<a class="buttonEnabled" href="http://go.pkgdoc.org/github.com/emicklei/hopwatch">About</a>
			</div>
		</div>
	</div>
	<div id="page" class="wide">
		<div id="watchlines">
			<table id="output"></table>
		</div>
		<!-- div id="gosource-wrapper">
			<pre id="gosource" />
		</div -->
	</div>
	<div id="footer" style="float:left; width:100%">
		&copy; 2012. <a href="http://github.com/emicklei/hopwatch">hopwatch on github.com</a>
	</div>
</body>
</html>
`)
	return
}
