package hopwatch

import (
	"log"
	"testing"
)

func TestWatchpoint_Caller(t *testing.T) {
	go shortCircuit(commandResume())
	Caller().Break()
}

func commandResume() command {
	return command{Action: "resume"}
}

func shortCircuit(next command) {
	cmd := <-toBrowserChannel
	log.Printf("send to browser:%#v\n", cmd)
	log.Printf("received from browser:%#v\n", next)
	fromBrowserChannel <- next
}
