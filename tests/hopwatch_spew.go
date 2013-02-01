package main

import (
	"github.com/emicklei/hopwatch"
)

type node struct {
	label string
	parent *node
	children []node
}

func main() {
	tree := node{label:"parent", children:[]node{node{label:"child"}}}
	
	hopwatch.Spew().Printf("spew %v","it")	
	hopwatch.Spew().Dump(tree)
}