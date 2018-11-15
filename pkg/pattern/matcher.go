package pattern

type Pattern struct {
	*pattern
}

type capture [2]int

func (cap *capture) start(pos int) { cap[0], cap[1] = pos, -1 }

func (cap *capture) close(pos int) { cap[1] = pos }

type pattern struct {
	caps []capture
	inst []instr
	head bool
	tail bool
}

func (patt *pattern) MatchIndexAll(src string, limit int) (captures [][]int) {
	for start, count := 0, 0; start <= len(src); start++ {
		if end, match := patt.match(src, start, 0); match {
			captures = append(captures, []int{start, end})
			for _, cap := range patt.caps {
				captures[count] = append(
					captures[count],
					cap[0],
					cap[1],
				)
			}
			if count++; limit > 0 && count >= limit {
				break
			}
			start = end
		}
		if patt.head {
			break
		}
	}
	return captures
}

func (patt *pattern) MatchIndex(src string) (captures []int) {
	for sp := 0; sp <= len(src); sp++ {
		if pos, match := patt.match(src, sp, 0); match {
			captures = append(captures, sp, pos)
			for _, cap := range patt.caps {
				captures = append(captures, cap[0], cap[1])
			}
			return captures
		}
		if patt.head {
			break
		}
	}
	return nil
}

func (patt *pattern) MatchAll(src string, limit int) (captures [][]string) {
	for _, caps := range patt.MatchIndexAll(src, limit) {
		var capture []string
		for cap := 0; cap < len(caps); cap += 2 {
			capture = append(capture, src[caps[cap]:caps[cap+1]])
		}
		captures = append(captures, capture)
	}
	return captures
}

func (patt *pattern) Match(src string) (captures []string) {
	if loc := patt.MatchIndex(src); loc != nil {
		for i := 0; i < len(loc); i += 2 {
			captures = append(captures, src[loc[i]:loc[i+1]])
		}
	}
	return captures
}

func (patt *pattern) match(src string, sp, ip int) (pos int, ok bool) {
	for {
		switch inst := patt.inst[ip]; inst.code {
		case opClass:
			traceVM(src, sp, ip, inst)
			if sp >= len(src) || !inst.item.matches(rune(src[sp])) {
				return sp, false
			}
			sp++
			ip++
		case opMatch:
			traceVM(src, sp, ip, inst)
			if patt.tail {
				return sp, sp == len(src)
			}
			return sp, true
		case opStart:
			traceVM(src, sp, ip, inst)
			patt.caps[inst.x].start(sp)
			ip++
		case opClose:
			traceVM(src, sp, ip, inst)
			patt.caps[inst.x].close(sp)
			ip++
		case opSplit:
			traceVM(src, sp, ip, inst)
			if pos, ok = patt.match(src, sp, inst.x); ok {
				return pos, ok
			}
			ip = inst.y
		case opJump:
			ip = inst.x
		}
	}
	panic("unreachable")
}
