package hopwatch

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
)

// Dump delegates to spew.Fdump, see https://github.com/davecgh/go-spew
func Dump(a ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 3}
	return wp.Dump(a...)
}

// Dumpf delegates to spew.Fprintf, see https://github.com/davecgh/go-spew
func Dumpf(format string, a ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 3}
	return wp.Dumpf(format, a...)
}

// Dump delegates to spew.Fdump, see https://github.com/davecgh/go-spew
func (self *Watchpoint) Dump(a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	spew.Fdump(writer, a...)
	return self.printcontent(string(writer.Bytes()))
}

// Dumpf delegates to spew.Fprintf, see https://github.com/davecgh/go-spew
func (self *Watchpoint) Dumpf(format string, a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	_, err := spew.Fprintf(writer, format, a...)
	if err != nil {
		return Printf("[hopwatch] error in spew.Fprintf:%v", err)
	}
	return self.printcontent(string(writer.Bytes()))
}
