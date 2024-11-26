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

func TestNewRegexCaptureNames_DefaultSyntax(t *testing.T) {
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

func TestNewRegexCaptureNames_PythonSyntax(t *testing.T) {
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
