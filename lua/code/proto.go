package code

type (
	Const interface{}

	Proto struct {
		UpVars []*UpVar // information about the function's upvalues
		Protos []*Proto // functions defined inside the function
		Consts []Const  // constants used by the function
		Instrs []Instr  // function instructions
		Vararg bool     // true variable number of arguments
		ParamN int      // number of fixed parameters
		StackN int      // number of registers need by this function

		// debug information
		Locals []*Local // local variable information
		PcLine []int32  // pc -> line
		Source string   // source name
		SrcPos int      // line defined
		EndPos int      // last line defined
	}

	UpVar struct {
		Name  string
		Stack bool
		Index int
	}

	Local struct {
		Name string
		Live int32
		Dead int32
	}
)
