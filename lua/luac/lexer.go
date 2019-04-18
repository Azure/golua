package luac

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = -1

// lexical state
type lexical struct {
	*scanner

	active []*variable
	labels []*label
	gotos  []*label
	fs     *function
	// for debugging
	indent int
}

type scanner struct {
	source *bufio.Reader
	buffer bytes.Buffer
	token  token
	peek0  token
	name   string
	char   rune
	line   int
	last   int
}

func (s *lexical) tok2str(char rune) string {
	if char < reserved { // single-byte symbols?
		return fmt.Sprintf("'%c'", char)
	}
	str := tokens[int(char)-reserved]
	if char < tEOS {
		return fmt.Sprintf("'%s'", str)
	}
	return str
}

func (s *lexical) addline() {
	s.assert(isNewLine(s.char))
	old := s.char
	s.read()
	if isNewLine(s.char) && s.char != old {
		s.read()
	}
	if s.line++; s.line >= maxInt {
		s.syntaxErr("chunk has too many lines")
	}
}

func (s *lexical) consume() { s.accept(s.char) }

func (s *lexical) accept(char rune) {
	s.save(char)
	s.read()
}

func (s *lexical) read() {
	if c, err := s.source.ReadByte(); err != nil {
		s.char = eof
	} else {
		s.char = rune(c)
	}
}

// check
func (s *lexical) expect(tok rune) {
	if s.token.char != tok {
		s.expectErr(tok)
	}
}

func (s *lexical) oneof(chars string) bool {
	if s.char == 0 || !strings.ContainsRune(chars, s.char) {
		return false
	}
	s.consume()
	return true
}

// check_match
func (s *lexical) match(what, who rune, where int) {
	if !s.test(what) {
		if where == s.line {
			s.expectErr(what)
		} else {
			errMsg := "%s expected (to close %s at line %d)"
			s.syntaxErr(fmt.Sprintf(errMsg, s.tok2str(what), s.tok2str(who), where))
		}
	}
}

func (s *lexical) test(want rune) (ok bool) {
	if ok = (want == s.token.char); ok {
		s.next()
	}
	return
}

func (s *lexical) save(char rune) {
	if err := s.buffer.WriteByte(byte(char)); err != nil {
		s.scanErr("lexical element too long", 0)
	}
}

func (s *lexical) scan() token {
	var comment bool
	for {
		if isNewLine(s.char) {
			s.addline()
			continue
		}
		if isSpace(s.char) {
			s.read()
			continue
		}
		switch char := s.char; char {
		case eof:
			return token{char: tEOS}
		case '-':
			if s.read(); s.char != '-' {
				return token{char: '-'}
			}
			if s.read(); s.char == '[' {
				if sep := s.skipSep(); sep >= 0 {
					_ = s.multiline(comment, sep)
					break
				}
				s.buffer.Reset()
			}
			for !isNewLine(s.char) && s.char != eof {
				s.read()
			}
		case '[': // long string or simple '['
			sep := s.skipSep()
			if sep >= 0 {
				return token{char: tString, sval: s.multiline(false, sep)}
			}
			s.buffer.Reset()
			if sep != -1 {
				s.scanErr("invalid long string delimiter", tString)
			}
			return token{char: '['}
		case '=':
			if s.read(); s.char != '=' {
				return token{char: '='}
			}
			s.read()
			return token{char: tEq}
		case '<':
			switch s.read(); s.char {
			case '=':
				s.read()
				return token{char: tLe}
			case '<':
				s.read()
				return token{char: tShl}
			}
			return token{char: '<'}
		case '>':
			switch s.read(); s.char {
			case '=':
				s.read()
				return token{char: tGe}
			case '>':
				s.read()
				return token{char: tShr}
			}
			return token{char: '>'}
		case '~':
			if s.read(); s.char != '=' {
				return token{char: '~'}
			}
			s.read()
			return token{char: tNe}
		case ':':
			if s.read(); s.char != ':' {
				return token{char: ':'}
			}
			s.read()
			return token{char: tColon2}
		case '/':
			if s.read(); s.char != '/' {
				return token{char: '/'}
			}
			s.read()
			return token{char: tDivI}
		case '"', '\'':
			return s.strlit(char)
		case '.':
			if s.consume(); s.oneof(".") {
				if s.oneof(".") {
					s.buffer.Reset()
					return token{char: tDots}
				}
				s.buffer.Reset()
				return token{char: tConcat}
			}
			if !isDigit(s.char) {
				s.buffer.Reset()
				return token{char: '.'}
			}
			return s.numlit()
		case 0:
			s.read()
		default:
			if isDigit(s.char) {
				return s.numlit()
			}
			if isIdent(s.char) {
				s.consume()
				return s.ident()
			}
			s.read()
			return token{char: char}
		}
	}
}

