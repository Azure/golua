// Patterns in Lua are described by regular strings, which are interpreted
// as patterns by the pattern-matching functions string.find, string.gmatch,
// string.gsub, and string.match. This file implements the syntax and the
// meaning (that is, what they match) of these strings.
//
// See https://www.lua.org/manual/5.3/manual.html#6.4.1
package strs

//import (
	//"strings"
	//"regexp"
	//"fmt"
//)

// A character class is used to represent a set of characters. The following
// combinations are allowed in describing a character class:
//
//     x: (where x is not one of the magic characters ^$()%.[]*+-?) represents the character x itself.
//     .: (a dot) represents all characters.
//    %a: represents all letters.
//    %c: represents all control characters.
//    %d: represents all digits.
//    %g: represents all printable characters except space.
//    %l: represents all lowercase letters.
//    %p: represents all punctuation characters.
//    %s: represents all space characters.
//    %u: represents all uppercase letters.
//    %w: represents all alphanumeric characters.
//    %x: represents all hexadecimal digits.
//    %x: (where x is any non-alphanumeric character) represents the character x. This is
//		  the standard way to escape the magic characters. Any non-alphanumeric character
//		  (including all punctuation characters, even the non-magical) can be preceded by
//		  a '%' when used to represent itself in a pattern.
// [set]: represents the class which is the union of all characters in set. A range of characters
//		  can be specified by separating the end characters of the range, in ascending order, with
//		  a '-'. All classes %x described above can also be used as components in set. All other
//		  characters in set represent themselves. For example, [%w_] (or [_%w]) represents all
//		  alphanumeric characters plus the underscore, [0-7] represents the octal digits, and [0-7%l%-]
//		  represents the octal digits plus the lowercase letters plus the '-' character.
//
// You can put a closing square bracket in a set by positioning it as the first character in the set.
// You can put a hyphen in a set by positioning it as the first or the last character in the set.
// You can also use an escape for both cases.
//
// The interaction between ranges and classes is not defined. Therefore, patterns like [%a-z] or [a-%%] have
// no meaning.
//
// [^set]: represents the complement of set, where set is interpreted as above.
//
// For all classes represented by single letters (%a, %c, etc.), the corresponding uppercase letter represents
// the complement of the class. For instance, %S represents all non-space characters.
//
// The definitions of letter, space, and other character groups depend on the current locale. In particular, the
// class [a-z] may not be equivalent to %l.
// var class = map[byte]string{
// 	'a': "[[:alpha:]]",  'A': "[[:^alpha:]]",  // letters
// 	'c': "[[:cntrl:]]",  'C': "[[:^cntrl:]]",  // control
// 	'd': "[[:digit:]]",  'D': "[[:^digit:]]",  // digits
// 	'g': "[[:graph:]]",  'G': "[[:^graph:]]",  // printable
// 	'l': "[[:lower:]]",  'L': "[[:^lower:]]",  // lower-case
// 	'p': "[[:punct:]]",  'P': "[[:^punct:]]",  // punctuation
// 	's': "[[:space:]]",  'S': "[[:^space:]]",  // spaces
// 	'u': "[[:upper:]]",  'U': "[[:^upper:]]",  // upper-case
// 	'w': "[[:word:]]",   'W': "[[:^word:]]",   // alphanumeric
// 	'x': "[[:xdigit:]]", 'X': "[[:^xdigit:]]", // hex digits
// }

// func Compile(pattern string) (*regexp.Regexp, error) {
// 	var str strings.Builder
// 	var lvl int
// 	for len(pattern) > 0 {
// 		var b0, b1 byte
// 		switch b0, pattern = pattern[0], pattern[1:]; b0 {
// 			case '{', '}', '\\':
// 				str.WriteByte('\\')
// 				str.WriteByte(b0)
// 			case '%':
// 				fmt.Println("TODO: HERE")
// 				_ = b1
// 			case '[':
// 				lvl++
// 				str.WriteByte(b0)
// 			case ']':
// 				lvl--
// 				str.WriteByte(b0)
// 			case '-':
// 				if lvl == 0 {
// 					str.WriteByte(b0)
// 				} else {
// 					str.WriteString("*?")
// 				}
// 			case '^':
// 				if lvl != 0 && str.Len() > 0 {
// 					str.WriteString("\\")
// 				}
// 				str.WriteByte(b0)
// 			case '$':
// 				if len(pattern) > 0 {
// 					str.WriteString("\\")
// 				}
// 				str.WriteByte(b0)
// 			default:
// 				str.WriteByte(b0)
// 		}
// 	}
// 	return nil, nil
// }

// func Gsub(subject, pattern string, replaceN int) (string, int) {
// 	fmt.Printf("subject: %q, pattern: %q, replace: %d\n", subject, pattern, replaceN)
// 	re, err := Compile(pattern)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_ = re
// 	return "", 0
// }