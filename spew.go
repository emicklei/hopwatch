package hopwatch

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
)

// Dump delegates to spew.Fdump, see https://github.com/davecgh/go-spew
func Dump(a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	spew.Fdump(writer, a)
	wp := &Watchpoint{offset: 2}
	return wp.printcontent(string(writer.Bytes()))
}

// Dumpf delegates to spew.Fprintf, see https://github.com/davecgh/go-spew
func Dumpf(format string, a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	_, err := spew.Fprintf(writer, format, a)
	if err != nil {
		return Printf("[hopwatch] error in spew.Fprintf:%v", err)
	}
	wp := &Watchpoint{offset: 2}
	return wp.printcontent(string(writer.Bytes()))
}
