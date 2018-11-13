package packer

import (
	"encoding/binary"
	"bytes"
	"fmt"
)

type state struct {
	out bytes.Buffer
	opt []option
	arg int
}

func newState(format string) (*state, error) {
	var (
		s = scan(format)
		p = new(state)
	)
	L: for {
		switch opt := s.nextOpt(); opt.typ {
		case optErr:
			s.drain()
			return nil, fmt.Errorf(opt.value)
		case optEnd:
			break L
		default:
			p.opt = append(p.opt, opt)
		}
	}
	return p, nil
}

func (p *state) Option() Option {
	if p.arg < 0 || p.arg >= len(p.opt) {
		return nil
	}
	return p.opt[p.arg]
}

func (p *state) Unpack(values ...interface{}) (n int, err error) {
	return 0, fmt.Errorf("pack.Unpack: TODO")
}

func (p *state) Pack(values ...interface{}) ([]byte, error) {
	p.out.Reset()
	p.arg = 0
	for _, opt := range p.opt {
		switch {
		case opt.typ == optPad:
			p.out.WriteByte(0)
			continue
		case p.arg >= len(values):
			return nil, fmt.Errorf("bad argument #%d to 'pack'", p.arg)
		default:
			if err := p.packValue(opt, values[p.arg]); err != nil {
				return nil, err
			}
			p.arg++
		}
	}
	return p.out.Bytes(), nil
}

func (p *state) Size() (size int, err error) {
	for i, opt := range p.opt {
		if opt.typ == optVarLen || opt.typ == optPrefix {
			err = fmt.Errorf("bad argument #%d to 'packsize' (variable-length format)", i)
			break
		}
		size += int(opt.width)
	}
	return size, err
}


func (p *state) packValue(o option, v interface{}) error {
	switch o.typ {
	case optVarLen, optPrefix, optFixed:
		return p.packString(o, v)
	case optUint, optInt:
		return p.packInt(o, v)
	case optFloat:
		return p.packFloat(o, v)
	}
	panic("unreachable")
}

func (p *state) packString(o option, v interface{}) error {
	fmt.Println("pack: string")
	// return p.pack(o, v)
	return nil
}

func (p *state) packFloat(o option, v interface{}) error {
	fmt.Println("pack: float")
	// return p.pack(o, v)
	return nil
}

func (p *state) packInt(o option, v interface{}) error {
	// fmt.Printf("packInt (%t) => %T\n", o.typ == optUint, v)
	if v, ok := v.(Packer); ok {
		return v.Pack(p)
	}
	// return p.pack(o, v)
	return nil
}

func (p *state) pack(o option, v interface{}) error {
	return binary.Write(&p.out, o.order, v)
}