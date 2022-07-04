package hopwatch

import (
	"fmt"
	"log"
	"runtime"
)

// watchpoint is a helper to provide a fluent style api.
// This allows for statements like hopwatch.Display("var",value).Break()
type Watchpoint struct {
	disabled bool
	offset   int // offset in the caller stack for highlighting source
}

// CallerOffset (default=2) allows you to change the file indicator in hopwatch.
func (w *Watchpoint) CallerOffset(offset int) *Watchpoint {
	if hopwatchEnabled && (offset < 0) {
		log.Panicf("[hopwatch] ERROR: illegal caller offset:%v . watchpoint is disabled.\n", offset)
		w.disabled = true
	}
	w.offset = offset
	return w
}

// Printf formats according to a format specifier and writes to the debugger screen.
func (w *Watchpoint) Printf(format string, params ...interface{}) *Watchpoint {
	w.offset += 1
	var content string
	if len(params) == 0 {
		content = format
	} else {
		content = fmt.Sprintf(format, params...)
	}
	return w.printcontent(content)
}

// Printf formats according to a format specifier and writes to the debugger screen.
func (w *Watchpoint) printcontent(content string) *Watchpoint {
	_, file, line, ok := runtime.Caller(w.offset)
	cmd := command{Action: "print"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	cmd.addParam("line", content)
	channelExchangeCommands(cmd)
	return w
}

// Display sends variable name,value pairs to the debugger. Values are formatted using %#v.
// The parameter nameValuePairs must be even sized.
func (w *Watchpoint) Display(nameValuePairs ...interface{}) *Watchpoint {
	_, file, line, ok := runtime.Caller(w.offset)
	cmd := command{Action: "display"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	if len(nameValuePairs)%2 == 0 {
		for i := 0; i < len(nameValuePairs); i += 2 {
			k := nameValuePairs[i]
			v := nameValuePairs[i+1]
			cmd.addParam(fmt.Sprint(k), fmt.Sprintf("%#v", v))
		}
	} else {
		log.Printf("[hopwatch] WARN: missing variable for Display(...) in: %v:%v\n", file, line)
		w.disabled = true
		return w
	}
	channelExchangeCommands(cmd)
	return w
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func (w Watchpoint) Break(conditions ...bool) {
	suspend(w.offset, conditions...)
}