func (s *lexical) peek() rune {
	s.assert(s.peek0.char == tEOS)
	s.peek0 = s.scan()
	return s.peek0.char
}

func (s *lexical) next() bool {
	if s.last = s.line; s.peek0.char != tEOS {
		s.token = s.peek0
		s.peek0.char = tEOS
	} else {
		s.token = s.scan()
	}
	return (s.char != eof)
}

func (s *scanner) init(name string, input interface{}) *scanner {
	b, err := readSource(name, input)
	if err != nil {
		panic(err)
	}
	// s.source = skipComment(bufio.NewReader(b))
	s.source = bufio.NewReader(b)
	s.peek0.char = tEOS
	s.buffer.Reset()
	s.name = name
	s.line = 1
	s.last = 1
	return skipComment(s)
}

func (s *lexical) ident() token {
	for isIdent(s.char) || isDigit(s.char) {
		s.consume()
	}
	str := s.buffer.String()
	s.buffer.Reset()

	if kwd, ok := keywords[str]; ok {
		return token{char: kwd, sval: str}
	}
	return token{char: tName, sval: str}
}

func (s *lexical) numlit() token {
	s.assert(isDigit(s.char))
	var char = s.char
	var expo = "eE"
	s.consume()
	if char == '0' && s.oneof("xX") {
		expo = "pP"
	}
	for {
		if s.oneof(expo) {
			s.oneof("-+")
		}
		if isHexDigit(s.char) {
			s.consume()
		} else if s.char == '.' {
			s.consume()
		} else {
			break
		}
	}
	n, ok := str2num(s.buffer.String())
	if s.buffer.Reset(); !ok {
		s.scanErr("malformed number", tFloat)
	}
	if i, ok := n.(int64); ok {
		return token{char: tInt, ival: i}
	}
	return token{char: tFloat, nval: float64(n.(float64))}
}

func (s *lexical) strlit(delimiter rune) token {
	for s.consume(); s.char != delimiter; {
		switch c := s.char; {
		case c == eof:
			s.scanErr("unfinished string", tEOS)
		case isNewLine(c):
			s.scanErr("unfinished string", tString)
		case c == '\\':
			s.read()
			s.escape()
		default:
			s.consume()
		}
	}
	s.consume()
	str := s.buffer.String()
	s.buffer.Reset()
	return token{char: tString, sval: str[1 : len(str)-1]}
}

func (s *lexical) hexlit(x float64) (n float64, c rune, i int) {
	if c, n = s.char, x; !isHexDigit(c) {
		return
	}
	for {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return
		}
		s.read()
		c, n, i = s.char, n*16.0+float64(c), i+1
	}
}

func (s *lexical) digits() (char rune) {
	for char = s.char; isDigit(char); char = s.char {
		s.consume()
	}
	return
}

func (s *lexical) escape() {
	if esc, ok := escapes[s.char]; ok {
		s.read()
		s.save(esc)
		return
	}
	switch c := s.char; {
	case isNewLine(c):
		s.addline()
		s.save('\n')
	case c == 'x':
		s.save(s.escape16())
	case c == 'z':
		for s.read(); unicode.IsSpace(s.char); {
			if isNewLine(s.char) {
				s.addline()
			} else {
				s.read()
			}
		}
	case c == 'u':
		s.escapeUTF8()
	case c == eof:
		// nothing todo
		return
	default:
		if !isDigit(c) {
			s.escapeErr([]rune{c}, "invalid escape sequence")
		}
		s.save(s.escape10())
	}
}

