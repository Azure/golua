package str

import (
	"regexp"
	"fmt"

	"github.com/Azure/golua/lua"
)

var fmtRE = regexp.MustCompile(`%[-+ #0]?([0-9][0-9]?)?(\.[0-9][0-9]?)?[%cdiouxXaAeEfgGqs]`)

func fmtarg(state *lua.State, opt string, arg int) string {
	switch verb := opt[len(opt)-1]; verb {
		case 'a', 'A', 'e', 'E', 'f', 'g', 'G':
			return fmt.Sprintf(opt, state.CheckNumber(arg))
		case 'x', 'X':
			return fmt.Sprintf(opt, uint(state.ToInt(arg)))
		case 'd', 'o':
			return fmt.Sprintf(opt, state.CheckInt(arg))
		case 'c':
			return string([]byte{byte(state.ToInt(arg))})
		case 'i':
			return fmt.Sprintf("%d", state.CheckInt(arg))
		case 'u':
			return fmt.Sprintf("%d", uint(state.CheckInt(arg)))
		case 's', 'q':
			return fmt.Sprintf(opt, state.ToString(arg))
		default:
			panic(fmt.Errorf("invalid option '%%%c' to 'format'", verb))
	}
}

func scanfmt(format string) (opts []string) {
	for len(format) > 0 {
		span := fmtRE.FindStringIndex(format)
		if span == nil {
			return append(opts, format)
		}

		prefix := format[:span[0]]
		option := format[span[0]:span[1]]
		suffix := format[span[1]:]

		if !fmtRE.MatchString(option) {
			panic(fmt.Errorf("invalid format"))
		}

		if prefix != "" {
			opts = append(opts, prefix)
		}

		opts = append(opts, option)
		format = suffix
	}
	return opts
}