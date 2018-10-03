package binary

// Protos   []Prototype
// PcLnTab  []uint32
// Locals   []LocalVar
// UpNames  []string

func encodeProto(w *writer, p *Prototype) {
	w.writeStr(p.Source)
	w.writeU32(p.SrcPos)
	w.writeU32(p.EndPos)
	w.writeByte(p.Params)
	w.writeByte(p.Vararg)
	w.writeByte(p.Stack)
	w.writeCode(p.Code)
	w.writeConsts(p.Consts)
	w.writeUpValues(p.UpValues)
	w.writeProtos(p.Protos)
	w.writePcLnInfo(p.PcLnTab)
	w.writeLocalVars(p.Locals)
	w.writeUpValueNames(p.UpNames)
}