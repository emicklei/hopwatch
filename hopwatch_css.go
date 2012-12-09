package hopwatch

import (
	"io"
	"net/http"
)

func css(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	io.WriteString(w, `
	.time    {background-color:#DDD;font-family:"Lucida Console", Monaco, monospace;font-size:small;white-space:nowrap}
	.watch 	 {background-color:#FFF;font-family:"Lucida Console", Monaco, monospace;font-size:small;width:100%;}
	.goline  {background-color:#FFF;color:#888;font-family:"Lucida Console", Monaco, monospace;font-size:small;}
	.err 	 {background-color:#FF3300;}
	.info 	 {background-color:#CCFFCC;}
	.suspend {background-color:#E0EBF5;}
	body 	 {background-color:#EEE;}`)
	return
}
