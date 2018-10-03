package binary

import (
	"encoding/binary"
	"bytes"
	"math"
	"fmt"
)

var order = binary.LittleEndian

// TODO: binary writer to write binary chunks.
// TODO: binary reader to read binary chunks.

type writer struct {
	b *bytes.Buffer
}

func (w *writer) writeHeader() {
	h := &Header{
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
	binary.Write(w.b, order, h)
}

func (w *writer) writeByte(b byte) {
	w.b.WriteByte(b)
}

func (w *writer) writeU32(u32 uint32) {
	binary.Write(w.b, order, u32)
}

func (w *writer) writeU64(u64 uint64) {
	binary.Write(w.b, order, u64)
}

func (w *writer) writeF64(f64 float64) {
	w.writeU64(math.Float64bits(f64))
}

func (w *writer) writeI64(i64 int64) {
	w.writeU64(uint64(i64))
}

func (w *writer) writeStr(str string) {
	strlen := len(str)
	if strlen == 0 {
		w.writeByte(0)
		return
	}
	strlen++
	if strlen >= 0xFF {
		w.writeByte(0xFF)
		w.writeU64(uint64(strlen))
	} else {
		w.writeByte(byte(strlen))
	}
	binary.Write(w.b, order, []byte(str))
}

func (w *writer) writeCode(codes []uint32) {
	w.writeU32(uint32(len(codes)))
	for _, code := range codes {
		w.writeU32(code)
	}
}

func (w *writer) writeConsts(consts []interface{}) {
	w.writeU32(uint32(len(consts)))
	for _, kst := range consts {
		switch kst := kst.(type) {
			case float64:
				w.writeByte(LUA_NUM_FLOAT)
				w.writeF64(kst)
			case string:
				if len(kst) > 40 {
					w.writeByte(LUA_STR_SHORT)
				} else {
					w.writeByte(LUA_STR_LONG)
				}
				w.writeStr(kst)
			case int64:
				w.writeByte(LUA_NUM_INT)
				w.writeI64(kst)
			case bool:
				w.writeByte(LUA_TYPE_BOOL)
				b := byte(0)
				if kst { b = 1 }
				w.writeByte(b)
			case nil:
				w.writeByte(LUA_TYPE_NIL)
			default:
				panic(fmt.Errorf("write: unknown constant type: %T", kst))
		}
	}
}

func (w *writer) writeUpValues(upvalues []UpValue) {
	w.writeU32(uint32(len(upvalues)))
	for _, upvalue := range upvalues {
		w.writeByte(upvalue.InStack)
		w.writeByte(upvalue.Index)
	}
}

func (w *writer) writeProtos(protos []Prototype) {
	w.writeU32(uint32(len(protos)))
	for _, proto := range protos {
		encodeProto(w, &proto)
	}
}

func (w *writer) writePcLnInfo(pcln []uint32) {
	w.writeU32(uint32(len(pcln)))
	for _, ln := range pcln {
		w.writeU32(ln)
	}
}

func (w *writer) writeLocalVars(vars []LocalVar) {
	w.writeU32(uint32(len(vars)))
	for _, loc := range vars {
		w.writeStr(loc.Name)
		w.writeU32(loc.Live)
		w.writeU32(loc.Dead)
	}
}

func (w *writer) writeUpValueNames(names []string) {
	w.writeU32(uint32(len(names)))
	for _, name := range names {
		w.writeStr(name)
	}
}
