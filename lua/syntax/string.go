package syntax

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	hfpRE = regexp.MustCompile(`^([0-9a-fA-F]+(\.[0-9a-fA-F]*)?|([0-9a-fA-F]*\.[0-9a-f]+))(p[+\-]?[0-9]+)?$`)
	intRE = regexp.MustCompile(`^[+-]?[0-9]+$|^-?0x[0-9a-fA-F]+$`)
)

func StrToF64(str string) (f64 float64, isFloat bool) {
	return parseFloat(str)
}

func StrToI64(str string) (i64 int64, isInt bool) {
	return parseInteger(str)
}

func parseInteger(str string) (int64, bool) {
	if str = strings.TrimSpace(strings.ToLower(str)); !intRE.MatchString(str) { // not int?
		return 0, false
	}
	if str[0] == '+' {
		str = str[1:]
	}
	if strings.Index(str, "0x") < 0 { // decimal
		i, err := strconv.ParseInt(str, 10, 64)
		return i, err == nil
	}

	// hex
	var sign int64 = 1
	if str[0] == '-' {
		sign = -1
		str = str[3:]
	} else {
		str = str[2:]
	}

	if len(str) > 16 {
		str = str[len(str)-16:] // cut long hex string
	}

	i, err := strconv.ParseUint(str, 16, 64)
	return sign * int64(i), err == nil
}

func parseFloat(str string) (float64, bool) {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	if strings.Contains(str, "nan") || strings.Contains(str, "inf") {
		return 0, false
	}
	if strings.HasPrefix(str, "0x") && len(str) > 2 {
		return parseHexFloat(str[2:])
	}
	if strings.HasPrefix(str, "+0x") && len(str) > 3 {
		return parseHexFloat(str[3:])
	}
	if strings.HasPrefix(str, "-0x") && len(str) > 3 {
		f, ok := parseHexFloat(str[3:])
		return -f, ok
	}
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}

// (0x)ABC.DEFp10
func parseHexFloat(str string) (float64, bool) {
	var i16, f16, p10 float64 = 0, 0, 0

	if !hfpRE.MatchString(str) {
		return 0, false
	}

	// decimal exponent
	if idxOfP := strings.Index(str, "p"); idxOfP > 0 {
		digits := str[idxOfP+1:]
		str = str[:idxOfP]

		var sign float64 = 1
		if digits[0] == '-' {
			sign = -1
		}
		if digits[0] == '-' || digits[0] == '+' {
			digits = digits[1:]
		}

		if len(str) == 0 || len(digits) == 0 {
			return 0, false
		}

		for i := 0; i < len(digits); i++ {
			if x, ok := parseDigit(digits[i], 10); ok {
				p10 = p10*10 + x
			} else {
				return 0, false
			}
		}

		p10 = sign * p10
	}

	// fractional part
	if idxOfDot := strings.Index(str, "."); idxOfDot >= 0 {
		digits := str[idxOfDot+1:]
		str = str[:idxOfDot]
		if len(str) == 0 && len(digits) == 0 {
			return 0, false
		}
		for i := len(digits) - 1; i >= 0; i-- {
			if x, ok := parseDigit(digits[i], 16); ok {
				f16 = (f16 + x) / 16
			} else {
				return 0, false
			}
		}
	}

	// integral part
	for i := 0; i < len(str); i++ {
		if x, ok := parseDigit(str[i], 16); ok {
			i16 = i16*16 + x
		} else {
			return 0, false
		}
	}

	// (i16 + f16) * 2^p10
	f := i16 + f16
	if p10 != 0 {
		f *= math.Pow(2, p10)
	}
	return f, true
}

func parseDigit(digit byte, base int) (float64, bool) {
	if base == 10 || base == 16 {
		switch digit {
		case '0':
			return 0, true
		case '1':
			return 1, true
		case '2':
			return 2, true
		case '3':
			return 3, true
		case '4':
			return 4, true
		case '5':
			return 5, true
		case '6':
			return 6, true
		case '7':
			return 7, true
		case '8':
			return 8, true
		case '9':
			return 9, true
		}
	}
	if base == 16 {
		switch digit {
		case 'a':
			return 10, true
		case 'b':
			return 11, true
		case 'c':
			return 12, true
		case 'd':
			return 13, true
		case 'e':
			return 14, true
		case 'f':
			return 15, true
		}
	}
	return -1, false
}
