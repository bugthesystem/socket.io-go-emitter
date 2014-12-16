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

type EmitterOptions struct {
	Host           string
	Password       string
	Key            string
}

type Emitter struct {
	_opts    EmitterOptions
	_key     string
	_flags   map[string]string
	_rooms   map[string]bool
}

func New(opts EmitterOptions) Emitter {
	validateOptions(opts)

	emitter := Emitter{_opts:opts}

	if opts.Key != "" {
		emitter._key = opts.Key+"#emitter"
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
	//Pack & Publish
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

func validateOptions(opts EmitterOptions) () {
	if opts.Host == "" {
		panic("Missing redis `host`")
	}
}

func dial(opts EmitterOptions) (redis.Conn, error) {
	c, err := redis.DialTimeout("tcp", opts.Host, 0, 10*time.Second, 10*time.Second)
	if err != nil {
		return nil, err
	}

	if opts.Password != "" {
		if _, err := c.Do("AUTH", opts.Password); err != nil {
			c.Close()
			return nil, err
		}
		return c, err
	}

	return c, nil
}


func publish(opts EmitterOptions, channel, value interface{}) {
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