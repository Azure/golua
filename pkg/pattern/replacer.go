package pattern

// A Replacer replaces strings that matches a pattern.
type Replacer interface {
	Replace(string)string
}

func (patt *pattern) ReplaceAll(text string, repl Replacer, limit int) (string, int) {
	return patt.replaceAll(text, repl, limit)
}

func (patt *pattern) Replace(text string, repl Replacer) (string, int) {
	return patt.replaceAll(text, repl, 0)
}

func (patt *pattern) replaceAll(text string, repl Replacer, limit int) (string, int) {
	return "", 0
}