// string: 's' (prefix), 'z' (varlen), 'c' (fixed)
// float:  'f', 'd', 'n'
// sint:   'b', 'h', 'l', 'j', 'i'
// uint:   'B', 'H', 'L', 'J', 'I', 'T'
//
// padding, space, config options: "xX <=>!"
//
// values options: "bBhHlLjJTiIfdnczs"
//		i[n], I[n]
//		b, B
//		h, H
//		l, L
//		j, J
//		T
//		f
//		d
//		n
//		z
//		cn
//		s[n]
//
// options "!n", "sn", "in", "In" ; 1 <= n <= 16
//
// Alignment works as follows: For each option, the format gets extra
// padding until the data starts at an offset that is a multiple of
// the minimum between the option size and the maximum alignment; this
// minimum must be a power of 2. Options "c" and "z" are not aligned;
// option "s" follows the alignment of its starting integer.
//
// Any format string starts as if prefixed by "!1=", that is, with maximum
// alignment of 1 (no alignment) and native endianness.
package packer

type (
	Unpacker interface {
		Unpack(state State) error
	}

	Packer interface {
		Pack(state State) error
	}
	
	Option interface {
		Format() (verb rune)
		Width() (width uint)
		Align() (align uint)
	}

	State interface {
		Option() Option
	}
)

// Unpack returns the values packed in string s (see Pack) according to the
// format string. An optional pos marks where to start reading in s (default
// is 1). After the read values, this function also returns the index of the
// first unread byte in s.
func Unpack(format string, values ...interface{}) (int, error) {
	p, err := newState(format)
	if err != nil {
		return 0, err
	}
	return p.Unpack(values...)
}

// Pack returns a binary string containing the values v1, v2, etc. packed
// (that is, serialized in binary form) according to the format string
// format string.
func Pack(format string, values ...interface{}) ([]byte, error) {
	p, err := newState(format)
	if err != nil {
		return nil, err
	}
	return p.Pack(values...)
}

// Size returns the size of a string resulting from string.pack with the
// given format string. The format string cannot have the variable-length
// options 's' or 'z'.
func Size(format string) (int, error) {
	p, err := newState(format)
	if err != nil {
		return 0, err
	}
	return p.Size()
}