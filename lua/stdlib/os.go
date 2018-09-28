package stdlib

import (
	"time"
    "fmt"
    "os"
    "github.com/Azure/golua/lua"
)

var _ = fmt.Println
var _ = os.Exit

//
// Lua Standard Library -- os
//

// OpenOS opens the Lua standard OS library. This library provides operating system facilities.
//
// See https://www.lua.org/manual/5.3/manual.html#6.9
func OpenOS(state *lua.State) int {
	// Create 'os' table
    var osFuncs = map[string]lua.Func{
		"clock":     lua.Func(osClock),
		"date":      lua.Func(osDate),
		"difftime":  lua.Func(osDiffTime),
		"execute":   lua.Func(osExecute),
		"exit":      lua.Func(osExit),
		"getenv":    lua.Func(osGetEnv),
		"remove":    lua.Func(osRemove),
		"rename":    lua.Func(osRename),
		"setlocale": lua.Func(osSetLocale),
		"time":      lua.Func(osTime),
		"tmpname":   lua.Func(osTmpName),
    }
	state.NewTableSize(0, len(osFuncs))
    state.SetFuncs(osFuncs, 0)

	// Return 'os' table
    return 1
}

// os.clock ()
//
// Returns an approximation of the amount in seconds of CPU time used by the program.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.clock
func osClock(state *lua.State) int {
	state.Push(time.Since(epoch).Seconds())
	return 1
}

// os.date ([format [, time]])
//
// Returns a string or a table containing date and time, formatted according to the given string format.
//
// If the time argument is present, this is the time to be formatted (see the os.time function for a description
// of this value). Otherwise, date formats the current time.
//
// If format starts with '!', then the date is formatted in Coordinated Universal Time. After this optional character,
// if format is the string "*t", then date returns a table with the following fields: year, month (1–12), day (1–31),
// hour (0–23), min (0–59), sec (0–61), wday (weekday, 1–7, Sunday is 1), yday (day of the year, 1–366), and isdst
// (daylight saving flag, a boolean). This last field may be absent if the information is not available.
//
// If format is not "*t", then date returns the date as a string, formatted according to the same rules as the
// ISO C function strftime.
//
// When called without arguments, date returns a reasonable date and time representation that depends on the host system
// and on the current locale. (More specifically, os.date() is equivalent to os.date("%c").)
//
// In non-POSIX systems, this function may be not thread safe because of its reliance on C function gmtime and C function
// localtime.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.date
func osDate(state *lua.State) int {
	unimplemented("os.date")
	return 0
}

// os.difftime (t2, t1)
//
// Returns the difference, in seconds, from time t1 to time t2 (where the times are values returned by os.time). In POSIX,
// Windows, and some other systems, this value is exactly t2-t1.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.difftime
func osDiffTime(state *lua.State) int {
	unimplemented("os.difftime")
	return 0
}

// os.execute ([command])
//
// This function is equivalent to the ISO C function system. It passes command to be executed by an operating system shell.
// Its first result is true if the command terminated successfully, or nil otherwise. After this first result the function
// returns a string plus a number, as follows:
//
// "exit": the command terminated normally; the following number is the exit status of the command.
// "signal": the command was terminated by a signal; the following number is the signal that terminated the command.
//
// When called without a command, os.execute returns a boolean that is true if a shell is available.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.execute
func osExecute(state *lua.State) int {
	unimplemented("os.execute")
	return 0
}

// os.exit ([code [, close]])
//
// Calls the ISO C function exit to terminate the host program. If code is true, the returned status is EXIT_SUCCESS;
// if code is false, the returned status is EXIT_FAILURE; if code is a number, the returned status is this number.
// The default value for code is true.
//
// If the optional second argument close is true, closes the Lua state before exiting.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.exit
func osExit(state *lua.State) int {
	unimplemented("os.exit")
	return 0
}

// os.getenv (varname)
//
// Returns the value of the process environment variable varname, or nil if the variable is not defined.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.getenv
func osGetEnv(state *lua.State) int {
	unimplemented("os.getenv")
	return 0
}

// os.remove (filename)
//
// Deletes the file (or empty directory, on POSIX systems) with the given name. If this function fails, it returns nil,
// plus a string describing the error and the error code. Otherwise, it returns true.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.remove
func osRemove(state *lua.State) int {
	unimplemented("os.remove")
	return 0
}

// os.rename (oldname, newname)
//
// Renames the file or directory named oldname to newname. If this function fails, it returns nil, plus a string describing
// the error and the error code. Otherwise, it returns true.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.rename
func osRename(state *lua.State) int {
	unimplemented("os.rename")
	return 0
}

// os.setlocale (locale [, category])
//
// Sets the current locale of the program. locale is a system-dependent string specifying a locale; category is an optional
// string describing which category to change: "all", "collate", "ctype", "monetary", "numeric", or "time"; the default category
// is "all". The function returns the name of the new locale, or nil if the request cannot be honored.
//
// If locale is the empty string, the current locale is set to an implementation-defined native locale. If locale is the string
// "C", the current locale is set to the standard C locale.
//
// When called with nil as the first argument, this function only returns the name of the current locale for the given category.
//
// This function may be not thread safe because of its reliance on C function setlocale.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.setlocale
func osSetLocale(state *lua.State) int {
	state.Push("") // TODO
	return 1
}

// os.time ([table])
//
// Returns the current time when called without arguments, or a time representing the local date and time specified by the
// given table. This table must have fields year, month, and day, and may have fields hour (default is 12), min (default is 0),
// sec (default is 0), and isdst (default is nil). Other fields are ignored. For a description of these fields, see the os.date function.
//
// The values in these fields do not need to be inside their valid ranges. For instance, if sec is -10, it means -10 seconds from the time
// specified by the other fields; if hour is 1000, it means +1000 hours from the time specified by the other fields.
//
// The returned value is a number, whose meaning depends on your system. In POSIX, Windows, and some other systems, this number counts the
// number of seconds since some given start time (the "epoch"). In other systems, the meaning is not specified, and the number returned by
// time can be used only as an argument to os.date and os.difftime.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.time
func osTime(state *lua.State) int {
	if state.IsNoneOrNil(1) { // no args?
		state.Push(time.Now().Unix())
		return 1
	}
	// state.CheckType(1, lua.TableType)
	// state.SetTop(1)
	// var (
	//		sec   = 
	//		min   =
	//		day   =
	//		hour  =
	//		year  = 
	//		month = time.Month(...)
	//		isdst = 
	// )
	// date := time.Date(year, month, day, hour, min, sec, 0, time.Local)
	// state.Push(date.Unix())
	//
	// time.Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time

	fmt.Println("os.time: TODO")
	state.Debug(true)
	return 1
}

// os.tmpname ()
//
// Returns a string with a file name that can be used for a temporary file. The file must be explicitly opened before its use and explicitly
// removed when no longer needed.
//
// In POSIX systems, this function also creates a file with that name, to avoid security risks. (Someone else might create the file with
// wrong permissions in the time between getting the name and creating the file.) You still have to open the file to use it and to remove
// it (even if you do not use it).
//
// When possible, you may prefer to use io.tmpfile, which automatically removes the file when the program ends.
//
// See https://www.lua.org/manual/5.3/manual.html#pdf-os.tmpname
func osTmpName(state *lua.State) int {
	unimplemented("os.tmpname")
	return 0
}

var epoch time.Time // start time.
func init() { epoch = time.Now() }