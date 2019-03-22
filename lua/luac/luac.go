package luac

import (
	"fmt"
	"os"

	"github.com/fibonacci1729/golua/lua/code"
)

var _ = fmt.Println
var _ = os.Exit

var Defaults = &Config{}

type Config struct {}

func Compile(config *Config, file string, src interface{}) (chunk *code.Chunk, err error) {
	defer func(err *error) {
		if r := recover(); r != nil {
			if e, ok := r.(code.Error); ok {
				*err = e
				return
			}
			panic(r)
		}
	}(&err)
	ls := &lexical{scanner: new(scanner).init(file, src)}
	fn := new(parser).mainFunc(ls)
	// TODO: assert
	return &code.Chunk{fn}, err
}

func Bundle(config *Config, files []string) (*code.Chunk, error) {
    bundle := &code.Proto{
		UpVars: []*code.UpVar{
			&code.UpVar{
				Name:  "_ENV",
				Stack: true,
				Index: 0,
			},
		},
		Source: "=(glua)",
		Vararg: true,
		StackN: 2,
	}
	if len(files) == 1 {
		return Compile(config, "@"+files[0], nil)
	}
    for _, file := range files {
    	fn, err := Compile(config, "@"+file, nil)
		if err != nil {
			return nil, err
		}
        if len(fn.Main.UpVars) > 0 {
            fn.Main.UpVars[0].Stack = false
        }
        bundle.Protos = append(bundle.Protos, fn.Main)
    }
	for pc := 0; pc < len(bundle.Protos); pc++ {
		bundle.Instrs = append(bundle.Instrs,
			code.MakeABx(code.CLOSURE, 0, pc),
			code.MakeABC(code.CALL, 0, 1, 1),
		)
		bundle.PcLine = append(bundle.PcLine, 1, 1)
	}
	bundle.Instrs = append(bundle.Instrs, code.MakeABC(code.RETURN, 0, 1, 0))
	bundle.PcLine = append(bundle.PcLine, 1)
	return &code.Chunk{bundle}, nil
}

func Must(chunk *code.Chunk, err error) *code.Chunk {
	if err != nil {
		panic(err)
	}
	return chunk
}