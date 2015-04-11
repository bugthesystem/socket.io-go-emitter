package emitter

import (
	"fmt"
	"time"
	"gopkg.in/vmihailenco/msgpack.v1"
	"github.com/garyburd/redigo/redis"
)

const (
	event              = 2
	binaryEvent        = 5
	redisPoolMaxIdle   = 80
	redisPoolMaxActive = 12000 // max number of connections
)

type EmitterOptions struct {
	Host       string
	Password   string
	Key        string
}

type Emitter struct {
	_opts      EmitterOptions
	_key       string
	_flags     map[string]string
	_rooms     map[string]bool
	_pool      *redis.Pool
}

func New(opts EmitterOptions) Emitter {
	emitter := Emitter{_opts:opts}

	initRedisConnPool(&emitter, opts)

	if opts.Key != "" {
		emitter._key = fmt.Sprintf("%s#emitter", opts.Key)
	}else {
		emitter._key = "socket.io#emitter"
	}

	emitter._rooms = make(map[string]bool)
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

	//TODO: Goroutines
	//Pack & Publish
	b, err := msgpack.Marshal([]interface{}{packet, extras})
	if err != nil {
		panic(err)
	}else {
		publish(emitter, emitter._key, b)
	}

	emitter._rooms = make(map[string]bool)
	emitter._flags = make(map[string]string)

	return true
}


func (emitter Emitter) hasBin(args ...interface{}) bool {
	//NOT implemented yet!
	return true
}

func initRedisConnPool(emitter *Emitter , opts EmitterOptions) () {
	if opts.Host == "" {
		panic("Missing redis `host`")
	}

	emitter._pool = newPool(opts)
}

func newPool(opts EmitterOptions) *redis.Pool {
	return &redis.Pool{
		MaxIdle: redisPoolMaxIdle,
		MaxActive: redisPoolMaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", opts.Host)
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

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}


func publish(emitter Emitter, channel string, value interface{}) {
	c := emitter._pool.Get()
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
