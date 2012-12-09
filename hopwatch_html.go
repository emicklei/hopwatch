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
	<h2>Hopwatch Go Debugger</h2>
	<p>
		<a href="javascript:sendResume();">[ Resume ]</a>
		<!-- a href="javascript:sendQuit();">[ Disconnect ]</a -->				
	</p>
	<table id="output" style="width:100%"></table>
</body>
</html>
`)
	return
}
