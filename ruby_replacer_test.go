package onig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRubyRegexReplacer(t *testing.T) {
	r := MustNewRegexWithSyntax(`hello (.*)`, SyntaxRuby)
	assert.Equal(t, "goodbye hello world", r.MustReplaceAll("hello world", "goodbye \\0"))
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", "goodbye \\1"))

	r = MustNewRegexWithSyntax(`hello (?<name>.*)`, SyntaxRuby)
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", "goodbye \\k<name>"))
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", "goodbye \\k<1>"))
	assert.Equal(t, "goodbye world1", r.MustReplaceAll("hello world", "goodbye \\k<1>1"))
	assert.Equal(t, "goodbye hello world", r.MustReplaceAll("hello world", "goodbye \\k<0>"))
	assert.Equal(t, "goodbye hello world0", r.MustReplaceAll("hello world", "goodbye \\k<0>0"))
}
