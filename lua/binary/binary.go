package binary

import (
	"bytes"
	"fmt"
	"io"
)

var (
	head = [...]byte{0x1B, 0x4C, 0x75, 0x61}
	tail = [...]byte{0x19, 0x93, '\r', '\n', 0x1A, '\n'}
)

const (
	LUA_SIGNATURE 	 = "\x1bLua"
	LUAC_VERSION  	 = 0x53
	LUAC_FORMAT   	 = 0
	LUAC_DATA 	  	 = "\x19\x93\r\n\x1a\n"
	CINT_SIZE 	  	 = 4
	CSIZET_SIZE   	 = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT 		 = 0x5678
	LUAC_NUM 		 = 370.5
)

const (
	LUA_TYPE_NONE = iota - 1
	LUA_TYPE_NIL
	LUA_TYPE_BOOL
	LUA_TYPE_LUDATA
	LUA_TYPE_NUMBER
	LUA_TYPE_STRING
	LUA_TYPE_TABLE
	LUA_TYPE_FUNC
	LUA_TYPE_UDATA
	LUA_TYPE_THREAD
)

// lua-5.3.4/src/lobject.h
const (
	LUA_NUM_FLOAT  = LUA_TYPE_NUMBER | (0 << 4) // floating-point numbers
	LUA_NUM_INT    = LUA_TYPE_NUMBER | (1 << 4) // integer numbers
	LUA_STR_SHORT  = LUA_TYPE_STRING | (0 << 4) // short strings
	LUA_STR_LONG   = LUA_TYPE_STRING | (1 << 4) // long strings
	LUA_CLOSURE    = LUA_TYPE_FUNC | (0 << 4)   // lua closure
	LUA_GO_FUNC    = LUA_TYPE_FUNC | (1 << 4)   // light go func
	LUA_GO_CLOSURE = LUA_TYPE_FUNC | (2 << 4)   // go closure
)

type (
	Prototype struct {
		Source   string
		SrcPos   uint32
		EndPos   uint32
		Params   byte
		Vararg   byte
		Stack    byte
		Code 	 []uint32
		Consts   []interface{}
		UpValues []UpValue
		Protos   []Prototype
		PcLnTab  []uint32
		Locals   []LocalVar
		UpNames  []string
	}

	Header struct {
		Signature  [4]byte
		Version    byte
		Format     byte
		LuacData   [6]byte
		GoIntSize  byte
		SizetSize  byte
		InstrSize  byte
		LuaIntSize byte
		LuaNumSize byte
		LuacIntEnc int64
		LuacNumEnc float64
	}

	LocalVar struct {
		Name string
		Live uint32
		Dead uint32
	}

	UpValue struct {
		InStack byte
		Index   byte
	}

	Chunk struct {
		Header Header
		Entry  Prototype
	}
)

func (proto *Prototype) NumParams() int { return int(proto.Params) }
func (proto *Prototype) StackSize() int { return int(proto.Stack) }
func (proto *Prototype) IsVararg() bool { return int(proto.Vararg) == 1 }
func (proto *Prototype) Const(index int) interface{} { return proto.Consts[index] }
func (proto *Prototype) Proto(index int) *Prototype { return &proto.Protos[index] }

func (upval *UpValue) IsLocal() bool { return upval.InStack == 1 }
func (upval *UpValue) AtIndex() int { return int(upval.Index) }

func IsChunk(data []byte) bool { return len(data) > 4 && string(data[:4]) == LUA_SIGNATURE }

func Load(data []byte) (chunk Chunk, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				if e != io.EOF {
					err = e
					return
				}
			}
		}
	}()
	decode(bytes.NewBuffer(data), &chunk)
	return chunk, err
}

func Dump(proto *Prototype, strip bool) []byte {
	var b bytes.Buffer
	w := &writer{b: &b}
	w.writeHeader()
	n := len(proto.UpValues)
	w.writeByte(byte(n))
	encodeProto(w, proto)
	return b.Bytes()
}

func assert(cond bool, mesg string) {
	if !cond {
		panic(fmt.Errorf(mesg))
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
