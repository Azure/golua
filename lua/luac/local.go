package luac

import "fmt"

var _ = fmt.Println

func adjustAssign(fs *function, varsN, exprN int, e *expr) {
	extra := varsN - exprN
	if e.retsX() {
		if extra++; extra < 0 { // include call itself
			extra = 0
		}
		fs.code.returnN(fs, e, extra) // last expr provides the difference
		if extra > 1 {
			fs.reserve(extra - 1)
		}
	} else {
		if e.kind != vvoid {
			fs.code.expr2next(fs, e) // close last expr
		}
		if extra > 0 {
			reg := fs.free
			fs.reserve(extra)
			fs.code.codeNils(fs, reg, extra)
		}
	}
	if exprN > varsN {
		fs.free -= exprN - varsN // remove extra values
	}
}

func adjustLocals(fs *function, varsN int) {
	for fs.active = fs.active + varsN; varsN > 0; varsN-- {
		fs.local(fs.active - varsN).Live = int32(fs.pc)
	}
}

func removeVars(fs *function, toLevel int) {
	count := len(fs.ls.active) - (fs.active - toLevel)
	for fs.active > toLevel {
		fs.local(fs.active - 1).Dead = int32(fs.pc)
		fs.active--
	}
	fs.ls.active = fs.ls.active[:count]
}
