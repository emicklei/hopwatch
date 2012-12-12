package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	hopwatch.Display("8",8)
	hopwatch.Display("9",9).Break()
	inside()
}
func inside() {
	hopwatch.Display("13",14)
	hopwatch.Display("14",14).Break()
}
