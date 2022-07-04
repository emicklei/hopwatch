package main

import (
	"io"
	"net/http"

	"github.com/emicklei/hopwatch"
)

func main() {
	hopwatch.Break()
	resp, _ := http.Get("http://localhost:8998/hop")
	data, _ := io.ReadAll(resp.Body)
	hopwatch.Display("response", string(data)).Break()
}
