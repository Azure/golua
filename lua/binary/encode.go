package binary

import (
	"encoding/binary"
	"io"
)

func encode(w io.Writer, p *Prototype) {
	binary.Write(w, order, head)
	binary.Write(w, order, byte(len(p.UpValues)))	
}