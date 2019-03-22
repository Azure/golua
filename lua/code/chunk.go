package code

import (
	"bytes"
	"io"
)

type Chunk struct {
	Main *Proto
}

// func Undump(name string, src interface{}) (*Chunk, error) {
// 	if len(name) > 1 && (name[0] == '@' || name[0] == '=') {
// 		name = name[1:]
// 	}
// 	b, err := readSource(name, src)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var main Proto
// 	decodeChunk(b, &main, name)
// 	// TODO: verifycode
// 	return &Chunk{&main}, nil
// }

func (chunk *Chunk) Dump(out io.Writer, strip bool) (int, error) {
	w := &source{ord: order, src: new(bytes.Buffer), strip: strip}
	n := byte(len(chunk.Main.UpVars))
	must(w.writeHeader())
	must(w.write(n))
	w.writeProto(chunk.Main, "")
	return out.Write(w.src.Bytes())
}

func (chunk *Chunk) Print(w io.Writer, full bool) {
	printFunc(w, chunk.Main, full)
}