package str

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/Azure/golua/lua"
)

func format(state *lua.State, format string, argc int) string {
	var (
		str strings.Builder
		arg = 2
	)
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			str.WriteByte(format[i])
			continue
		}
		i++                   // at a '%'
		if format[i] == '%' { // "%%" ?
			str.WriteByte('%')
			continue
		}
		if arg > argc {
			state.Errorf("bad argument #%d to 'format' (no value)", arg)
			return ""
		}
		o := fmtOpt(state, format[i:])
		i += len(o) - 2

		str.WriteString(fmtArg(state, o, arg, format[i]))
		arg++
	}
	return str.String()
}

// TODO: 'a', 'A'
func fmtArg(state *lua.State, opt string, arg int, verb byte) string {
	switch verb {
	case 'e', 'E', 'f', 'g', 'G':
		return fmt.Sprintf(opt, state.CheckNumber(arg))
	case 'o', 'x', 'X':
		return fmt.Sprintf(opt, uint(state.CheckInt(arg)))
	case 'i':
		opt = opt[:len(opt)-1] + "d"
		return fmt.Sprintf(opt, state.CheckInt(arg))
	case 'u':
		opt = opt[:len(opt)-1] + "d"
		return fmt.Sprintf(opt, uint(state.CheckInt(arg)))
	case 'c':
		return string([]byte{byte(state.ToInt(arg))})
	case 'd':
		return fmt.Sprintf(opt, state.CheckInt(arg))
	case 'q':
		var (
			q = state.TypeAt(arg) == lua.StringType
			b = new(bytes.Buffer)
		)
		s, ok := state.TryString(arg)
		if !ok && s == "" {
			state.ArgError(arg, "value has no literal form")
		}
		if q {
			b.WriteByte('"')
		}
		for i := 0; i < len(s); i++ {
			switch s[i] {
			case '"', '\\', '\n':
				b.WriteByte('\\')
				b.WriteByte(s[i])
			default:
				if 0x20 <= s[i] && s[i] != 0x7f {
					b.WriteByte(s[i])
				} else if i+1 < len(s) && unicode.IsDigit(rune(s[i+1])) {
					fmt.Fprintf(b, "\\%03d", s[i])
				} else {
					fmt.Fprintf(b, "\\%d", s[i])
				}
			}
		}
		if q {
			b.WriteByte('"')
		}
		opt = opt[:len(opt)-1] + "s"
		return fmt.Sprintf(opt, b.String())
	case 's':
		s := state.ToStringMeta(arg)
		if len(opt) > 2 {
			// If the option has any modifier (flags, width, length),
			// the string argument should not contain embedded zeros.
			state.ArgCheck(strings.Count(s, "\x00") == 0, arg, "strings contains zeros")
		}
		return fmt.Sprintf(opt, s)
	default:
		state.Errorf("invalid option '%%%c' to 'format'", verb)
		return ""
	}
}

// func fmtArg(state *lua.State, arg int, verb rune) string {
// 	switch gofmt := fmt.Sprintf("%%%c", verb); verb {
// 		case 'a', 'A', 'e', 'E', 'f', 'g', 'G':
// 			return fmt.Sprintf(gofmt, state.CheckNumber(arg))
// 		case 'x', 'X':
// 			return fmt.Sprintf(gofmt, uint(state.CheckInt(arg)))
// 		case 'd', 'o':
// 			return fmt.Sprintf(gofmt, state.CheckInt(arg))
// 		case 'c':
// 			return string([]byte{byte(state.ToInt(arg))})
// 		case 'i':
// 			return fmt.Sprintf("%d", state.CheckInt(arg))
// 		case 'u':
// 			return fmt.Sprintf("%d", uint(state.CheckInt(arg)))
// 		case 's', 'q':
// 			return fmt.Sprintf(gofmt, state.ToString(arg))
// 		default:
// 			state.Errorf("invalid option '%%%c' to 'format'", verb)
// 			return ""
// 	}
// }

