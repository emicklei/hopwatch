package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	for i := 0; i < 10; i++ {
		j := i * i
		hopwatch.Watch("i",i).Watch("j",j).Break()
	}
}
