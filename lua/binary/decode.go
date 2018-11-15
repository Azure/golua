package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func decode(r *bytes.Buffer, c *Chunk) {
	var h Header
	must(binary.Read(r, order, &h))

	assert(h.Signature == head, "not a precompiled chunk")
	assert(h.Version == LUAC_VERSION, "version mismatch")
	assert(h.Format == LUAC_FORMAT, "format mismatch")
	assert(h.LuacData == tail, "corrupted")
	assert(h.GoIntSize == CINT_SIZE, "int size mismatch")
	assert(h.SizetSize == CSIZET_SIZE, "size_t size mismatch")
	assert(h.InstrSize == INSTRUCTION_SIZE, "instruction size mismatch")
	assert(h.LuaIntSize == LUA_INTEGER_SIZE, "lua integer size mismatch")
	assert(h.LuaNumSize == LUA_NUMBER_SIZE, "lua number size mismatch")
	assert(h.LuacIntEnc == LUAC_INT, "endianess mismatch")
	assert(h.LuacNumEnc == LUAC_NUM, "float format mismatch")
	c.Header = h

	// decode size_upvalues (?)
	_, err := r.ReadByte()
	must(err)

	// decode container closure prototype
	decodePrototype(r, &c.Entry)
}

func decodePrototype(r *bytes.Buffer, proto *Prototype) {
	// decode source name string (b[0] == length)
	// b[0] == 0xFF ? uint64 : size
	proto.Source = decodeString(r)

	// decode line start
	must(binary.Read(r, order, &proto.SrcPos))

	// decode line end
	must(binary.Read(r, order, &proto.EndPos))

	// decode number of parameters
	must(binary.Read(r, order, &proto.Params))

	// decode is varadic
	must(binary.Read(r, order, &proto.Vararg))

	// decode maximum stack size
	must(binary.Read(r, order, &proto.Stack))

	// decode instruction bytecode
	//
	// leading 4-bytes is number of instructions
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.Code = make([]uint32, num)
		for i := range proto.Code {
			must(binary.Read(r, order, &proto.Code[i]))
		}
	}

	// decode constants
	//
	// leading 4-bytes is number of constants
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.Consts = make([]interface{}, num)
		for i := range proto.Consts {
			proto.Consts[i] = decodeConst(r)
		}
	}

	// decode upvalues
	//
	// leading 4-bytes is number of upvalues
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.UpValues = make([]UpValue, num, num)
		for i := range proto.UpValues {
			var (
				upv UpValue
				err error
			)
			upv.InStack, err = r.ReadByte()
			must(err)

			upv.Index, err = r.ReadByte()
			must(err)

			proto.UpValues[i] = upv
		}
	}

	// decode nested closure prototypes
	//
	// leading 4-bytes is number of prototypes
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.Protos = make([]Prototype, num)
		for i := range proto.Protos {
			var fn Prototype
			decodePrototype(r, &fn)
			proto.Protos[i] = fn
		}
	}

	// decode line info (pc -> line)
	//
	// leading 4-bytes is number of pcln entries
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.PcLnTab = make([]uint32, num)
		for i := range proto.PcLnTab {
			must(binary.Read(r, order, &proto.PcLnTab[i]))
		}
	}

	// decode local variables
	//
	// leading 4-bytes is number of pcln entries
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.Locals = make([]LocalVar, num)
		for i := range proto.Locals {
			local := LocalVar{Name: decodeString(r)}
			must(binary.Read(r, order, &local.Live))
			must(binary.Read(r, order, &local.Dead))
			proto.Locals[i] = local
		}
	}

	// decode upvalue names
	//
	// leading 4-bytes is number of pcln entries
	{
		var num uint32
		must(binary.Read(r, order, &num))
		proto.UpNames = make([]string, num)
		for i := range proto.UpNames {
			proto.UpNames[i] = decodeString(r)
		}
	}
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

func decodeConst(r *bytes.Buffer) interface{} {
	t, err := r.ReadByte()
	must(err)
	switch t {
	case LUA_TYPE_NIL: // NIL
		return nil
	case LUA_TYPE_BOOL: // BOOL
		b, err := r.ReadByte()
		must(err)
		var v bool
		if b == 0 {
			v = false
		} else if b == 1 {
			v = true
		} else {
			panic(fmt.Errorf("invalid bool constant: %d", b))
		}
		return v
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
	default:
		panic(fmt.Errorf("unexpected constant type: %d", t))
	}
}
