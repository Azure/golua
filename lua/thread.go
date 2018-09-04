package lua

// Lua execution thread.
type Thread struct {
	*State
}

func (x *Thread) String() string { return "thread" }
func (x *Thread) Type() Type { return ThreadType }