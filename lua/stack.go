package lua

import (
	"strings"
	"fmt"
)

type stack []Value

func (s *stack) reverse(from, to int) {
	for from < to {
		(*s)[from], (*s)[to] = (*s)[to], (*s)[from]
		from++
		to--
	}
}

func (s *stack) insert(i int, v Value) {
	if i--; i >= 0 && i < len(*s) {
		*s = append(*s, nil)
		copy((*s)[i+1:], (*s)[i:])
		(*s)[i] = v
	}
}

func (s *stack) push(v Value) {
	*s = append(*s, v)
}

func (s *stack) popN(n int) []Value {
	vs := make([]Value, n)
	for i := n - 1; i >= 0; i-- {
		vs[i] = s.pop()
	}
	return vs
}

func (s *stack) pop() (v Value) {
	v, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return v
}

func (s *stack) get(i int) Value {
	if i--; i >= 0 && i < len(*s) {
		return (*s)[i]
	}
	return nil
}

func (s *stack) set(i int, v Value) {
	if i--; i >= 0 && i < len(*s) {
		(*s)[i] = v
	}
}

func (s *stack) cut(i int) {
	if i--; i >= 0 || i < len(*s) {
		copy((*s)[i:], (*s)[i+1:])
		(*s)[len(*s)-1] = nil
		*s = (*s)[:len(*s)-1]
	}
}

func (s *stack) top() int { return len(*s) }

func (s *stack) valid(index int) bool { return index > 0 && index <= len(*s) }

// String returns the dump of the stack s.
func (s *stack) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n#### STACK <%d/%d>\n", len(*s), cap(*s))
	for i := 1; i <= s.top(); i++ {
		fmt.Fprintf(&b, "%04d: %v\n", i, s.get(i))
	}
	return b.String()
}

// func (l *State) reallocStack(newSize int) {
// 	l.assert(newSize <= maxStack || newSize == errorStackSize)
// 	l.assert(l.stackLast == len(l.stack)-extraStack)
// 	l.stack = append(l.stack, make([]value, newSize-len(l.stack))...)
// 	l.stackLast = len(l.stack) - extraStack
// 	l.callInfo.next = nil
// 	for ci := l.callInfo; ci != nil; ci = ci.previous {
// 		if ci.isLua() {
// 			top := ci.top
// 			ci.frame = l.stack[top-len(ci.frame) : top]
// 		}
// 	}
// }