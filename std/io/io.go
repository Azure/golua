package io

import (
    "fmt"
    "os"

    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- io
//

// IO opens the Lua standard IO library. The I/O library provides two different
// styles for file manipulation. The first one uses implicit file handles; that
// is, there are operations to set a default input file and a default output file,
// and all input/output operations are over these default files. The second style
// uses explicit file handles.
// 
// When using implicit file handles, all operations are supplied by table io.
// When using explicit file handles, the operation io.open returns a file handle
// and then all operations are supplied as methods of the file handle.
//
// The table io also provides three predefined file handles with their usual
// meanings from C: io.stdin, io.stdout, and io.stderr. The I/O library never
// closes these files.
//
// Unless otherwise stated, all I/O functions return nil on failure (plus an error
// message as a second result and a system-dependent error code as a third result)
// and some value different from nil on success. In non-POSIX systems, the computation
// of the error message and error code in case of errors may be not thread safe,
// because they rely on the global C variable errno.
//
// See https://www.lua.org/manual/5.3/manual.html#6.8
func Open(state *lua.State) int {
	// Create 'io' table.
	var ioFuncs = map[string]lua.Func{
		"close":   lua.Func(ioClose),
		"flush":   lua.Func(ioFlush),
		"input":   lua.Func(ioInput),
		"lines":   lua.Func(ioLines),
		"open":    lua.Func(ioOpen),
		"output":  lua.Func(ioOutput),
		"popen":   lua.Func(ioPopen),
		"read":    lua.Func(ioRead),
		"tmpfile": lua.Func(ioTmpFile),
		"type":    lua.Func(ioType),
		"write":   lua.Func(ioWrite),
	}
	state.NewTableSize(0, len(ioFuncs))
	state.SetFuncs(ioFuncs, 0)
	createFileMetaTable(state)

	// Create (and set) default standard files.
	createStdFile(state, os.Stdin, "input", "stdin")
	createStdFile(state, os.Stdout, "output", "stdout")
	createStdFile(state, os.Stderr, "", "stderr")

	// Return 'io' table.
	return 1
}

// createFileMetaTable creates the metatable for file handles.
func createFileMetaTable(state *lua.State) {
	var funcs = map[string]lua.Func{
		"close":      lua.Func(fileClose),
		"flush":      lua.Func(fileFlush),
		"lines":      lua.Func(fileLines),
		"read": 	  lua.Func(fileRead),
		"seek": 	  lua.Func(fileSeek),
		"setvbuf":	  lua.Func(fileSetvbuf),
		"write": 	  lua.Func(fileWrite),
		"__gc":   	  lua.Func(fileGC),
		"__tostring": lua.Func(fileToString),
	}
	state.NewMetaTable(fileTypeName)
	state.PushIndex(-1) // push metatable
	state.SetField(-2, "__index") // metatable.__index = metatable
	state.SetFuncs(funcs, 0)
	state.Pop()
}

// createStdFile creates (and sets) the default standard files.
func createStdFile(state *lua.State, file *os.File, field, fname string) {
	newStream(state, file, lua.Func(noClose))
	if field != "" {
		state.PushIndex(-1)
		state.SetField(lua.RegistryIndex, field)
	}
	state.SetField(-2, fname)
}

// io.close ([file])
//
// Equivalent to file:close(). Without a file, closes the default output file.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.close
func ioClose(state *lua.State) int {
	if state.IsNone(1) { // no argument?
		state.GetField(lua.RegistryIndex, "output") // use standard output
	}
	return fileClose(state)
}

// io.flush ()
//
// Equivalent to io.output():flush().
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.flush
func ioFlush(state *lua.State) int {
	unimplemented("io.flush")
	return 0
}

// io.input ([file])
//
// When called with a file name, it opens the named file (in text mode), and sets
// its handle as the default input file. When called with a file handle, it simply
// sets this file handle as the default input file. When called without arguments,
// it returns the current default input file.
//
// In case of errors this function raises the error, instead of returning an error
// code.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.input
func ioInput(state *lua.State) int {
	return getOrSetStdFile(state, "input", "r")
}

// io.output ([file])
//
// Similar to io.input, but operates over the default output file.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.output
func ioOutput(state *lua.State) int {
	return getOrSetStdFile(state, "output", "w")
}

// io.lines ([filename, ···])
// 
// Opens the given file name in read mode and returns an iterator function that
// works like file:lines(···) over the opened file. When the iterator function
// detects the end of file, it returns no values (to finish the loop) and automatically
// closes the file.
// 
// The call io.lines() (with no file name) is equivalent to io.input():lines("*l");
// that is, it iterates over the lines of the default input file. In this case, the
// iterator does not close the file when the loop ends.
// 
// In case of errors this function raises the error, instead of returning an error
// code.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.lines
func ioLines(state *lua.State) int {
	unimplemented("io.lines")
	return 0
}

// io.open (filename [, mode])
//
// This function opens a file, in the mode specified in the string mode.
//
// In case of success, it returns a new file handle.
//
// The mode string can be any of the following:
//
//	 "r": read mode (the default);
//	 "w": write mode;
//	 "a": append mode;
//	"r+": update mode, all previous data is preserved;
//	"w+": update mode, all previous data is erased;
//	"a+": append update mode, previous data is preserved, writing is only allowed
//		  at the end of file.
//
// The mode string can also have a 'b' at the end, which is needed in some systems
// to open the file in binary mode.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.open
func ioOpen(state *lua.State) int {
	filename := state.CheckString(1)

	flags, err := mode2flags(state.OptString(2, "r"))
	if err != nil {
		panic(fmt.Errorf("bad argument #2 to 'open' (invalid mode)"))
	}
	stream := newFile(state)
	stream.file, err = os.OpenFile(filename, flags, 0666)
	if err == nil {
		return 1
	}
	return state.FileResult(err, filename)
}

// io.popen (prog [, mode])
//
// This function is system dependent and is not available on all platforms.
//
// Starts program prog in a separated process and returns a file handle that you
// can use to read data from this program (if mode is "r", the default) or to write
// data to this program (if mode is "w").
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.popen
func ioPopen(state *lua.State) int {
	unimplemented("io.popen")
	return 0
}

// io.read (···)
//
// Equivalent to io.input():read(···).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.read
func ioRead(state *lua.State) int {
	unimplemented("io.read")
	return 0
}

// io.tmpfile ()
//
// In case of success, returns a handle for a temporary file. This file is opened
// in update mode and it is automatically removed when the program ends.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.tmpfile
func ioTmpFile(state *lua.State) int {
	unimplemented("io.tmpfile")
	return 0
}

// io.type (obj)
//
// Checks whether obj is a valid file handle. Returns the string "file" if obj is
// an open file handle, "closed file" if obj is a closed file handle, or nil if
// obj is not a file handle.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.type
func ioType(state *lua.State) int {
	state.CheckAny(1)
	if stream, ok := state.TestUserData(1, fileTypeName).(*stream); ok {
		if stream.close == nil {
			state.Push("closed file")
		} else {
			state.Push("file")
		}
	} else {
		state.Push(nil)
	}
	return 1
}

// io.write (···)
//
// Equivalent to io.output():write(···).
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-io.write
func ioWrite(state *lua.State) int {
	unimplemented("io.write")
	return 0
}

func unimplemented(msg string) { panic(fmt.Errorf(msg)) }