package onig

// NewRubyRegexReplacer creates a new RubyRegexReplacer given a replacement pattern.
// See https://ruby-doc.org/core-2.5.1/String.html#method-i-gsub for more information.
func NewRubyRegexReplacer(pattern string) RegexReplacer {
	return newGenericRegexReplacer(pattern, byte('k'))
}
