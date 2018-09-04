package syntax

import (
	"io/ioutil"
	"fmt"
	"io"
)

// Source opens and reads the token source from source.
// If source == nil, the filename is used.
//
// The bytes of source is returned or error.
func Source(filename string, source interface{}) (data []byte, err error) {
    switch source := source.(type) {
        case io.Reader:
            data, err = ioutil.ReadAll(source)
        case string:
            data = []byte(source)
        case []byte:
            data = source
        case nil:
            data, err = ioutil.ReadFile(filename)
        default:
            return nil, fmt.Errorf("invalid source: %T", source)
    }
    if err != nil {
        return nil, fmt.Errorf("reading %s: %v", filename, err)
    }
    return data, nil
}