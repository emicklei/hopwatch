# Hopwatch, a debugging tool for Go

Hopwatch is a simple tool in HTML5 that can help debug Go programs. 
It works by communicating to a WebSockets based client in Javascript.
When your program calls the Break function, it sends debug information to the browser page and waits for user interaction.
Using the functions Display, Printf or Dump (go-spew), you can log information on the browser page.
On the hopwatch page, the developer can view debug information and choose to resume the execution of the program.

[First announcement](https://ernestmicklei.com/2012/12/hopwatch-a-debugging-tool-for-go/)

![How](hopwatch_how.png)


## Distributed (work in progress)

Hopwatch can be used in a distributed services architecture to hop and trace between services following an incoming request to downstream service.

Consider the setup where the browser is sending a HTTP request to a GraphQL endpoint which calls a gRPC backend service, which calls a PostgreSQL Database server to perform a query. The result of that query needs to be transformed into a gRPC response which in turn needs to be transformed into a GraphQL response before transporting it back to the browser.

We want to jump from client to server to server and back, for a given request. 
To signal the upstream services that it should break on this request, the request must be annotated using a special HTTP header `x-hopwatch : your-correlation-name`.

Each upstream server must have the hopwatch agent package included:

    import _ "github.com/emicklei/hopwatch/agent"


&copy; 2012-2022, http://ernestmicklei.com. MIT License 