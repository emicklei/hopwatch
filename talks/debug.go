package main

import (
	"github.com/emicklei/hopwatch"
)

var DEBUG = true

func debug(doBreak bool, vars ...interface{}) {
	if !DEBUG {
		return
	}
	watchPoint := hopwatch.CallerOffset(3) // compensate for stack
	if doBreak {
		watchPoint.Break()
	}
}
