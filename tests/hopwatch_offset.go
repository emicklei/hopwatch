package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	hopwatch.Display("8",8)
	hopwatch.Display("9",9).Break()
	inside()
	indirectDisplay("11",11)
	indirectBreak()
}
func inside() {
	hopwatch.Display("15",15)
	hopwatch.Display("16",16).Break()
}
func indirectDisplay(args ...interface{}) {
	hopwatch.CallerOffset(2).Display(args...)
}
func indirectBreak() {
	hopwatch.CallerOffset(3).Break()
}