func fmtOpt(state *lua.State, format string) string {
	index := 0
	digit := func() {
		if unicode.IsDigit(rune(format[index])) {
			index++
		}
	}
	flags := "-+ #0"
	for index < len(format) && strings.ContainsRune(flags, rune(format[index])) {
		index++
	}
	if index >= len(flags) {
		state.Errorf("invalid format (repeated flags)")
		return ""
	}
	digit()
	digit()

	if format[index] == '.' {
		index++
		digit()
		digit()
	}
	if unicode.IsDigit(rune(format[index])) {
		state.Errorf("invalid format (width or precision too long)")
		return ""
	}
	index++
	return "%" + format[:index]
}

// func format(state *lua.State, format string, argc int) (result string) {
// 	var (
// 		end = len(format)
// 		arg = 2
// 	)
// 	FMT: for i := 0; i < end; {
// 		pos := i
// 		for i < end && format[i] != '%' {
// 			i++
// 		}
// 		if i > pos {
// 			result += string(format[pos:i])
// 		}
// 		if i >= end {
// 			break
// 		}
// 		i++
// 		OPT: for ; i < end; i++ {
// 			switch b := format[i]; b {
// 				case '0':
// 					// fmt.Println("opts.zero")
// 					// fmts.opts.zero = !fmts.opts.minus
// 				case '#':
// 					// fmt.Println("opts.sharp")
// 					// fmts.opts.sharp = true
// 				case ' ':
// 					// fmt.Println("opts.space")
// 					// fmts.opts.space = true
// 				case '-':
// 					// fmt.Println("opts.minus")
// 					// fmts.opts.minus = true
// 					// fmts.opts.zero = false
// 				case '+':
// 					// fmt.Println("opts.plus")
// 					// fmts.opts.plus = true
// 				default:
// 					if 'a' <= b && b <= 'z' && arg <= argc {
// 						result += fmtArg(state, arg, rune(b))
// 						arg++
// 						i++
// 						continue FMT
// 					}
// 					break OPT
// 			}
// 		}

// 		// Handle width/precision (2 digits max).
// 		// if i + 1 < end && isDigit(format[i]) {
// 		// 	var ok bool
// 		// 	fmts.opts.width, ok, i = parseNum(format, i, end)
// 		// 	if !ok {
// 		// 		fmts.opts.width = -1
// 		// 	}
// 		// }
// 		// if i + 1 < end && format[i] == '.' {
// 		// 	var ok bool
// 		// 	i++
// 		// 	fmts.opts.prec, ok, i = parseNum(format, i, end)
// 		// 	if !ok {
// 		// 		fmts.opts.prec = -1
// 		// 	}
// 		// }

// 		verb, size := rune(format[i]), 1
// 		if verb >= utf8.RuneSelf {
// 			verb, size = utf8.DecodeRuneInString(format[i:])
// 		}
// 		i += size

// 		switch {
// 			case arg > argc:
// 				state.Errorf("%%!%c(MISSING)", verb)
// 			case verb == '%':
// 				result += "%"
// 			default:
// 				result += fmtArg(state, arg, verb)
// 				arg++
// 		}
// 	}
// 	return
// }

// // 'c', 'd', 'i', 'o', 'u'
// // 'x', 'X'
// // 'a', 'A'
// // 'e', 'E', 'f'
// // 'g', 'G'
// // 'q'
// // 's'
// func fmtArg(state *lua.State, arg int, verb rune) string {
// 	switch gofmt := fmt.Sprintf("%%%c", verb); verb {
// 		case 'a', 'A', 'e', 'E', 'f', 'g', 'G':
// 			return fmt.Sprintf(gofmt, state.CheckNumber(arg))
// 		case 'x', 'X':
// 			return fmt.Sprintf(gofmt, uint(state.CheckInt(arg)))
// 		case 'd', 'o':
// 			return fmt.Sprintf(gofmt, state.CheckInt(arg))
// 		case 'c':
// 			return string([]byte{byte(state.ToInt(arg))})
// 		case 'i':
// 			return fmt.Sprintf("%d", state.CheckInt(arg))
// 		case 'u':
// 			return fmt.Sprintf("%d", uint(state.CheckInt(arg)))
// 		case 's', 'q':
// 			return fmt.Sprintf(gofmt, state.ToString(arg))
// 		default:
// 			state.Errorf("invalid option '%%%c' to 'format'", verb)
// 			return ""
// 	}
// }

// var ErrMissing = errors.New("format: missing arguments for format string")
