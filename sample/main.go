package main

import (
	"fmt"
	"github.com/ziyasal/socket.io-go-emitter/emitter"
)

func main() {
	opts := map[string]string{
		"host" : "127.0.0.1",
		"port" : "6379",
		"key"  : "socket.io",
	}
	sio := emitter.New((map[string]string)(opts))

	fmt.Println("Emit :", sio.Emit("broadcast event", "Hello from socket.io-go-emitter"))

	fmt.Println("In Emit :", sio.In("test").Emit("broadcast event", "Hello from socket.io-go-emitter"))

	fmt.Println("To Emit : ", sio.To("test").Emit("broadcast event", "Hello from socket.io-go-emitter"))

	fmt.Println("Of Emit :", sio.Of("/nsp").Emit("broadcast event", "Hello from socket.io-go-emitter"))
}
