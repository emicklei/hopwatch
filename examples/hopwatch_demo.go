package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	for i := 0; i < 6; i++ {
		hopwatch.Display("i",i)
		j := i * i
		hopwatch.Display("i",i, "j", j ).Break(j > 10)
		hopwatch.Break()
	}
}
