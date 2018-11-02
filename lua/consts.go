package lua

const (
	// The Lua version to which this implementation conforms
	Version = "Lua 5.3"

	// Option for multiple returns in 'lua_pcall' and 'lua_call'
	MultRets = -1
)

const (
	// RegistryIndex is the pseudo-index for the registry table.
	RegistryIndex = -DefaultStackMax - 1000

	// MainThreadIndex is the registry index for the main thread of the main state.
	MainThreadIndex = 1

	// GlobalsIndex is the registry index for the global environment.
	GlobalsIndex = 2

	// Key, in the registry, for table of loaded modules.
	LoadedKey = "_LOADED"

	// Key, in the registry, for table or preloaded loaders.
	PreloadKey = "_PRELOAD"
)

const (
	// Maximum valid index and maximum size of stack.
	DefaultStackMax = 1000000

	// Minimum Lua stack available to a function.
	DefaultStackMin = 20

	// Size allocated for new stacks.
	InitialStackNew = 2 * DefaultStackMin

	// Initial space allocate for UpValues.
	InitialFreeMax = 5

	// Extra stack space to handle metamethod calls and some other extras.
	ExtraStack = 5

	// Maximum number of upvalues in a closure (both lua and go).
	// Value must fit in a VM register.
	MaxUpValues = 255

	// Limit for table tag-method chains (to avoid loops).
	MaxMetaChain = 2000

	// Limit for table tag-method chains (to avoid loops).
	metaLoopMax = 10

	// Maximum depth for nested Go calls and syntactical nested non-terminals
	// in a program.
	//
	// Value must be < 255.
	MaxCalls = 255

	// Number of list items to accumulate before a SETLIST instruction.
	FieldsPerFlush = 50
)