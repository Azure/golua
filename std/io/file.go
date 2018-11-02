package io

import (
	"fmt"

	"github.com/Azure/golua/lua"
)

const fileTypeName = "FILE*"

// file:close ()
//
// Closes file. Note that files are automatically closed when their handles are
// garbage collected, but that takes an unpredictable amount of time to happen.
//
// When closing a file handle created with io.popen, file:close returns the same
// values returned by os.execute.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:close
func fileClose(state *lua.State) int {
	toFile(state) // make sure argument is an open stream
	return closer(state)
}

// file:flush ()
//
// Saves any written data to file.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:flush
func fileFlush(state *lua.State) int {
	unimplemented("file:flush")
	return 0
}

// file:lines (···)
//
// Returns an iterator function that, each time it is called, reads the file
// according to the given formats. When no format is given, uses "l" as a default.
//
// As an example, the construction
//
//      for c in file:lines(1) do body end
//
// will iterate over all characters of the file, starting at the current position.
// Unlike io.lines, this function does not close the file when the loop ends.
//
// In case of errors this function raises the error, instead of returning an error code.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:lines
func fileLines(state *lua.State) int {
	unimplemented("file:lines")
	return 0
}

// file:read (···)
//
// Reads the file file, according to the given formats, which specify what to read.
// For each format, the function returns a string or a number with the characters
// read, or nil if it cannot read data with the specified format. (In this latter
// case, the function does not read subsequent formats.) When called without formats,
// it uses a default format that reads the next line (see below).
//
// The available formats are:
//
// 	"n": reads a numeral and returns it as a float or an integer, following the
// 		 lexical conventions of Lua. (The numeral may have leading spaces and a
//		 sign.)
//
//		 This format always reads the longest input sequence that is a valid prefix
//		 for a numeral; if that prefix does not form a valid numeral (e.g., an empty
// 	 	 string, "0x", or "3.4e-"), it is discarded and the function returns nil.
//
// 	"a": reads the whole file, starting at the current position.
//		 On end of file, it returns the empty string.
//
// 	"l": reads the next line skipping the end of line, returning nil on end of file.
//		 This is the default format.
//
// 	"L": reads the next line keeping the end-of-line character (if present),
//		 returning nil on end of file.
//
// number: reads a string with up to this number of bytes, returning nil on end
//		   of file. If number is zero, it reads nothing and returns an empty string,
//		   or nil on end of file. The formats "l" and "L" should be used only for
//		   text files.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:read
func fileRead(state *lua.State) int {
	unimplemented("file:read")
	return 0
}

// file:seek ([whence [, offset]])
//
// Sets and gets the file position, measured from the beginning of the file,
// to the position given by offset plus a base specified by the string whence,
// as follows:
// 
// 	"set": base is position 0 (beginning of the file);
// 	"cur": base is current position;
// 	"end": base is end of file;
//
// In case of success, seek returns the final file position, measured in bytes
// from the beginning of the file. If seek fails, it returns nil, plus a string
// describing the error.
// 
// The default value for whence is "cur", and for offset is 0. Therefore, the call
// file:seek() returns the current file position, without changing it; the call
// file:seek("set") sets the position to the beginning of the file (and returns 0);
// and the call file:seek("end") sets the position to the end of the file, and
// returns its size.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:seek
func fileSeek(state *lua.State) int {
	file := toFile(state)

	arg1 := state.OptString(2, "cur")
	arg2 := state.OptNumber(3, 0)

	whence, valid := whence[arg1]
	if !valid {
		panic(fmt.Errorf("bad argument #1 to 'seek' (invalid option)"))
	}
	offset := int64(arg2)
	if float64(offset) != arg2 {
		panic(fmt.Errorf("bad argument #1 to 'seek' (not an integer in proper range)"))
	}
	ret, err := file.Seek(int64(whence), int(offset))
	if err != nil {
		return state.FileResult(err, "")
	}
	state.Push(ret)
	return 1
}

// file:setvbuf (mode [, size])
//
// Sets the buffering mode for an output file.
//
// There are three available modes:
// 
// 	  "no": no buffering; the result of any output operation appears immediately.
// 	"full": full buffering; output operation is performed only when the buffer is
//			full or when you explicitly flush the file (see io.flush).
// 	"line": line buffering; output is buffered until a newline is output or there
//		    is any input from some special files (such as a terminal device).
//
// For the last two cases, size specifies the size of the buffer, in bytes.
// The default is an appropriate size.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:setvbuf
func fileSetvbuf(state *lua.State) int {
	unimplemented("file:setvbuf")
	return 0
}

// file:write (···)
//
// Writes the value of each of its arguments to file.
//
// The arguments must be strings or numbers.
//
// In case of success, this function returns file. Otherwise it
// returns nil plus a string describing the error.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-file:write
func fileWrite(state *lua.State) int {
	unimplemented("file:write")
	return 0
}

// See https://www.lua.org/manual/5.3/manual.html#pdf-file:__gc
func fileGC(state *lua.State) int {
	unimplemented("file:__gc")
	return 0
}

// See https://www.lua.org/manual/5.3/manual.html#pdf-file:__tostring
func fileToString(state *lua.State) int {
	if stream := toStream(state); stream.close == nil {
		state.Push("file (closed)")
	} else {
		state.Push(fmt.Sprintf("file (%p)", stream.file))
	}
	return 1
}