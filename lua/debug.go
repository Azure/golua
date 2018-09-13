package lua

import (
	"strings"
	"fmt"
	"os"
)

func Debug(state *State) {
	 var b strings.Builder

	fr := state.frame()

    fmt.Fprintf(&b, "\nframe#%d <prev=%p|next=%p>\n", fr.depth, fr.prev, fr.next)
    fmt.Fprintf(&b, "    %s", fr.closure)
	if fr.closure != nil {
		fmt.Fprintf(&b, " (# up = %d)", len(fr.closure.upvals))
	}
	fmt.Fprintln(&b)
    fmt.Fprintf(&b, "            * savedpc = %d\n", fr.pc)
    fmt.Fprintf(&b, "            * returns = %d\n\n", fr.rets)
    // fmt.Fprintf(&b, "            upvalues\n")
    // for i, upval := range fr.closure.upvals {
    //         fmt.Fprintf(&b, "                [%d] %v\n", i, *upval)
    // }
    // fmt.Fprintf(&b, "            end\n\n")
    fmt.Fprintf(&b, "            varargs\n")
    for i, extra := range fr.vararg {
            fmt.Fprintf(&b, "                [%d] %v\n", i, extra)
    }
    fmt.Fprintf(&b, "            end\n")
    fmt.Fprintf(&b, "    end\n\n")
    fmt.Fprintf(&b, "    locals (len=%d, cap=%d, top=%d)\n", len(fr.locals), cap(fr.locals), fr.gettop())
	for i := fr.gettop() - 1; i >= 0; i-- {
		fmt.Fprintf(&b, "        [%d] %v\n", i+1, fr.locals[i])
	}

    fmt.Fprintf(&b, "    end\n")
    fmt.Fprintf(&b, "end\n")
 	
	fmt.Println(b.String())
	os.Exit(1)
}