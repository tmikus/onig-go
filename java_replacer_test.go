package onig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJavaRegexReplacer(t *testing.T) {
	r := MustCompileWithSyntax(`hello (.*)`, SyntaxJava)
	assert.Equal(t, `goodbye \`, r.MustReplaceAll("hello world", `goodbye \\`))
	assert.Equal(t, "goodbye hello world", r.MustReplaceAll("hello world", `goodbye $0`))
	assert.Equal(t, "goodbye world", r.MustReplaceAll("hello world", `goodbye $1`))
	assert.Equal(t, `goodbye $1`, r.MustReplaceAll("hello world", `goodbye \$1`))
	assertReplaceError(t, r, "hello world", `goodbye \`)
	assertReplaceError(t, r, "hello world", `goodbye $`)
	assertReplaceError(t, r, "hello world", `goodbye $ 0`)
	assertReplaceError(t, r, "hello world", `goodbye $asdf`)

	// By the looks of it, the Oniguruma library doesn't support named groups in Java syntax.
	//r = MustCompileWithSyntax("hello (?<name>.*)", SyntaxJava)
	//assert.Equal(t, `goodbye world`, r.MustReplaceAll("hello world", `goodbye ${name}`))
	//assert.Equal(t, `goodbye {} world`, r.MustReplaceAll("hello world", `goodbye {} ${name}`))
	//assert.Equal(t, `goodbye world}`, r.MustReplaceAll("hello world", `goodbye ${name}}`))
	//assert.Equal(t, `goodbye $world`, r.MustReplaceAll("hello world", `goodbye \$${name}`))
	//assertReplaceError(t, r, "hello world", `goodbye $ {name}`)
	//assertReplaceError(t, r, "hello world", `goodbye $\{name}`)
	//assertReplaceError(t, r, "hello world", `goodbye ${name\}`)
}