func (s *lexical) escape10() (r rune) {
	b := make([]rune, 3)

	for c, i := s.char, 0; i < len(b) && isDigit(c); i, c = i+1, s.char {
		b[i], r = c, 10*r+c-'0'
		s.read()
	}
	if r > math.MaxUint8 {
		s.escapeErr(b[:], "decimal escape too large")
	}
	return
}

func (s *lexical) escape16() (r rune) {
	s.read()

	b := [3]rune{'x'}

	for i, c := 1, s.char; i < len(b); i, c, r = i+1, s.char, r<<4+c {
		switch b[i] = c; {
		case 'a' <= c && c <= 'z':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'Z':
			c = c - 'A' + 10
		case isDigit(c):
			c = c - '0'
		default:
			s.escapeErr(b[:i+1], "hexadecimal digit expected")
		}
		s.read()
	}
	return
}

// TODO: this aint right
func (s *lexical) escapeUTF8() {
	s.read() // skip 'u'
	s.checkEsc(s.char == '{', "missing '{'")
	s.read() // skip '{'
	s.checkEsc(isHexDigit(s.char), "hexadecimal digit expected")
	r := hexvalue(s.char)
	s.read() // must have at least one digit
	for isHexDigit(s.char) {
		r = (r << 4) + hexvalue(s.char)
		s.checkEsc(r <= utf8.MaxRune, "UTF-8 value too large")
		s.read()
	}
	s.checkEsc(s.char == '}', "missing '}'")
	s.read() // skip '}'
	for _, r := range []byte(fmt.Sprintf("%c", r)) {
		s.save(rune(r))
	}
}

// skip a sequence '[=*[' or ']=*]'; if sequence is well-formed, return
// its number of '=''s; otherwise, return a negative number (-1 iff there
// are no '=''s after initial bracket).
func (s *lexical) skipSep() int {
	s.assert(s.char == '[' || s.char == ']')
	delim := s.char
	count := 0
	s.consume()
	for s.char == '=' {
		s.consume()
		count++
	}
	if s.char == delim {
		return count
	}
	return -count - 1
}

func (s *lexical) multiline(comment bool, sep int) (str string) {
	if s.consume(); isNewLine(s.char) {
		s.addline()
	}
	for {
		switch s.char {
		case eof:
			if comment {
				s.scanErr("unfinished long comment", tEOS)
			} else {
				s.scanErr("unfinished long string", tEOS)
			}
		case ']':
			if s.skipSep() == sep {
				if s.consume(); !comment {
					str = s.buffer.String()
					str = str[2+sep : len(str)-(2+sep)]
				}
				s.buffer.Reset()
				return
			}
		case '\r':
			s.char = '\n'
			fallthrough
		case '\n':
			s.save(s.char)
			s.addline()
		default:
			if !comment {
				s.save(s.char)
			}
			s.read()
		}
	}
}

// ** reads the first character of file 'f' and skips an optional BOM mark
// ** in its beginning plus its first line if it starts with '#'. Returns
// ** true if it skipped the first line.  In any case, '*cp' has the
// ** first "valid" character of the file (after the optional BOM and
// ** a first-line comment).
// func skipComment(r *bufio.Reader) *bufio.Reader {
// 	if b, err := r.ReadByte(); err == nil {
// 		if b == '#' {
// 			r.ReadString('\n')
// 		} else {
// 			r.UnreadByte()
// 		}
// 	}
// 	return r
// }
func skipComment(s *scanner) *scanner {
	if b, err := s.source.ReadByte(); err == nil {
		if b == '#' {
			s.source.ReadString('\n')
			s.line++
			s.last++
		} else {
			s.source.UnreadByte()
		}
	}
	return s
}

func hexvalue(r rune) rune {
	if isDigit(r) {
		return r - '0'
	}
	return (unicode.ToLower(r) - 'a') + 10
}
