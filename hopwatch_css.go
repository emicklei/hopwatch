package hopwatch

import (
	"io"
	"net/http"
)

func css(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	io.WriteString(w, `
	.time   {background-color:#DDD;font-family:"Lucida Console", Monaco, monospace;font-size:small;white-space:nowrap}
	.watch 	{background-color:#FFF;font-family:"Lucida Console", Monaco, monospace;font-size:small;width:100%;}
	.err 	{background-color:#FF3300;}
	.info 	{background-color:#CCFFCC;}
	body 	{background-color:#EEE;}`)
	return
}
