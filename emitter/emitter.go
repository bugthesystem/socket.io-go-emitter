package emitter

import (
	"fmt"
	"io"
	"reflect"
	"github.com/ugorji/go/codec"
	//"github.com/garyburd/redigo/redis"
)

const (
	event       = 2
	binaryEvent = 5
)

type Emitter struct {
	_key     string
	_client  string
	_flags  map[string]string
	_rooms   []string
}

func New(opts map[string]string) Emitter {
	emitter := Emitter{}

	if name, ok := opts["client"]; ok {
		emitter._client = name
		fmt.Println(name, ok)
	}

	if name, ok := opts["key"]; ok {
		emitter._key = name
		fmt.Println(name, ok)
	}

	emitter._rooms = make([]string, 0)
	emitter._flags = make(map[string]string)

	return emitter
}

func (emitter Emitter) In(room string) Emitter {
	emitter._rooms = append(emitter._rooms, room)
	return emitter
}

func (emitter Emitter) Emit(args ...interface{}) bool {

	fmt.Println("Args : ", args)
	fmt.Println("Emitter key: ", emitter._key)
	fmt.Println("Emitter cleint: ", emitter._client)
	fmt.Println("Emitter flags: ", emitter._flags)
	fmt.Println("Emitter rooms: ", emitter._rooms)

	// create and configure Handle
	var (
		mh codec.MsgpackHandle
	)
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	// create and use decoder/encoder
	var (
		w io.Writer
		b []byte
		h = &mh
	)

	var enc = codec.NewEncoder(w, h)
	enc = codec.NewEncoderBytes(&b, h)
	var err = enc.Encode(args)
	err = enc.Encode(args)

	fmt.Println("err : ", err)
	fmt.Println("mh : ", mh)
	fmt.Println("b : ", b)

	emitter._rooms = make([]string, 0)
	emitter._flags = make(map[string]string)

	return true
}
