socket.io-go-emitter
========================

A [Go](https://golang.org/) implementation of socket.io-emitter

[socket.io](http://socket.io/) provides a hook point to easily allow you to emit events to browsers from anywhere so `socket.io-go-emitter` communicates with [socket.io](http://socket.io/) servers through redis.


## How to use

**Install**

```go
go get github.com/ziyasal/socket.io-go-emitter/emitter
```

## API

### Emitter(opts)

The following options are allowed:
- `key`: the name of the key to pub/sub events on as prefix (`socket.io`)
- `host`: host to connect to redis on (`localhost`)
- `port`: port to connect to redis on (`6379`)

**Important** Make sure to supply the`host` and `port` options.

Specifies a specific `room` that you want to emit to.

**Initialize emitter**
```go
import "github.com/ziyasal/socket.io-go-emitter/emitter"

//....

opts := map[string]string{
	"host" : "127.0.0.1",
	"port" : "6379",
	"key"  : "socket.io",
 }
sio := emitter.New((map[string]string)(opts))
```

###Emitter#Emit(channel,message):Emitter
```go
  sio.Emit("broadcast event", "Hello from socket.io-go-emitter")
```


### Emitter#In(room):Emitter
```go
  sio.In("test").Emit("broadcast event", "Hello from socket.io-go-emitter")
```
### Emitter#To(room):Emitter
```go
 sio.To("test").Emit("broadcast event", "Hello from socket.io-go-emitter")
```

### Emitter#Of(namespace):Emitter
Specifies a specific namespace that you want to emit to.
```go
 sio.Of("/nsp").Emit("broadcast event", "Hello from socket.io-go-emitter")
```

##Bugs
If you encounter a bug, performance issue, or malfunction, please add an [Issue](https://github.com/ziyasal/socket.io-go-emitter/issues) with steps on how to reproduce the problem.

##TODO
- Add more tests
- Add samples

###Open Source Projects in Use
* [redigo](https://github.com/garyburd/redigo) by Gary Burd @garyburd
* [msgpack](https://github.com/vmihailenco/msgpack) by Vladimir Mihailenco @vmihailenco

##License
Code and documentation are available according to the *MIT* License (see [LICENSE](https://github.com/ziyasal/socket.io-go-emitter/blob/master/LICENSE)).


