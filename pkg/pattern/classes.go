package pattern

import (
	"unicode"
	"strings"
	"fmt"
)

type class interface {
	matches(rune) bool
}

func (i item) matches(r rune) (match bool) {
	if i.typ == itemText { return i.val == string(r) }

	switch strings.ToLower(i.val) {
		case "a":
			match = isalpha(r)
		case "c":
			match = iscntrl(r)
		case "d":
			match = isdigit(r)
		case "g":
			match = isgraph(r)
		case "l":
			match = islower(r)
		case "p":
			match = ispunct(r)
		case "s":
			match = isspace(r)
		case "u":
			match = isupper(r)
		case "w":
			match = isalnum(r)
		case "x":
			match = isxdigit(r)
		case ".":
			return true
		default:
			panic(fmt.Errorf("unhandled character class %q", i.val))
	}
	// fmt.Printf("class(%s) matches (%s) = %t\n", i.val, string(r), match)
	if strings.ToUpper(i.val) == i.val {
		return !match
	}
	return match
}

func classID(item item) (id string) {
	if item.typ == itemText { return "exact" }

	switch strings.ToLower(item.val) {
	case ".":
		return "single"
	case "a":
		id = "alpha"
	case "c":
		id = "cntrl"
	case "d":
		id = "digit"
	case "g":
		id = "print"
	case "l":
		id = "lower"
	case "p":
		id = "punct"
	case "s":
		id = "space"
	case "u":
		id = "upper"
	case "w":
		id = "alnum"
	case "x":
		id = "digit16"
	default:
		panic(fmt.Errorf("unhandled character class %q", item.val))
	}
	if strings.ToUpper(item.val) == item.val {
		id = "!" + id
	}
	return id
}

func isclass(r rune) bool { return strings.ContainsRune(classes, unicode.ToLower(r)) }
func ispunct(r rune) bool { return isprint(r) && !isalnum(r) && !isspace(r) }
func isalpha(r rune) bool { return islower(r) || isupper(r) }
func isalnum(r rune) bool { return isalpha(r) || isdigit(r) }
func iscntrl(r rune) bool { return unicode.IsControl(r) }
func isspace(r rune) bool { return unicode.IsSpace(r) }
func isprint(r rune) bool { return unicode.IsPrint(r) }
func isdigit(r rune) bool { return '0' <= r && r <= '9' }
func isgraph(r rune) bool { return '!' <= r && r <= '~' }
func islower(r rune) bool { return 'a' <= r && r <= 'z' }
func isupper(r rune) bool { return 'A' <= r && r <= 'Z' }
func isxdigit(r rune) bool { return isdigit(r) || ('a' <= r && r <= 'f') ||  ('A' <= r && r <= 'F') }