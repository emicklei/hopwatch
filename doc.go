/*
Hopwatch is an experimental HTML5 application and Go package that can help debugging Go programs. 

It works by communicating to a WebSockets based agent in Javascript.
When your program calls the Break function, it sends debug information to the browser page and waits for user interaction.
On the hopwatch page, the developer can view debug information and choose to proceed or terminate the execution of your program.

Usage:

	import (
		"github.com/emicklei/hopwatch"
	)

	func foo() {
		bar := "john"
		// suspends execution until hitting "Resume" in the browser
		hopwatch.Display("foo", bar).Break()
	}

Tool:

Open the Hopwatch debugger on http://localhost:23456/hopwatch.html.
Your browser must support WebSockets ; it has been tested with Chrome.

Resources:

[project]	https://github.com/emicklei/hopwatch

[blog]		http://ernestmicklei.com

(c) 2012, http://ernestmicklei.com. MIT License
*/
package hopwatch
