package hopwatch

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
)

// Spewpoint provides the integration with go-spew from Dave Collins <dave@davec.name>
// https://github.com/davecgh/go-spew
type Spewpoint struct {
}

func Spew() *Spewpoint {
	return &Spewpoint{}
}

// Printf delegates to spew.Fprintf, see https://github.com/davecgh/go-spew
func (sp *Spewpoint) Printf(format string, a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	_, err := spew.Fprintf(writer, format, a)
	if err != nil {
		return Printf("[hopwatch] error in spew.Fprintf:%v", err)
	}
	wp := &Watchpoint{offset: 2}
	return wp.printcontent(string(writer.Bytes()))
}

// Dump delegates to spew.Fdump, see https://github.com/davecgh/go-spew
func (sp *Spewpoint) Dump(a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	spew.Fdump(writer, a)
	wp := &Watchpoint{offset: 2}
	return wp.printcontent(string(writer.Bytes()))
}
