// Copyright 2012,2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
)

// Dump displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// Delegates to spew.Fdump, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func Dump(a ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Dump(a...)
}

// Dumpf formats and displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// delegates to spew.Fprintf, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func Dumpf(format string, a ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Dumpf(format, a...)
}

// Dump displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// Delegates to spew.Fdump, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func (w *Watchpoint) Dump(a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	spew.Fdump(writer, a...)
	return w.printcontent(string(writer.Bytes()))
}

// Dumpf formats and displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// Delegates to spew.Fprintf, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func (w *Watchpoint) Dumpf(format string, a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	_, err := spew.Fprintf(writer, format, a...)
	if err != nil {
		return Printf("[hopwatch] error in spew.Fprintf:%v", err)
	}
	return w.printcontent(string(writer.Bytes()))
}
