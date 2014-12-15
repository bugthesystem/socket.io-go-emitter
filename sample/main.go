package main

import (
	"fmt"
	"github.com/ziyasal/socket.io-go-emitter/emitter"
)

func main() {
	opts := map[string]string{
		"client":"client",
		"key"   : "key",
	}
	sio := emitter.New((map[string]string)(opts))

	fmt.Println("Emit Result: ", sio.In("test").Emit("broadcast event", "Hello from socket.io-go-emitter"))
}
