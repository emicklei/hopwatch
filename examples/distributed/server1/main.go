package main

import (
	"io"
	"log"
	"net/http"

	"github.com/emicklei/hopwatch/agent"
)

func main() {
	http.HandleFunc("/hop", handleRequest) // dont listen root, browsers want icons so bad
	log.Println("listing for HTTP on http://localhost:8998/hop")
	http.ListenAndServe(":8998", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	agent.Display("request", r).Break()
	io.WriteString(w, "hello hopper")
}
