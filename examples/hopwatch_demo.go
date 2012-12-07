package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	for i := 0; i < 10; i++ {
		hopwatch.Break("hopwatch_demo.go:main", "i", i)
	}
}
