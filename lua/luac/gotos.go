package luac

import "fmt"

func closegoto(ls *lexical, g int, lbl *label) {
	goto_ := ls.gotos[g]
	ls.assert(goto_.label == lbl.label)
	if goto_.level < lbl.level {
		ls.semanticErr(
			fmt.Sprintf("<goto %s> at line %d jumps into the scope of local '%s'",
				goto_.label,
				goto_.line,
				ls.fs.local(goto_.level).Name))
	}
	// fmt.Println(len(ls.fs.instrs), goto_.pc, lbl.pc)
	patchList(ls.fs, goto_.pc, lbl.pc)
	// remove goto from pending list
	// copy(ls.gotos[:g], ls.gotos[g+1:])
	// ls.gotos[len(ls.gotos)-1] = nil
	// ls.gotos = ls.gotos[:len(ls.gotos)-1]

	copy(ls.gotos[g:], ls.gotos[g+1:])
	k := len(ls.gotos) - (g + 1) + g
	n := len(ls.gotos)
	for k < n {
		ls.gotos[k] = nil
		k++
	}
	ls.gotos = ls.gotos[:len(ls.gotos)-(g+1)+g]
}

// findgotos checks whether the new label matches any pending gotos in
// the current block; solves forward jumps.
func findgotos(ls *lexical, lbl *label) {
	for g := ls.fs.block.gotos0; g < len(ls.gotos); {
		if ls.gotos[g].label == lbl.label {
			closegoto(ls, g, lbl)
		} else {
			g++
		}
	}
}

// movegotos exports pending gotos to the outer level, to check them against
// outer labels; if the block being existed has upvalues, and the goto exits
// the scope of any variable (which can be the upvalue), close those variables
// being exited.
func movegotos(fs *function, b *block) {
	// correct pending gotos to current block
	// and try to close it with visible labels.
	for i, ls := b.gotos0, fs.ls; i < len(ls.gotos); {
		if g := ls.gotos[i]; g.level > b.active {
			if b.hasup {
				patchClose(ls.fs, g.pc, b.active)
			}
			g.level = b.active
		}
		if !findlabel(ls, i) {
			i++ // move to next one
		}
	}
}

// findlabel tries to close a goto with existing labels; this solves
// backward jumps.
func findlabel(ls *lexical, g int) bool {
	// check labels in current block for a match.
	for i := ls.fs.block.label0; i < len(ls.labels); i++ {
		var (
			lbl = ls.labels[i]
			gto = ls.gotos[g]
		)
		if lbl.label == gto.label { // correct label?
			if gto.level > lbl.level {
				if ls.fs.block.hasup || len(ls.labels) > ls.fs.block.label0 {
					patchClose(ls.fs, gto.pc, lbl.level)
				}
			}
			closegoto(ls, g, lbl) // close it
			return true
		}
	}
	return false // label not found; cannot close goto
}
