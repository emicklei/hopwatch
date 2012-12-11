package hopwatch

import (
	"io"
	"net/http"
)

func css(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	io.WriteString(w, `
	body {
		margin: 0;
		font-family: Helvetica, Arial, sans-serif;
		font-size: 16px;
		color: #222;
	}
	td { vertical-align:top }
	a { text-decoration:none; color: #375EAB; }
	div#heading {
		float: left;
		margin: 0 0 10px 0;
		padding: 21px 0;
		font-size: 20px;
		font-weight: normal;
	}
	div#heading a {
		color: #222;
		text-decoration: none;
	}	
	div#topbar {
		background: #E0EBF5;
		height: 64px;
	}
	div#page,
	div#topbar > .container {
		clear: both;
		text-align: left;
		margin-left: auto;
		margin-right: auto;
		padding: 0 20px;
		width: 900px;
	}
	div#page,
	div#topbar > .wide {
		width: auto;
	}	
	div#menu {
		float: left;
		min-width: 590px;
		padding: 10px 0;
		text-align: right;
		margin-top: 10px;
	}
	div#menu > a {
		margin-right: 5px;
		margin-bottom: 10px;
		padding: 10px;				
	}
	.buttonEnabled {
		color: white;
		background: #375EAB;
	}
	.buttonDisabled {
		color: #375EAB;
		background: white;
	}	
	div#menu > a,
	div#menu > input {
		padding: 10px;	
		text-decoration: none;
		font-size: 16px;	
		-webkit-border-radius: 5px;
		-moz-border-radius: 5px;
		border-radius: 5px;
	}
	div#footer {
		text-align: center;
		color: #666;
		font-size: 14px;
		margin: 40px 0;
	}		
	.toggle  {padding-left:4px;padding-right:4px;margin-left:4px;margin-right:4px;background-color:#375EAB;color:#FFF;}	
	.stack   {font-family:"Lucida Console", Monaco, monospace;font-size:small;background-color:#F8EEB9;border:1;padding:4px}
	.time    {font-family:"Lucida Console", Monaco, monospace;font-size:small;color:#AAA;white-space:nowrap}
	.watch 	 {font-family:"Lucida Console", Monaco, monospace;font-size:small;width:100%;}
	.goline  {font-family:"Lucida Console", Monaco, monospace;font-size:small;color:#888;padding-left:8px;padding-right:8px;}
	.err 	 {font-family:"Lucida Console", Monaco, monospace;font-size:small;background-color:#FF3300;width:100%;}
	.info 	 {font-family:"Lucida Console", Monaco, monospace;font-size:small;width:100%;}
	.suspend {font-family:"Lucida Console", Monaco, monospace;font-size:small;}`)
	return
}
