// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

Connect:

	Open the Hopwatch debugger on http://localhost:23456/hopwatch.html after starting your program.
	Your browser must support WebSockets ; it has been tested with Chrome.

Other examples:

	hopwatch.Break(i > 10,  j < 100)	// zero or more conditions ; suspends program (or goroutine)
	hopwatch.Display("i",i , "j",j")	// zero or more name,value pairs ; no program suspend
	hopwatch.Caller().Display("a",a)	// fixes the caller offset when called inside a wrapper function

The flags are:

	-hopwatch	if present and set to false then hopwatch is disabled and will not connect to the debugger.

Resources:

	https://github.com/emicklei/hopwatch

(c) 2012, http://ernestmicklei.com. MIT License
*/
package hopwatch
