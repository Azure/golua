package vm

type Code uint8

const (
	MOVE	 Code = iota
	LOADK
	LOADKX
	LOADBOOL
	LOADNIL
	GETUPVAL
	GETTABUP
	GETTABLE
	SETTABUP
	SETUPVAL
	SETTABLE
	NEWTABLE
	SELF
	ADD
	SUB
	MUL
	MOD
	POW
	DIV
	IDIV
	BAND
	BOR
	BXOR
	SHL
	SHR
	UNM
	BNOT
	NOT
	LEN
	CONCAT
	JMP
	EQ
	LT
	LE
	TEST
	TESTSET
	CALL
	TAILCALL
	RETURN
	FORLOOP
	FORPREP
	TFORCALL
	TFORLOOP
	SETLIST
	CLOSURE
	VARARG
	EXTRAARG
)

var names = [...]string{
	MOVE: 	  "MOVE",
	LOADK:    "LOADK",
	LOADKX:   "LOADKX",
	LOADBOOL: "LOADBOOL",
	LOADNIL:  "LOADNIL",
	GETUPVAL: "GETUPVAL",
	GETTABUP: "GETTABUP",
	GETTABLE: "GETTABLE",
	SETTABUP: "SETTABUP",
	SETUPVAL: "SETUPVAL",
	SETTABLE: "SETTABLE",
	NEWTABLE: "NEWTABLE",
	SELF:     "SELF",
	ADD:      "ADD",
	SUB:      "SUB",
	MUL:      "MUL",
	MOD:      "MOD",
	POW:      "POW",
	DIV:      "DIV",
	IDIV:     "IDIV",
	BAND:     "BAND",
	BOR:      "BOR",
	BXOR:     "BXOR",
	SHL:      "SHL",
	SHR:      "SHR",
	UNM:      "UNM",
	BNOT:     "BNOT",
	NOT:      "NOT",
	LEN:      "LEN",
	CONCAT:   "CONCAT",
	JMP:      "JMP",
	EQ: 	  "EQ",
	LT: 	  "LT",
	LE: 	  "LE",
	TEST: 	  "TEST",
	TESTSET:  "TESTSET",
	CALL: 	  "CALL",
	TAILCALL: "TAILCALL",
	RETURN:   "RETURN",
	FORLOOP:  "FORLOOP",
	FORPREP:  "FORPREP",
	TFORCALL: "TFORCALL",
	TFORLOOP: "TFORLOOP",
	SETLIST:  "SETLIST",
	CLOSURE:  "CLOSURE",
	VARARG:   "VARARG",
	EXTRAARG: "EXTRAARG",
}

func (op Code) Mask() Mask { return masks[op] }

func (op Code) Mode() Mode { return masks[op].Mode() }

func (op Code) String() string { return names[op] }