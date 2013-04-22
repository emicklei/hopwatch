package main

import (
	"github.com/emicklei/hopwatch"
)

var ENABLED = true

func watch(doBreak bool, vars ...interface{}) {
	if !ENABLED {
		return
	}
	watchPoint := hopwatch.CallerOffset(2 + 1) // compensate for indirect call
	if doBreak {
		watchPoint.Dump(vars...).Break()
	}
}
