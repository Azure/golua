package pattern

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const classes = "acdglpsuwx"
const escape = '%'
const eos = rune(-1)

type itemType int

const (
	itemEnd itemType = iota
	itemErr
	itemText
	itemClass
	itemStartCapture
	itemCloseCapture
)

var itemTypes = [...]string{
	itemEnd:          "end of string",
	itemClass:        "character class",
	itemStartCapture: "start capture",
	itemCloseCapture: "close capture",
}

type repeatOp int

const (
	single repeatOp = iota
	greedy
	minimum
	maximum
	optional
)

func (op repeatOp) String() string {
	switch op {
	case greedy:
		return "greedy"
	case minimum:
		return "minimum"
	case maximum:
		return "maximum"
	case optional:
		return "optional"
	}
	return "single"
}

type item struct {
	typ itemType
	pos int
	val string
	rep repeatOp
}

func (i item) String() (s string) {
	switch i.typ {
	case itemErr:
		return fmt.Sprintf("error: %s", i.val)
	case itemText:
		return fmt.Sprintf("text: %q", i.val)
	case itemClass:
		return fmt.Sprintf("class: %q", i.val)
	}
	return itemTypes[i.typ]
}

type stateFn func(*scanner) stateFn

type scanner struct {
	item  chan item
	expr  string
	head  bool
	tail  bool
	size  int
	ncap  int
	caps  int
	pos   int
	start int
}

func scan(expr string) *scanner {
	s := &scanner{expr: expr, item: make(chan item)}
	go s.run()
	return s
}

func (scan *scanner) errorf(format string, args ...interface{}) stateFn {
	scan.item <- item{itemErr, scan.start, fmt.Sprintf(format, args...), 0}
	return nil
}

func (scan *scanner) nextItem() item { return <-scan.item }

func (scan *scanner) backup() { scan.pos -= scan.size }

func (scan *scanner) ignore() { scan.start = scan.pos }

func (scan *scanner) drain() {
	for range scan.item {
	}
}

func (scan *scanner) emit(typ itemType, rep repeatOp) {
	scan.item <- item{typ, scan.start, scan.expr[scan.start:scan.pos], rep}
	scan.start = scan.pos
}

func (scan *scanner) next() (r rune) {
	if scan.pos >= len(scan.expr) {
		scan.size = 0
		return eos
	}
	r, scan.size = utf8.DecodeRuneInString(scan.expr[scan.pos:])
	scan.pos += scan.size
	return r
}

func (scan *scanner) peek() (r rune) {
	r = scan.next()
	scan.backup()
	return r
}

func (scan *scanner) rep() repeatOp {
	switch scan.next() {
	case '+':
		return greedy
	case '?':
		return optional
	case '-':
		return minimum
	case '*':
		return maximum
	}
	scan.backup()
	return single
}

func (scan *scanner) run() {
	for state := scanText; state != nil; {
		state = state(scan)
	}
	if scan.ncap > 0 {
		scan.errorf("unfinished capture")
	}
	scan.emit(itemEnd, 0)
	close(scan.item)
}

func scanText(scan *scanner) stateFn {
	trace("scanText")
	switch r := scan.next(); r {
	case '(', ')': // start capture
		scan.backup()
		return scanCapture
	case escape:
		return scanEscape
	case eos:
		return nil
	case '^':
		if scan.start == 0 {
			scan.head = true
			scan.ignore()
			return scanText
		}
	case '$':
		if scan.peek() == eos {
			scan.tail = true
			scan.ignore()
			return scanText
		}
	}
	scan.backup()
	return scanSingle
}

func scanSingle(scan *scanner) stateFn {
	trace("scanSingle")
	switch r := scan.next(); r {
	case '[':
		return scanUnion
	case '.':
		scan.item <- item{itemClass, scan.start, ".", scan.rep()}
		scan.ignore()
	default:
		scan.item <- item{itemText, scan.start, string(r), scan.rep()}
		scan.ignore()
	}
	return scanText
}

func scanEscape(scan *scanner) stateFn {
	trace("scanEscape")
	switch r := scan.next(); {
	case unicode.IsDigit(r):
		return scan.errorf("todo: capture index")
	case r == 'b':
		return scan.errorf("todo: balance")
	case r == 'f':
		return scan.errorf("todo: frontier")
	case r == eos:
		return scan.errorf("malformed pattern (%s)", "ends with '%'")
	default:
		var (
			typ itemType = itemText
			rep repeatOp
			lit string
			pos int
		)
		if lit, pos = string(r), scan.start+1; isclass(r) {
			typ = itemClass
			lit = string(r)
			rep = scan.rep()
		}
		scan.item <- item{typ, pos, lit, rep}
		scan.ignore()
	}
	return scanText
}

func scanCapture(scan *scanner) stateFn {
	// TODO: positional captures
	trace("scanCapture")
	switch scan.next() {
	case '(': // start capture
		scan.emit(itemStartCapture, 0)
		scan.caps++
		scan.ncap++
	case ')': // end capture
		scan.emit(itemCloseCapture, 0)
		scan.ncap--
	}
	return scanText
}

// [set]: represents the class which is the union of all characters in set.
// A range of characters can be specified by separating the end characters
// of the range, in ascending order, with a '-'. All classes %x described
// above can also be used as components in set. All other characters in set
// represent themselves. For example, [%w_] (or [_%w]) represents all
// alphanumeric characters plus the underscore, [0-7] represents the octal
// digits, and [0-7%l%-] represents the octal digits plus the lowercase
// letters plus the '-' character.
//
// You can put a closing square bracket in a set by positioning it as the
// first character in the set. You can put a hyphen in a set by positioning
// it as the first or the last character in the set. (You can also use an
// escape for both cases.)
//
// The interaction between ranges and classes is not defined. Therefore,
// patterns like [%a-z] or [a-%%] have no meaning.
//
// [^set]: represents the complement of set, where set is interpreted as above.
func scanUnion(scan *scanner) stateFn {
	return scan.errorf("union set classses not implemented")
}
