package code

import (
	"encoding/binary"
	"bytes"
	"fmt"
)

const (
	// maximum length for short strings, that is, strings that are internalized.
	//
	// Cannot be smaller than reserved words or tags for metamethods, as these
	// strings must be internalized; #("function") = 8, #("__newindex") = 10.
	maxShortLen = 40

	// signature is the mark for precompiled code ('<esc>Lua')
	signature = "\x1bLua"
)

var order = binary.LittleEndian

var (
	head = [...]byte{0x1B, 0x4C, 0x75, 0x61}
	tail = [...]byte{0x19, 0x93, '\r', '\n', 0x1A, '\n'}
)

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
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

type header struct {
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

type source struct {
	ord   binary.ByteOrder
	src   *bytes.Buffer
	name  string
	strip bool
}

func (bin *source) write(data interface{}) error {
	return binary.Write(bin.src, bin.ord, data)
}

func (bin *source) writeHeader() error {
	h := &header{
		Version:    LUAC_VERSION,
		Format:     LUAC_FORMAT,
		GoIntSize:  CINT_SIZE,
		SizetSize:  CSIZET_SIZE,
		InstrSize:  INSTRUCTION_SIZE,
		LuaIntSize: LUA_INTEGER_SIZE,
		LuaNumSize: LUA_NUMBER_SIZE,
		LuacIntEnc: LUAC_INT,
		LuacNumEnc: LUAC_NUM,
	}
	copy(h.Signature[:], LUA_SIGNATURE)
	copy(h.LuacData[:], LUAC_DATA)
	return bin.write(h)
}

func (bin *source) writeProto(fn *Proto, srcID string) {
	if bin.strip || fn.Source == srcID {
		bin.writeStr("")
	} else {
		bin.writeStr(fn.Source)
	}
	bin.write(uint32(fn.SrcPos))
	bin.write(uint32(fn.EndPos))
	bin.write(byte(fn.ParamN))
	bin.write(byte(b2i(fn.Vararg)))
	bin.write(byte(fn.StackN))

	bin.writeInstrs(fn)
	bin.writeConsts(fn)
	bin.writeUpVars(fn)
	bin.writeProtos(fn)
	bin.writeDebug(fn)
}

func (bin *source) writeInstrs(fn *Proto) {
	bin.write(uint32(len(fn.Instrs)))
	for _, inst := range fn.Instrs {
		bin.write(inst)
	}
}

func (bin *source) writeConsts(fn *Proto) {
	bin.write(uint32(len(fn.Consts)))
	for _, kst := range fn.Consts {
		switch kst := kst.(type) {
			case string:
				if len(kst) <= maxShortLen { // short string?
					bin.write(byte(LUA_STR_SHORT))
				} else {
					bin.write(byte(LUA_STR_LONG))
				}
				bin.writeStr(string(kst))
			case float64:
				bin.write(byte(LUA_NUM_FLOAT))
				bin.write(float64(kst))
			case bool:
				bin.write(byte(LUA_TYPE_BOOL))
				bin.write(byte(b2i(bool(kst))))
			case int64:
				bin.write(byte(LUA_NUM_INT))
				bin.write(int64(kst))
			default:
				if kst == nil {
					bin.write(byte(LUA_TYPE_NIL))
				} else {
					panic(fmt.Errorf("unexpected constant %T", kst))
				}
		}
	}
}

func (bin *source) writeUpVars(fn *Proto) {
	bin.write(uint32(len(fn.UpVars)))
	for _, up := range fn.UpVars {
		bin.write(byte(b2i(up.Stack)))
		bin.write(byte(up.Index))
	}
}

func (bin *source) writeProtos(fn *Proto) {
	bin.write(uint32(len(fn.Protos)))
	for _, fn := range fn.Protos {
		bin.writeProto(fn, "")
	}
}

func (bin *source) writeDebug(fn *Proto) {
	var (
		pclines = len(fn.PcLine)
		localsN = len(fn.Locals)
		upvarsN = len(fn.UpVars)
	)
	if bin.strip {
		pclines = 0
		localsN = 0
		upvarsN = 0
	}
	bin.write(uint32(pclines))
	for i := 0; i < pclines; i++ {
		bin.write(uint32(fn.PcLine[i]))
	}
	bin.write(uint32(localsN))
	for i := 0; i < localsN; i++ {
		bin.writeStr(fn.Locals[i].Name)
		bin.write(uint32(fn.Locals[i].Live))
		bin.write(uint32(fn.Locals[i].Dead))
	}
	bin.write(uint32(upvarsN))
	for i := 0; i < upvarsN; i++ {
		bin.writeStr(fn.UpVars[i].Name)
	}
}

func (bin *source) writeStr(s string) {
	if s == "" {
		bin.write(byte(1))
	} else {
		if size := len(s)+1; size < 0xFF {
			bin.write(byte(size))
		} else {
			bin.write(byte(0xFF))
			bin.write(uint64(size))
		}
		bin.write([]byte(s))
	}
}

func checkHeader(buf *bytes.Buffer) {
	var h header
	must(binary.Read(buf, order, &h))
	check(h.Signature == head, "not a precompiled chunk")
	check(h.Version == LUAC_VERSION, "version mismatch")
	check(h.Format == LUAC_FORMAT, "format mismatch")
	check(h.LuacData == tail, "corrupted")
	check(h.GoIntSize == CINT_SIZE, "int size mismatch")
	check(h.SizetSize == CSIZET_SIZE, "size_t size mismatch")
	check(h.InstrSize == INSTRUCTION_SIZE, "instruction size mismatch")
	check(h.LuaIntSize == LUA_INTEGER_SIZE, "lua integer size mismatch")
	check(h.LuaNumSize == LUA_NUMBER_SIZE, "lua number size mismatch")
	check(h.LuacIntEnc == LUAC_INT, "endianess mismatch")
	check(h.LuacNumEnc == LUAC_NUM, "float format mismatch")
}

func decodeChunk(r *bytes.Buffer, main *Proto, src string) {
	checkHeader(r)

	// decode size_upvalues (?)
	_, err := r.ReadByte()
	must(err)

	// decode main prototype
	decodeProto(r, main, src)
}

func decodeProto(r *bytes.Buffer, p *Proto, src string) {
	// decode source name string (b[0] == length)
	// b[0] == 0xFF ? uint64 : size
	if p.Source = decodeString(r); p.Source == "" || p.Source == "=stdin" {
		p.Source = src
	}

	// decode line start
	{
		var u32 uint32
		must(binary.Read(r, order, &u32))
		p.SrcPos = int(u32)
	}

	// decode line end
	{
		var u32 uint32
		must(binary.Read(r, order, &u32))
		p.EndPos = int(u32)
	}

	// decode number of parameters
	{
		var u8 byte
		must(binary.Read(r, order, &u8))
		p.ParamN = int(u8)
	}

	// decode is varadic
	{
		var u8 byte
		must(binary.Read(r, order, &u8))
		p.Vararg = (u8 == 1)
	}

	// decode maximum stack size
	{
		var u8 byte
		must(binary.Read(r, order, &u8))
		p.StackN = int(u8)
	}

	// decode instruction bytecode
	//
	// leading 4-bytes is number of instructions
	{
		var num uint32
		must(binary.Read(r, order, &num))
		p.Instrs = make([]Instr, num)
		for i := range p.Instrs {
			var u32 uint32
			must(binary.Read(r, order, &u32))
			p.Instrs[i] = Instr(u32)
		}
	}

	// decode constants
	//
	// leading 4-bytes is number of constants
	{
		var num uint32
		must(binary.Read(r, order, &num))
		p.Consts = make([]Const, num)
		for i := range p.Consts {
			p.Consts[i] = decodeConst(r)
		}
	}

	// decode upvalues
	//
	// leading 4-bytes is number of upvalues
	{
		var num uint32
		must(binary.Read(r, order, &num))
		p.UpVars = make([]*UpVar, num)
		for i := range p.UpVars {
			var up UpVar
			// decode upvalue instack
			{
				var u8 byte
				must(binary.Read(r, order, &u8))
				up.Stack = (u8 == 1)
			}
			// decode upvalue index
			{
				var u8 byte
				must(binary.Read(r, order, &u8))
				up.Index = int(u8)
			}
			p.UpVars[i] = &up
	
		}
	}

	// decode nested closure prototypes
	//
	// leading 4-bytes is number of prototypes
	{
		var num uint32
		must(binary.Read(r, order, &num))
		p.Protos = make([]*Proto, num)
		for i := range p.Protos {
			var fn Proto
			decodeProto(r, &fn, src)
			p.Protos[i] = &fn
		}
	}

	// decode line info (pc -> line)
	//
	// leading 4-bytes is number of pcln entries
	{
		var num uint32
		must(binary.Read(r, order, &num))
		p.PcLine = make([]int32, num)
		for i := range p.PcLine {
			must(binary.Read(r, order, &p.PcLine[i]))
		}
	}

	// decode local variables
	//
	// leading 4-bytes is number of local entries
	{
		var num uint32
		must(binary.Read(r, order, &num))
		p.Locals = make([]*Local, num)
		for i := range p.Locals {
			local := &Local{Name: decodeString(r)}
			must(binary.Read(r, order, &local.Live))
			must(binary.Read(r, order, &local.Dead))
			p.Locals[i] = local
		}
	}

	// decode upvalue names
	//
	// leading 4-bytes is number of name entries
	{
		var num uint32
		must(binary.Read(r, order, &num))
		for i := range p.UpVars {
			p.UpVars[i].Name = decodeString(r)
		}
	}
}

func decodeConst(r *bytes.Buffer) Const {
	t, err := r.ReadByte()
	must(err)
	switch t {
		case LUA_TYPE_NIL: // NIL
			return nil

		case LUA_TYPE_BOOL: // BOOL
			b, err := r.ReadByte()
			must(err)
			if b == 0 {
				return false
			} 
			if b == 1 {
				return true
			}
			panic(fmt.Errorf("invalid bool constant: %d", b))

		case LUA_NUM_INT: // INT
			var i64 int64
			must(binary.Read(r, order, &i64))
			return i64

		case LUA_NUM_FLOAT: // FLOAT
			var f64 float64
			must(binary.Read(r, order, &f64))
			return f64

		case LUA_STR_SHORT, LUA_STR_LONG: // STRING
			return decodeString(r)
	}
	panic(fmt.Errorf("unexpected constant type: %d", t))
}

func decodeString(r *bytes.Buffer) string {
	b, err := r.ReadByte()
	must(err)
	switch {
	case b == 0x00:
		return ""
	case b == 0xFF:
		var u64 uint64
		must(binary.Read(r, order, &u64))
		return string(r.Next(int(u64) - 1))
	default:
		return string(r.Next(int(b) - 1))
	}
}

func (bin *source) error(why string) {
	msg := fmt.Sprintf("%s: %s precompiled chunk", bin.name, why)
	panic(Error(msg))
}

func isValid(data []byte) bool {
	return len(data) > 4 && string(data[:4]) == LUA_SIGNATURE
}

func check(cond bool, mesg string) {
	if !cond {
		panic(fmt.Errorf(mesg))
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}