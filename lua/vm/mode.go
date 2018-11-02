package vm

type Mode uint8

const (
    ModeABC Mode = iota
    ModeABx
    ModeAsBx
    ModeAx
)

var modes = [...]string{
	ModeABC:  "ABC",
	ModeABx:  "ABx",
	ModeAsBx: "AsBx",
	ModeAx:   "Ax",
}

func (mode Mode) String() string { return modes[mode] }