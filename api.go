package hopwatch

import "log"

// Printf formats according to a format specifier and writes to the debugger screen.
// It returns a new Watchpoint to send more or break.
func Printf(format string, params ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Printf(format, params...)
}

// Display sends variable name,value pairs to the debugger.
// The parameter nameValuePairs must be even sized.
func Display(nameValuePairs ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Display(nameValuePairs...)
}

// Break suspends the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func Break(conditions ...bool) {
	suspend(2, conditions...)
}

// CallerOffset (default=2) allows you to change the file indicator in hopwatch.
// Use this method when you wrap the .CallerOffset(..).Display(..).Break() in your own function.
func CallerOffset(offset int) *Watchpoint {
	return (&Watchpoint{}).CallerOffset(offset)
}

func Disable() {
	log.Print("[hopwatch] disabled by code.\n")
	hopwatchEnabled = false
}

func Enable() {
	log.Print("[hopwatch] enabled by code.\n")
	hopwatchEnabled = true
}
