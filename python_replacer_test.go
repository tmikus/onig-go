package onig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPythonRegexReplacer(t *testing.T) {
	r := MustNewRegexWithSyntax(`hello (.*)`, SyntaxPython)
	assert.Equal(t, "goodbye hello world", r.MustReplaceAll("hello world", `goodbye \0`))
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", `goodbye \1`))
	assert.Equal(t, `goodbye \1`, r.MustReplaceAll("hello world", `goodbye \\1`))

	r = MustNewRegexWithSyntax(`hello (?P<name>.*)`, SyntaxPython)
	assert.Equal(t, `goodbye \g <name>`, r.MustReplaceAll("hello world", `goodbye \g <name>`))
	assert.Equal(t, `goodbye \ g<name>`, r.MustReplaceAll("hello world", `goodbye \ g<name>`))
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", `goodbye \g<name>`))
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", `goodbye \g<1>`))
	assert.Equal(t, "goodbye world1", r.MustReplaceAll("hello world", `goodbye \g<1>1`))
	assert.Equal(t, "goodbye hello world", r.MustReplaceAll("hello world", `goodbye \g<0>`))
	assert.Equal(t, "goodbye hello world0", r.MustReplaceAll("hello world", `goodbye \g<0>0`))
}
