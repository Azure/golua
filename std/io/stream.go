package io

import (
	"strings"
	"fmt"
	"os"

	"github.com/Azure/golua/lua"
)

var whence = map[string]int{
	"set": os.SEEK_SET,
	"cur": os.SEEK_CUR,
	"end": os.SEEK_END,
}

type stream struct {
	file  *os.File
	close lua.Func
}

func newStream(state *lua.State, file *os.File, close lua.Func) *stream {
	stream := &stream{file: file, close: close}
	state.Push(stream)
	state.SetMetaTable(fileTypeName)
	return stream
}

func newFile(state *lua.State) *stream {
	return newStream(state, nil, lua.Func(func(state *lua.State) int {
		err := toStream(state).file.Close()
		return state.FileResult(err, "")
	}))
}

func mustOpen(state *lua.State, name, mode string) (file *os.File) {
	flags, err := mode2flags(mode)
	if err == nil {
		file, err = os.OpenFile(name, flags, 0666)
	}
	if err != nil {
		panic(fmt.Errorf("cannot open file '%s' (%s)", name, err.Error()))
	}
	return
}

func getOrSetStdFile(state *lua.State, file, mode string) int {
	if !state.IsNoneOrNil(1) {
		if name := state.ToString(1); name != "" {
			newFile(state).file = mustOpen(state, name, mode)
		} else {
			toFile(state) // check that it's a valid file handle
			state.PushIndex(1)
		}
		state.SetField(lua.RegistryIndex, file)
	}
	// return current value
	state.GetField(lua.RegistryIndex, file)
	return 1
}

func toFile(state *lua.State) *os.File {
	stream := toStream(state)
	if stream.close == nil {
		panic(fmt.Errorf("attempt to use a closed file"))
	}
	if stream.file == nil {
		panic(fmt.Errorf("file is nil"))
	}
	return stream.file
}

func toStream(state *lua.State) *stream {
	return state.CheckUserData(1, fileTypeName).(*stream)
}

func noClose(state *lua.State) int {
	toStream(state).close = noClose
	state.Push(nil)
	state.Push("cannot close standard file")
	return 2
}

func closer(state *lua.State) int {
	stream := toStream(state)
	closer := stream.close
	stream.close = nil
	return closer(state)
}

func mode2flags(mode string) (int, error) {
	if end := len(mode); end > 0 && mode[end-1] == 'b' {
		mode = mode[:end-1]
	}
	switch strings.TrimSpace(mode) {
 		// append update mode, previous data is preserved, writing
		// is only allowed at the end of the file.
		case "a+":
			return os.O_RDWR | os.O_CREATE | os.O_APPEND, nil
 		// update mode, all previous data is preserved.
		case "r+":
			return os.O_RDWR, nil
		// update mode, all previous data is erased.
		case "w+":
			return os.O_RDWR | os.O_CREATE | os.O_TRUNC, nil
		// append mode.
		case "a":
			return os.O_WRONLY | os.O_CREATE | os.O_APPEND, nil
		// read mode (the default).
		case "r":
			return os.O_RDONLY, nil
		// write mode.
		case "w":
			return os.O_WRONLY | os.O_CREATE | os.O_TRUNC, nil
		default:
			return -1, os.ErrInvalid		
	}
}