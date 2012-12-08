package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	for i := 0; i < 6; i++ {
		hopwatch.Watch("i",i).Break()
		j := i * i
		hopwatch.Watch("j",j).BreakIf(j > 10)
	}
}
