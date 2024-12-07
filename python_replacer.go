package onig

// NewPythonRegexReplacer creates a new PythonRegexReplacer given a replacement pattern.
// See https://docs.python.org/3/library/re.html#re.sub for more information.
func NewPythonRegexReplacer(pattern string) RegexReplacer {
	return newGenericRegexReplacer(pattern, byte('g'))
}
