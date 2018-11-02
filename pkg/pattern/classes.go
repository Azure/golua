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
			match = unicode.IsLetter(r)
		case "c":
			match = unicode.IsControl(r)
		case "d":
			match = unicode.IsDigit(r)
		case "g":
			match = unicode.IsPrint(r)
		case "l":
			match = unicode.IsLower(r)
		case "p":
			match = unicode.IsPunct(r)
		case "s":
			match = unicode.IsSpace(r)
		case "u":
			match = unicode.IsUpper(r)
		case "w":
			match = unicode.IsLetter(r) || unicode.IsDigit(r)
		// case "x": TODO: hexadecimal characters
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