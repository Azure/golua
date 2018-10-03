package binary

import (
	"text/template"
	"strings"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/Azure/golua/lua/ir"
	"github.com/Azure/golua/lua/op"
)

const ProtoTplStr = `
(( .Prefix )) <(( .Source )):(( .Line1 )),(( .Line2 ))>(( $num := len .Codes )) ( ((- $num )) (( $num | plural "instruction" "instructions" -)) )
((- $pclntab := .PcLn -))
((- range $ip, $code := .Codes ))
	(( $ip | add 1 ))	[(( pcln $pclntab $ip -))]
(( end ))

`

var (
	protoTpl = template.Must(template.New("proto").Delims("((", "))").Funcs(funcsMap).Parse(ProtoTplStr))

	funcsMap = map[string]interface{}{
		"plural": plural,
		"add":    add,
		"pcln":   pcln,
	}
)

type templateData struct {
	Prefix string
	Source string
	Line1  uint32
	Line2  uint32
	Codes  []uint32
	PcLn   []uint32
}

func (x *Prototype) String() string {
	var b strings.Builder

	// err := protoTpl.Execute(&b, &templateData{
	// 	Prefix: ifElse(x.SrcPos == 0, "main", "function").(string),
 	// 	Source: ifElse(x.Source == "", "=?", x.Source).(string)[1:],
	// 	Line1:  x.SrcPos,
	// 	Line2:  x.EndPos,
	// 	Codes:  x.Code,
	// 	PcLn:   x.PcLnTab,
	// })
	// if err != nil {
	// 	return err.Error()
	// }

	fmt.Fprintf(&b, "\n%s <%s:%d,%d> (%d instructions)\n",
		ifElse(x.SrcPos == 0, "main", "function").(string),
		ifElse(x.Source == "", "=?", x.Source).(string)[1:],
		x.SrcPos,
		x.EndPos,
		len(x.Code),
	)
	fmt.Fprintln(&b, Disasm(x))

	return b.String()
}

func (x *Chunk) String() string {
	return x.Entry.String()
}

func Disasm(x *Prototype) string {
	var (
		ink = func(x int) int { return x&^(1<<8) } 
		myk = func(x int) int { return -1 - x }
		isk = func(x int) bool { return x&(1<<8)!=0 }
	)

	var b strings.Builder
	for pc := 0; pc < len(x.Code); pc++ {
		var (
			in = ir.Instr(x.Code[pc])
			cc = in.Code()
			ln = "-"
		)

		if len(x.PcLnTab) > 0 {
			ln = fmt.Sprintf("%d", x.PcLnTab[pc])
		}

		fmt.Fprintf(&b, "\t%d\t[%s]\t%s \t", pc+1, ln, cc)

		switch mode, mask := cc.Mode(), cc.Mask(); mode {
			case op.ModeABC:
				fmt.Fprintf(&b, "%d", in.A())
				if !mask.B(op.ArgN) {
					if isk(in.B()) {
						fmt.Fprintf(&b, " %d", myk(ink(in.B())))
					} else {
						fmt.Fprintf(&b, " %d", in.B())
					}
				}
				if !mask.C(op.ArgN) {
					if isk(in.C()) {
						fmt.Fprintf(&b, " %d", myk(ink(in.C())))
					} else {
						fmt.Fprintf(&b, " %d", in.C())
					}
				}
				
			case op.ModeABx:
				fmt.Fprintf(&b, "%d", in.A())
				if mask.B(op.ArgK) {
					fmt.Fprintf(&b, " %d", myk(in.BX()))
				}
				if mask.B(op.ArgU) {
					fmt.Fprintf(&b, " %d", in.BX())
				}
				
			case op.ModeAsBx:
				fmt.Fprintf(&b, "%d %d", in.A(), in.SBX())
				
			case op.ModeAx:
				fmt.Fprintf(&b, "%d", myk(in.AX()))
				
			default:
				panic(fmt.Errorf("unknown instruction format: %d", mode))
		}
		fmt.Fprintln(&b)
	}
	return b.String()	
}

func ifElse(cond bool, arg1, arg2 interface{}) interface{} {
	if cond {
		return arg1
	}
	return arg2
}

func plural(one, many string, count int) string {
	if count == 1 {
		return one
	}
	return many
}

func add(xs ...int) (v int) {
	for _, x := range xs {
		v += x
	}
	return v
}

func pcln(pclntab []uint32, pc int) (line string) {
	line = "-"
	if len(pclntab) > 0 {
		line = fmt.Sprintf("%d", pclntab[pc])
	}
	return line
}

func Print(proto *Prototype) string {
	return spew.Sdump(proto)
}