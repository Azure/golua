package pattern

import (
	"errors"
	"fmt"
)

var _ = fmt.Println

func MustCompile(expr string) *Pattern {
	patt, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return patt
}

func Compile(expr string) (*Pattern, error) {
	patt, err := compile(scan(expr), 0)
	return &Pattern{patt}, err
}

type opcode int

const (
	opClass opcode = iota
	opSplit
	opMatch
	opStart
	opClose
	opJump
	// opSave
)

var opname = [...]string{
	opClass: "char",
	opSplit: "split",
	opMatch: "match",
	opStart: "start",
	opClose: "close",
	opJump:  "jump",
	// opSave:  "save",
}

func (op opcode) String() string { return opname[op] }

type instr struct {
	code opcode
	item item
	x, y int
}

func (inst instr) String() string {
	switch inst.code {
	case opClass:
		return fmt.Sprintf("class %v (%q, %s)", classID(inst.item), inst.item.val, inst.item.rep)
	case opSplit:
		return fmt.Sprintf("split %d, %d", inst.x, inst.y)
	case opStart:
		return fmt.Sprintf("start %d", inst.x)
	case opClose:
		return fmt.Sprintf("close %d", inst.x)
	case opMatch:
		return "match"
	case opJump:
		return fmt.Sprintf("jump %d", inst.x)
	// case opSave:
	// 	return fmt.Sprintf("save %d", inst.x)
	}
	return fmt.Sprintf("unknown opcode %d", inst.code)
}

type builder struct {
	caps []capture
	inst []instr
	cap  int
}

func (bldr *builder) init(maxcaps int) {
	if maxcaps <= 0 {
		maxcaps = 32
	}
	bldr.caps = make([]capture, 0, maxcaps)
	// bldr.start()
}

func (bldr *builder) start() {
	bldr.caps = append(bldr.caps, capture{})
	bldr.caps[bldr.cap].start(-1)
	bldr.caps[bldr.cap].close(-1)
	bldr.emit(instr{code: opStart, x: bldr.cap})
	bldr.cap++
}

func (bldr *builder) close() {
	bldr.emit(instr{code: opClose, x: bldr.cap-1})
	// bldr.cap--
}

func (bldr *builder) class(item item) {
	bldr.emit(instr{code: opClass, item: item})
}

func (bldr *builder) split(l1, l2 int) {
	bldr.emit(instr{code: opSplit, x: l1, y: l2})
}

// func (bldr *builder) save(cap int) {
// 	bldr.emit(instr{code: opSave, x: cap})
// }

func (bldr *builder) jump(dst int) {
	bldr.emit(instr{code: opJump, x: dst})
}

func (bldr *builder) done() {
	// bldr.close()
	bldr.emit(instr{code: opMatch})
}

func (bldr *builder) emit(insts ...instr) {
	bldr.inst = append(bldr.inst, insts...)
}

func (bldr *builder) code(item item) {
	switch pc := len(bldr.inst); item.rep {
		case optional: // ?
			bldr.split(pc+1, pc+2)
			bldr.class(item)
		case maximum: // *
			bldr.split(pc+1, pc+3)
			bldr.class(item)
			bldr.jump(pc)
		case minimum: // -
			bldr.split(pc+3, pc+1)
			bldr.class(item)
			bldr.jump(pc)
		case greedy: // +
			bldr.class(item)
			bldr.split(pc, pc+2)
		default: // single
			bldr.class(item)
	}
}

func compile(scan *scanner, maxcaps int) (patt *pattern, err error) {
	var bldr builder
	bldr.init(maxcaps)

	for item := scan.nextItem(); item.typ != itemEnd; item = scan.nextItem() {
		// fmt.Printf("%v (%s)\n", item, item.rep)
		switch item.typ {
		case itemClass, itemText:
			bldr.code(item)
		case itemStartCapture:
			bldr.start()
		case itemCloseCapture:
			bldr.close()
		case itemErr:
			scan.drain()
			return nil, errors.New(item.val)
		}
	}
	bldr.done()
	return &pattern{
		caps: bldr.caps,
		inst: bldr.inst,
		head: scan.head,
		tail: scan.tail,
	}, nil
}