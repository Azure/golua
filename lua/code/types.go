package code

type Type int8

const NoType = Type(-1)

const (
	NilType     = Type(0)
	BoolType    = Type(1)
	PointerType = Type(2)
	NumberType  = Type(3)
	StringType  = Type(4)
	TableType   = Type(5)
	FuncType    = Type(6)
	GoType      = Type(7)
	ThreadType  = Type(8)
	MaxType
)

const (
	FloatType = Type(NumberType | (0 << 4))
	IntType   = Type(NumberType | (1 << 4))
)

var typeNames = [...]string{
	NilType:     "nil",
	BoolType:    "boolean",
	PointerType: "light userdata",
	NumberType:  "number",
	StringType:  "string",
	TableType:   "table",
	FuncType:    "function",
	GoType:      "userdata",
	ThreadType:  "thread",
}

func (typ Type) String() string {
	switch 0x0F & typ {
	case NumberType:
		return "number"
	case NoType:
		return "no value"
	case NilType:
		return "nil"
	case BoolType:
		return "boolean"
	case PointerType:
		return "light userdata"
	case StringType:
		return "string"
	case TableType:
		return "table"
	case FuncType:
		return "function"
	case GoType:
		return "userdata"
	case ThreadType:
		return "thread"
	}
	panic(typ)
}
