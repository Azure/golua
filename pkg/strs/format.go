package strs

// import (
// 	"strings"
// 	"regexp"
// 	"fmt"

// 	"github.com/Azure/golua/lua"
// )

// // %[flags][width][.precision]specifer
// var fmtRE = regexp.MustCompile(`%[ #+-0]?[0-9]*(\.[0-9]+)?[cdeEfgGioqsuxX%]`)

// func Format(format string, args ...lua.Value) (str string, err error) {
// 	if len(format) <= 1 || strings.IndexByte(format, '%') < 0 {
// 		return format, nil
// 	}
// 	arg := 1
// 	for _, opt := range parseFmtOpts(format) {
// 		switch opt[0] {
// 			case '%':
// 				if opt == "%%" {
// 					str += "%"
// 				} else {
// 					if arg++; arg > len(args) {
// 						return "", fmt.Errorf("bad argument #%d to 'format' (no value)", arg)
// 					}
// 					s, err := fmtValue(opt, args[arg])
// 					if err != nil {
// 						return "", err
// 					}
// 					str += s
// 				}
// 			default:
// 				str += opt
// 		}
// 	}
// 	return str, nil
// }

// // Options A, a, E, e, f, G, and g all expect a number as argument.
// // Options c, d, i, o, u, X, and x expect an integer.
// func fmtValue(opt string, val lua.Value) (string, error) {
// 	switch o := opt[len(opt)-1]; o {
// 		//case 'e', 'E':
// 		//case 'g', 'G':
// 		//case 'a', 'A':
// 		case 'd', 'o':
// 		case 'x', 'X':
// 		case 's', 'q':
// 		case 'f':
// 		case 'o':
// 		case 'c':
// 		case 'i':
// 		case 'u':
// 		default:
// 			return "", fmt.Errorf("unhandle format specifier %s", o)
// 	}
// }

// func parseFmtOpts(format string) (opts []string) {
// 	opts = make([]string, 0, len(format)/2)
// 	for format != "" {
// 		pos := fmtRE.FindStringIndex(format)
// 		if pos == nil {
// 			opts = append(opts, format)
// 			break
// 		}
// 		option := format[pos[0]:pos[1]]
// 		prefix := format[:pos[0]]
// 		format = format[pos[1]:]
// 		if prefix != "" {
// 			opts = append(opts, prefix)
// 		}
// 		opts = append(opts, option)
// 	}
// 	return opts
// }