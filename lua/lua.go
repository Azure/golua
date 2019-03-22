package lua

import (
	"github.com/Azure/golua/lua/luac"
)

const EnvID = "_ENV"

type (
	Importer interface {
		Import(*Thread, string) error
	}

	Loader func(*Thread) (Value, error)

	Library struct {
		Funcs []*GoFunc
		Open  Loader
		Name  string
	}
)

type Op int

const (
	OpNone Op = iota

	// binary ops
	OpAdd
	OpSub
	OpMul
	OpMod
	OpPow
	OpDivF
	OpDivI
	OpBand
	OpBor
	OpBxor
	OpShl
	OpShr
	OpConcat
	OpEq
	OpLt
	OpLe
	OpNe
	OpGt
	OpGe
	OpAnd
	OpOr

	// unary ops
	OpMinus
	OpBnot
	OpNot
	OpLen
)

var opnames = [...]string{
	OpNone:   "none",
	OpAdd:    "add",
	OpSub:    "sub",
	OpMul:    "mul",
	OpMod:    "mod",
	OpPow:    "pow",
	OpDivF:   "fdiv",
	OpDivI:   "idiv",
	OpBand:   "band",
	OpBor:    "bor",
	OpBxor:   "bxor",
	OpShl:    "shl",
	OpShr:    "shr",
	OpConcat: "concat",
	OpEq:     "eq",
	OpLt:     "lt",
	OpLe:     "le",
	OpNe:     "ne",
	OpGt:     "gt",
	OpGe:     "ge",
	OpAnd:    "and",
	OpOr:     "or",
	OpMinus:  "minus",
	OpBnot:   "bnot",
	OpNot:    "not",
	OpLen:    "len",
}

func (op Op) String() string { return opnames[op] }

func LoadFile(t *Thread, file string) (*Func, error) {
	chunk, err := luac.Compile(luac.Defaults, file, nil)
	if err != nil {
		return nil, err
	}
	return t.Load(chunk), nil
}

func Must(ls *Thread, err error) *Thread {
	if err != nil {
		panic(err)
	}
	return ls
}

func Init(config *Config) (*Thread, error) {
	ls := new(runtime).init(config)
	ls.tt = &Thread{ls}
	return ls.tt, config.Stdlib(ls.tt)
}