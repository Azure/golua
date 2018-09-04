package op

// Masks for instruction arguments.
type ArgMask uint8

const (
	ArgN ArgMask = iota // argument is not used
	ArgU 				// argument is used
	ArgR 				// argument is a register or a jump offset
	ArgK 				// argument is a constant or register/constant
)

// ArgMask represents masks for instruction properties.
//
// The format is:
// bits 0-1: op mode
// bits 2-3: C arg mode
// bits 4-5: B arg mode
// bit 6: instruction set register A
// bit 7: operator is a test (next instruction must be a jump) 
type Mask uint8

var masks = [...]Mask{
    MOVE:     mask(0, 1, ArgR, ArgN, ModeABC),     // MOVE
    LOADK:    mask(0, 1, ArgK, ArgN, ModeABx),     // LOADK
    LOADKX:   mask(0, 1, ArgN, ArgN, ModeABx),     // LOADKX
    LOADBOOL: mask(0, 1, ArgU, ArgU, ModeABC),     // LOADBOOL
    LOADNIL:  mask(0, 1, ArgU, ArgN, ModeABC),     // LOADNIL
    GETUPVAL: mask(0, 1, ArgU, ArgN, ModeABC),     // GETUPVAL
    GETTABUP: mask(0, 1, ArgU, ArgK, ModeABC),     // GETTABUP
    GETTABLE: mask(0, 1, ArgR, ArgK, ModeABC),     // GETTABLE
    SETTABUP: mask(0, 0, ArgK, ArgK, ModeABC),     // SETTABUP
    SETUPVAL: mask(0, 0, ArgU, ArgN, ModeABC),     // SETUPVAL
    SETTABLE: mask(0, 0, ArgK, ArgK, ModeABC),     // SETTABLE
    NEWTABLE: mask(0, 1, ArgU, ArgU, ModeABC),     // NEWTABLE
    SELF:     mask(0, 1, ArgR, ArgK, ModeABC),     // SELF
    ADD:      mask(0, 1, ArgK, ArgK, ModeABC),     // ADD
    SUB:      mask(0, 1, ArgK, ArgK, ModeABC),     // SUB
    MUL:      mask(0, 1, ArgK, ArgK, ModeABC),     // MUL
    MOD:      mask(0, 1, ArgK, ArgK, ModeABC),     // MOD
    POW:      mask(0, 1, ArgK, ArgK, ModeABC),     // POW
    DIV:      mask(0, 1, ArgK, ArgK, ModeABC),     // DIV
    IDIV:     mask(0, 1, ArgK, ArgK, ModeABC),     // IDIV
    BAND:     mask(0, 1, ArgK, ArgK, ModeABC),     // BAND
    BOR:      mask(0, 1, ArgK, ArgK, ModeABC),     // BOR
    BXOR:     mask(0, 1, ArgK, ArgK, ModeABC),     // BXOR
    SHL:      mask(0, 1, ArgK, ArgK, ModeABC),     // SHL
    SHR:      mask(0, 1, ArgK, ArgK, ModeABC),     // SHR
    UNM:      mask(0, 1, ArgR, ArgN, ModeABC),     // UNM
    BNOT:     mask(0, 1, ArgR, ArgN, ModeABC),     // BNOT
    NOT:      mask(0, 1, ArgR, ArgN, ModeABC),     // NOT
    LEN:      mask(0, 1, ArgR, ArgN, ModeABC),     // LEN
    CONCAT:   mask(0, 1, ArgR, ArgR, ModeABC),     // CONCAT
    JMP:      mask(0, 0, ArgR, ArgN, ModeAsBx),    // JMP
    EQ:       mask(1, 0, ArgK, ArgK, ModeABC),     // EQ
    LT:       mask(1, 0, ArgK, ArgK, ModeABC),     // LT
    LE:       mask(1, 0, ArgK, ArgK, ModeABC),     // LE
    TEST:     mask(1, 0, ArgN, ArgU, ModeABC),     // TEST
    TESTSET:  mask(1, 1, ArgR, ArgU, ModeABC),     // TESTSET
    CALL:     mask(0, 1, ArgU, ArgU, ModeABC),     // CALL
    TAILCALL: mask(0, 1, ArgU, ArgU, ModeABC),     // TAILCALL
    RETURN:   mask(0, 0, ArgU, ArgN, ModeABC),     // RETURN
    FORLOOP:  mask(0, 1, ArgR, ArgN, ModeAsBx),    // FORLOOP
    FORPREP:  mask(0, 1, ArgR, ArgN, ModeAsBx),    // FORPREP
    TFORCALL: mask(0, 0, ArgN, ArgU, ModeABC),     // TFORCALL
    TFORLOOP: mask(0, 1, ArgR, ArgN, ModeAsBx),    // TFORLOOP
    SETLIST:  mask(0, 0, ArgU, ArgU, ModeABC),     // SETLIST
    CLOSURE:  mask(0, 1, ArgU, ArgN, ModeABx),     // CLOSURE
    VARARG:   mask(0, 1, ArgU, ArgN, ModeABC),     // VARARG
    EXTRAARG: mask(0, 0, ArgU, ArgU, ModeAx),      // EXTRAARG
}

func (mask Mask) B(barg ArgMask) bool { return barg == ArgMask((mask>>4)&3) }

func (mask Mask) C(carg ArgMask) bool { return carg == ArgMask((mask>>2)&3) }

func (mask Mask) Mode() Mode { return Mode(mask & 3) }

func (mask Mask) SetA() bool { return mask&(1<<6) == 1 }

func (mask Mask) Test() bool { return mask&(1<<7) == 1 }

func mask(t, a uint8, b, c ArgMask, m Mode) Mask {
    return Mask((((t)<<7) | ((a)<<6) | ((uint8(b))<<4) | ((uint8(c))<<2) | (uint8(m))))
}