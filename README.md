# Hopwatch, a debugging tool for Go

Hopwatch is an experimental tool in HTML5 that can help debugging Go programs. 
It works by communicating to a WebSockets based client in Javascript.
When your program calls the Break function, it sends debug information to the browser page and waits for user interaction.
On the hopwatch page, the developer can view debug information and choose to resume the execution of the program.

###Documentation

[http://go.pkgdoc.org/github.com/emicklei/hopwatch](http://go.pkgdoc.org/github.com/emicklei/hopwatch)


&copy; 2012, http://ernestmicklei.com. MIT License
	