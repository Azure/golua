package lua

import "fmt"

// ThreadStatus is a Lua thread status.
type ThreadStatus int

// thread statuses
const (
	ThreadOK ThreadStatus = iota // Thread is in success state 
	ThreadYield 				 // Thread is in suspended state
	ThreadError 				 // Thread finished execution with error
)

// String returns the canoncial string of the thread status.
func (status ThreadStatus) String() string {
	switch status {
		case ThreadOK:
			return "OK"
		case ThreadYield:
			return "YIELD"
		case ThreadError:
			return "ERROR"
	}
	return fmt.Sprintf("unknown thread status %d", status)
}

// Lua execution thread.
type Thread struct {
	*State
}

func (x *Thread) String() string { return "thread" }
func (x *Thread) Type() Type { return ThreadType }