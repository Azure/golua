package lua

import (
	"os"
)

const (
	// DefaultLuaLibsPrefix is the prefix for the directory containing lua shared libraries.	
	DefaultLuaHomePrefix = DefaultLuaRootPrefix + "share/lua/5.3/"

	// DefaultLuaLibsPrefix is the prefix for the directory containing lua library modules.
	DefaultLuaLibsPrefix = DefaultLuaRootPrefix + "lib/lua/5.3/"

	// DefaultLuaRootPrefix is the prefix for the default lua root directory.
	DefaultLuaRootPrefix = "/usr/local/"
)

const (
	// Name of the environment variable that Lua checks to set
	// package.gopath. This takes precedence over LUA_GOPATH.
	EnvVarLuaGoPath53 = "$LUA_GOPATH_5_3"

	// Name of the environment variable that Lua checks to set
	// package.gopath.
	EnvVarLuaGoPath = "$LUA_GOPATH"

	// Name of the environment variable that Lua checks to set
	// package.path. This takes precedence over LUA_PATH.
	EnvVarLuaPath53 = "$LUA_PATH_5_3"

	// Name of the environment variable that Lua checks to set
	// package.path.
	EnvVarLuaPath = "$LUA_PATH"

	// Versioned equivalent to LUA_INIT.
	EnvVarLuaInit53 = "$LUA_INIT_5_3"

	// When called without option -E, the interpreter checks for
	// an environment variable LUA_INIT_5_3 (or LUA_INIT if the
	// versioned name is not defined) before running any argument.
	// If the variable content has the format @filename, then lua
	// executes the file. Otherwise, lua executes the string itself.
	EnvVarLuaInit = "$LUA_INIT"

	// Name of the environment variable that Lua checks to set
	// the Lua root dir.
	EnvVarLuaRoot = "$LUA_ROOT"
)

const (
	// The default value used to initialize LuaPath if not provided and
	// the environment variables  LUA_PATH_5_3 and LUA_PATH are unset.
	DefaultLuaPath = DefaultLuaHomePrefix + "?.lua;" +  DefaultLuaHomePrefix + "?/init.lua;" +  DefaultLuaLibsPrefix + "?.lua;" +  DefaultLuaLibsPrefix + "?/init.lua;./?.lua;./?/init.lua"
	
	// The default value used to initialize LuaGoPath if not provided and
	// the environment variables LUA_GOPATH_5_3 and LUA_GOPATH are unset.
	DefaultLuaGoPath = DefaultLuaLibsPrefix + "?.so;" + DefaultLuaLibsPrefix + "loadall.so;./?.so"
)

var (
	// Config is a string describing some compile-time configurations for packages.
	// This string consists of a sequence of lines:
	//
	// The first line is the directory separator string. Default is '\' for Windows
	// and '/' for all other systems.
	//
	// The second line is the character that separates templates in a path. Default
	// is ';'.
	//
	// The third line is the string that marks the substitution points in a template.
	// Default is '?'.
	//
	// The fourth line is a string that, in a path in Windows, is replaced by the
	// executable's directory. Default is '!'.
	//
	// The fifth line is a mark to ignore all text after it when building the luaopen_
	// function name. Default is '-'.
	Config = "/\n;\n?\n!\n-"

	// Root is the path to the Lua root directory.
	//
	// On init, Root is set to the value of the environment variable LUA_ROOT;
	// otherwise Root is set to DefaultLuaRootPrefix.
	EnvRoot = expand(EnvVarLuaRoot, DefaultLuaRootPrefix)

	// Home is the path to the directory Lua used by require to search for lua loaders.
	//
	// At start-up, Lua initializes this variable using the value of the environment variable
	// LUA_PATH_5_3 or LUA_PATH. Any ";;" in the environment variables value is replaced by
	// the default path.
	//
	// On init, Home is set to the value of the environment variable LUA_PATH;
	// otherwise Home is set to DefaultLuaPath.
	EnvHome = expand(EnvVarLuaGoPath, DefaultLuaGoPath)

	// Path is the path to the directory used by Lua require to search for a go loader.
	//
	// Lua initializes the LuaGoPath (package.gopath) in the same way it initializes
	// LuaPath (package.path), using the environment variable LUA_GOPATH_5_3 or LUA_GOPATH.
	//
	// On init, Path is set to the value of the environment variable LUA_GOPATH;
	// otherwise Path is set to DefaultLuaGoPath.
	EnvPath = expand(EnvVarLuaPath, DefaultLuaPath)

	// EnvInit is the path to file or string used to initialize Lua.
	//
	// export LUA_INIT=@/path/to/file/init.lua
	EnvInit = expand(EnvVarLuaInit, "")
)

// expand expands the environment variable envvar return the value set in the
// system environment or the provided default.
func expand(envvar string, orElse string) string {
	if value := os.ExpandEnv(envvar); value != "" {
		return value
	}
	return orElse
}

// Option is an optional configuration for a Lua state.
type Option func(*config)

// config holds all configuration for a Lua state.
type config struct {
	errFn func(error)
	check bool
	trace bool
	debug bool
}

// WithChecks returns an Option that instruction a Lua state to perform API checks.
func WithChecks(enable bool) Option {
	return func(cfg *config) {
		cfg.check = enable
	}
}

// WithTrace returns an Option that toggles execution tracing.
func WithTrace(enable bool) Option {
	return func(cfg *config) {
		cfg.trace = enable
	}
}

// WithVerbose returns an Option that toggles verbose logging.
func WithVerbose(enable bool) Option {
	return func(cfg *config) {
		cfg.debug = enable
	}
}

// Mode is a set of flags (or 0). They control where Lua chunk loading is limited
// to binary chunks, text chunks, or both (default).
type Mode uint

const (
	BinaryMode Mode = 1 << iota // Only binary chunks
	TextMode 					// Only text chunks
)