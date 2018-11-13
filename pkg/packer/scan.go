package packer

import (
	"encoding/binary"
	"unicode/utf8"
	"math"
	"fmt"
)

var order = binary.LittleEndian

const align = 1
const eos = rune(-1)

type stateFn func(*scanner) stateFn

type scanner struct {
	order binary.ByteOrder
	opts  chan option
	fmts  string
	size  int
	pos   int
	start int
}

type optType int

const (
	optEnd   optType = iota // error
	optErr    			   // end of string
	optPad 				   // padding
	optInt  			   // signed integers
	optUint 	 	   	   // unsigned integers
	optFloat 			   // floating-point numbers
	optFixed 			   // fixed-length strings
	optPrefix			   // strings with prefixed length
	optVarLen 			   // variable-length strings
)

type option struct {
	order binary.ByteOrder
	value string
	align uint
	width uint
	start int
	verb  rune
	typ   optType
}

func scan(fmts string) *scanner {
	s := &scanner{order: order, fmts: fmts, opts: make(chan option)}
	go s.run()
	return s
}

func (scan *scanner) errorf(format string, args ...interface{}) stateFn {
	scan.opts <- option{nil, fmt.Sprintf(format, args...), 0, 0, scan.start, -1, optErr}
	return nil
}

func (scan *scanner) nextOpt() option { return <-scan.opts }
 
func (scan *scanner) backup() { scan.pos -= scan.size }

func (scan *scanner) ignore() { scan.start = scan.pos }

func (scan *scanner) drain() { for range scan.opts {} }

func (scan *scanner) emit(typ optType, verb rune, align, width uint) {
	scan.opts <- option{scan.order, scan.fmts[scan.start:scan.pos], align, width, scan.start, verb, typ}
	scan.start = scan.pos
}

func (scan *scanner) next() (r rune) {
	if scan.pos >= len(scan.fmts) {
		scan.size = 0
		return eos
	}
	r, scan.size = utf8.DecodeRuneInString(scan.fmts[scan.pos:])
	scan.pos += scan.size
	return r
}

func (scan *scanner) peek() (r rune) {
	r = scan.next()
	scan.backup()
	return r
}

func (scan *scanner) run() {
	for state := scanFmt; state != nil; {
		state = state(scan)
	}
	scan.emit(optEnd, -1, 0, 0)
	close(scan.opts)
}

// TODO: ![n] => set maximum alignment
// TODO: Xop => set alignment op
func scanFmt(scan *scanner) stateFn {
	L: switch r := scan.next(); r {
	case '<': // sets little endian
		scan.order = binary.LittleEndian
		scan.ignore()
		goto L
	case '>': // sets big endian
		scan.order = binary.BigEndian
		scan.ignore()
		goto L
	case '=': // sets native endian
		scan.order = order
		scan.ignore()
		goto L
	case '!': // ![n]: sets maximum alignment to n (default is native alignment) 
		n, err := optSize(scan, 4)
		if err != nil {
			return scan.errorf("%v", err)
		}
		_ = n // TODO
		scan.ignore()
		goto L
	case 'x': // one byte of padding
		scan.emit(optPad, r, 0, 1)
		scan.ignore()
		goto L
	case 'X': // Xop: an empty item that aligns according to option op (which is otherwise ignored)
		// TODO
		scan.ignore()
		goto L
	case ' ': // empty space
		scan.ignore()
		goto L
	case eos:
		return nil
	}
	scan.backup()
	return scanOpt
}

// TODO: limit
// TODO: align
func scanOpt(scan *scanner) stateFn {
	switch r := scan.next(); r {	
	case 'b': // a signed byte (char)
		scan.emit(optInt, r, 0, 1)
	case 'B': // an unsigned byte (char)
		scan.emit(optUint, r, 0, 1)
	case 'h': // a signed short (native size)
		scan.emit(optInt, r, 0, 2)
	case 'H': // an unsigned short (native size)
		scan.emit(optUint, r, 0, 2)
	case 'l': // a signed long (native size)
		scan.emit(optInt, r, 0, 8)
	case 'L': // an unsigned long (native size)
		scan.emit(optUint, r, 0, 8)
	case 'j': // a lua_Integer
		scan.emit(optInt, r, 0, 8)
	case 'J': // a lua_Unsigned
		scan.emit(optUint, r, 0, 8)
	case 'i': // i[n]: a signed int with n bytes (default is native size)
		n, err := optSize(scan, 4)
		if err != nil {
			return scan.errorf("%v", err)
		}
		scan.emit(optInt, r, 0, n)
	case 'I': // I[n]: an unsigned int with n bytes (default is native size)
		n, err := optSize(scan, 4)
		if err != nil {
			return scan.errorf("%v", err)
		}
		scan.emit(optUint, r, 0, n)
	case 'T': // a size_t (native size)
		scan.emit(optUint, r, 0, 8)
	case 'f': // a float (native size)
		scan.emit(optFloat, r, 0, 4)
	case 'd': // a double (native size)
		scan.emit(optFloat, r, 0, 8)
	case 'n': // a lua_Number
		scan.emit(optInt, r, 0, 8)
	case 'c': // cn: a fixed-sized string with n bytes
		n, err := number(scan, 0)
		if err != nil {
			return scan.errorf("%v", err)
		}
		if n == 0 {
			return scan.errorf("missing size for format option 'c'")
		}
		scan.emit(optFixed, r, 0, n)
	case 'z': // a zero-terminated string
		scan.emit(optVarLen, r, 0, 0)
	case 's': // s[n]: a string preceded by its length coded as an unsigned integer with n bytes (default is a size_t)
		n, err := optSize(scan, 8)
		if err != nil {
			return scan.errorf("%v", err)
		}
		scan.emit(optPrefix, r, 0, n)
	default:
		return scan.errorf("invalid format option '%c'", r)
	}
	return scanFmt
}

// For options "!n", "sn", "in", and "In", n can be any integer between 1 and 16.
func optSize(scan *scanner, opt uint) (size uint, err error) {
	if size, err = number(scan, opt); err != nil {
		return 0, err
	}
	if size < 1 || size > 16 {
		return 0, fmt.Errorf("integral size (%d) out of limits [1,16]", size)
	}
	return size, nil
}

func number(scan *scanner, opt uint) (num uint, err error) {
	if !isDigit(scan.peek()) {
		return opt, nil
	}
	for isDigit(scan.peek()) {
		if num = (num * 10) + uint(scan.next() - '0'); num > maxsize {
			return 0, fmt.Errorf("option size overflow")
		}
	}
	return num, nil
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isDigit(r rune) bool { return r >= '0' && r <= '9' }

const maxsize = math.MaxUint64 / 10