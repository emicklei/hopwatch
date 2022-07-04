package hopwatch

// command is used to transport message to and from the debugger.
type command struct {
	Action     string
	Parameters map[string]string
}

// addParam adds a key,value string pair to the command ; no check on overwrites.
func (c *command) addParam(key, value string) {
	if c.Parameters == nil {
		c.Parameters = map[string]string{}
	}
	c.Parameters[key] = value
}
