# hopping between running applications with 2 debuggers

Open 2 terminal sessions.
In the first, you start the server:

    cd server
    go run *.go

In the second, you start the client which will hit a breakpoint.

    cd client
    go run *.go -hopwatch.port=23455

Resuming that breakpoint will hit the breakpoint in the server.
Resuming that breakpoint will hit another breakpoint in the client.
