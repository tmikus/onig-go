package onig

// RegexReplacer is an object that can be used to replace capture groups in a string.
// It implements the Replace method that takes a Captures object and returns the replaced string
// using the replacement pattern that suits the syntax of the regex.
type RegexReplacer interface {
	Replace(captures *Captures) (string, error)
}

// RegexReplacerFactory is a function that creates a RegexReplacer given a replacement pattern.
type RegexReplacerFactory func(replacement string) RegexReplacer
