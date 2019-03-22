package lua

import (
	"fmt"
	"os"
	"strings"
)

const (
	// Some space for error handling
	errorStackSize = stackMax + 200

	// RegistryIndex is the pseudo-index for the registry table.
	registryIndex = -stackMax - 1000

	// MainThreadIndex is the registry index for the main thread of the main state.
	mainthread = Int(1) //= lua.Int(lua.MAIN_THREAD_INDEX)

	// GlobalsIndex is the registry index for the global environment.
	globals = Int(2) //= lua.Int(lua.GLOBALS_INDEX)

	// Key, in the registry, for table of loaded modules.
	loadedKey = "_LOADED"

	// Key, in the registry, for table or preloaded loaders.
	preloadKey = "_PRELOAD"

	// Initial space allocate for UpValues.
	initNumUps = 5

	// Maximum number of upvalues in a closure (both lua and go).
	// Value must fit in a VM register.
	maxUpVars = 255

	// Limit for table tag-method chains (to avoid loops).
	maxMetaLoop = 2000

	// Limit for table tag-method chains (to avoid loops).
	metaLoopMax = 10

	// Maximum depth for nested Go calls and syntactical nested non-terminals
	// in a program.
	//
	// Value must be < 255.
	maxCalls = 255

	// Maximum valid index and maximum size of stack.
	stackMax = 1000000

	// Minimum Lua stack available to a function.
	stackMin = 20

	// Size allocated for new stacks.
	stackNew = 2 * stackMin

	// Extra stack space to handle metamethod calls and some other extras.
	extraStack = 5

	// Number of list items to accumulate before a SETLIST instruction.
	fieldsPerFlush = 50
)

type Config struct {
	Stdlib func(*Thread) error
	Import Importer
	GoPath Path
	Path   Path
	Trace  bool
	NoEnv  bool
}

func (config *Config) init(rt *runtime) {
	if config.GoPath == "" {
		config.GoPath = Path(GOPATH_DEFAULT)
	}
	if config.Path == "" {
		config.Path = Path(LUAPATH_DEFAULT)
	}

	config.GoPath = Path(envvar(config, GOPATH, string(config.GoPath)))
	config.Path = Path(envvar(config, LUAPATH, string(config.Path)))
	rt.config = config
}

func envvar(config *Config, envVar, defVal string) (path string) {
	versioned := fmt.Sprintf("%s%s", envVar, "_5_3")
	if path = os.Getenv(versioned); path == "" {
		path = os.Getenv(envVar)
	}
	if path == "" {
		path = defVal
	} else {
		path = strings.Replace(path, ";;", "; ;", -1)
		path = strings.Replace(path, "; ;", defVal, -1)
	}
	return path
}
