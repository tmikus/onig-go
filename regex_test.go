package onig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRegex(t *testing.T) {
	regex, err := NewRegex("foo")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
}

func TestRegex_AllCaptures(t *testing.T) {
	regex, err := NewRegex(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	captures, err := regex.AllCaptures("a12b2")
	assert.NoError(t, err)
	assert.Len(t, captures, 2)
	assert.Equal(t, NewRange(1, 3), captures[0].Pos(0))
	assert.Equal(t, NewRange(4, 5), captures[1].Pos(0))
}

func TestRegex_CaptureNames_DefaultSyntax(t *testing.T) {
	regex, err := NewRegex("(he)(l+)(o)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Nil(t, regex.CaptureNames())

	regex, err = NewRegex("(?<foo>foo)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo"}, regex.CaptureNames())

	regex, err = NewRegex("(?<foo>foo)(?<bar>bar)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo", "bar"}, regex.CaptureNames())
}

func TestRegex_CaptureNames_PythonSyntax(t *testing.T) {
	regex, err := NewRegexWithOptions("(he)(l+)(o)", REGEX_OPTION_NONE, SyntaxPython)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Nil(t, regex.CaptureNames())

	regex, err = NewRegexWithOptions("(?P<foo>foo)", REGEX_OPTION_NONE, SyntaxPython)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo"}, regex.CaptureNames())

	regex, err = NewRegexWithOptions("(?P<foo>foo)(?P<bar>bar)", REGEX_OPTION_NONE, SyntaxPython)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo", "bar"}, regex.CaptureNames())
}

func TestRegex_Captures(t *testing.T) {
	regex, err := NewRegex("e(l+)|(r+)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	captures, err := regex.Captures("hello")
	assert.NoError(t, err)
	assert.NotNil(t, captures)
	assert.Equal(t, 3, captures.Len())
	assert.Equal(t, false, captures.IsEmpty())
	assert.Equal(t, NewRange(1, 4), captures.Pos(0))
	assert.Equal(t, NewRange(2, 4), captures.Pos(1))
	assert.Nil(t, captures.Pos(2))
	assert.Equal(t, "ell", captures.At(0))
	assert.Equal(t, "ll", captures.At(1))
	assert.Equal(t, "", captures.At(2))
}
