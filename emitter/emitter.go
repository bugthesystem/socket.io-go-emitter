package emitter

import (
	"fmt"
	"time"
	"gopkg.in/vmihailenco/msgpack.v1"
	"github.com/garyburd/redigo/redis"
)

const (
	event       = 2
	binaryEvent = 5
)

type Emitter struct {
	_opts    map[string]string
	_key     string
	_flags   map[string]string
	_rooms   map[string]bool
}

func New(opts map[string]string) Emitter {
	validateOptions(opts)

	emitter := Emitter{_opts:opts}

	if value, ok := opts["key"]; ok {
		emitter._key = value+"#emitter"
	}else {
		emitter._key = "socket.io#emitter"
	}

	emitter._rooms = make(map[string]bool, 0)
	emitter._flags = make(map[string]string)

	return emitter
}

func (emitter Emitter) In(room string) Emitter {
	if _, ok := emitter._rooms[room]; ok == false {
		emitter._rooms[room] = true
	}
	return emitter
}

func (emitter Emitter) To(room string) Emitter {
	return emitter.In(room)
}

func (emitter Emitter) Of(nsp string) Emitter {
	emitter._flags["nsp"] = nsp
	return emitter
}

func (emitter Emitter) Emit(args ...interface{}) bool {

	packet := make(map[string]interface{})
	extras := make(map[string]interface{})

	if ok := emitter.hasBin(args); ok {
		packet["type"] = binaryEvent
	}else {
		packet["type"] = event
	}

	packet["data"] = args

	if value, ok := emitter._flags["nsp"]; ok {
		packet["nsp"] = value
		delete(emitter._flags, "nsp")
	}else {
		packet["nsp"] = "/"
	}


	if ok := len(emitter._rooms); ok > 0 {
		//TODO:Cast??
		extras["rooms"] = getKeys(emitter._rooms)
	}else {
		extras["rooms"] = make([]string, 0, 0)
	}

	if ok := len(emitter._flags); ok > 0 {
		extras["flags"] = emitter._flags
	}else {
		extras["flags"] = make(map[string]string)
	}

	//TODO: Gorotunes
	//PAck & Publish
	b, err := msgpack.Marshal([]interface{}{packet, extras})
	if err != nil {
		panic(err)
	}else {
		publish(emitter._opts, emitter._key, b)
	}

	emitter._rooms = make(map[string]bool)
	emitter._flags = make(map[string]string)

	return true
}


func (emitter Emitter) hasBin(args ...interface{}) bool {
	//NOT implemented yet!
	return true
}

func validateOptions(opts map[string]string) () {
	if _, ok := opts["host"]; ok == false {
		panic("Missing redis `host`")
	}

	if _, ok := opts["port"]; ok == false {
		panic("Missing redis `port`")
	}
}

func dial(opts map[string]string) (redis.Conn, error) {
	connStr := fmt.Sprintf("%s:%s", opts["host"], opts["port"])
	c, err := redis.DialTimeout("tcp", connStr, 0, 10*time.Second, 10*time.Second)
	if err != nil {
		return nil, err
	}

	if value, ok := opts["password"]; ok {
		if _, err := c.Do("AUTH", value); err != nil {
			c.Close()
			return nil, err
		}
		return c, err
	}

	return c, nil
}


func publish(opts map[string]string, channel, value interface{}) {
	c, err := dial(opts)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.Do("PUBLISH", channel, value)

}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}


func dump(emitter Emitter, args ...interface{}) {
	fmt.Println("Emit params : ", args)
	fmt.Println("Emitter key: ", emitter._key)
	fmt.Println("Emitter flags: ", emitter._flags)
	fmt.Println("Emitter rooms: ", emitter._rooms)
}