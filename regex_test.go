package onig

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"strings"
	"testing"
)

func TestNewRegex(t *testing.T) {
	regex, err := Compile("foo")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
}

func TestRegex_AllCaptures(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	captures, err := regex.AllCaptures("a12b2")
	assert.NoError(t, err)
	assert.Len(t, captures, 2)
	assert.Equal(t, NewRange(1, 3), captures[0].Pos(0))
	assert.Equal(t, NewRange(4, 5), captures[1].Pos(0))
}

func TestRegex_CaptureNames_DefaultSyntax(t *testing.T) {
	regex, err := Compile("(he)(l+)(o)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Nil(t, regex.CaptureNames())

	regex, err = Compile("(?<foo>foo)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo"}, regex.CaptureNames())

	regex, err = Compile("(?<foo>foo)(?<bar>bar)")
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo", "bar"}, regex.CaptureNames())
}

func TestRegex_CaptureNames_PythonSyntax(t *testing.T) {
	regex, err := CompileWithSyntax("(he)(l+)(o)", SyntaxPython)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Nil(t, regex.CaptureNames())

	regex, err = CompileWithSyntax("(?P<foo>foo)", SyntaxPython)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo"}, regex.CaptureNames())

	regex, err = CompileWithSyntax("(?P<foo>foo)(?P<bar>bar)", SyntaxPython)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	assert.Equal(t, []string{"foo", "bar"}, regex.CaptureNames())
}

func TestRegex_Captures(t *testing.T) {
	regex, err := Compile("e(l+)|(r+)")
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
	assert.Equal(t, []string{"ell", "ll", ""}, slices.Collect(captures.All()))
}

func TestRegex_FindMatch(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	match, err := regex.FindMatch("a12b2")
	assert.NoError(t, err)
	assert.Equal(t, NewRange(1, 3), match)
}

func TestRegex_FindMatches(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	matches, err := regex.FindMatches("a12b2")
	assert.NoError(t, err)
	assert.Equal(t, []*Range{
		NewRange(1, 3),
		NewRange(4, 5),
	}, matches)
}

func TestRegex_FindMatches_OneZeroLength(t *testing.T) {
	regex, err := Compile(`\d*`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	matches, err := regex.FindMatches("a1b2")
	assert.NoError(t, err)
	assert.Equal(t, []*Range{
		NewRange(0, 0),
		NewRange(1, 2),
		NewRange(3, 4),
	}, matches)
}

func TestRegex_FindMatches_ManyZeroLength(t *testing.T) {
	regex, err := Compile(`\d*`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	matches, err := regex.FindMatches("a1bbb2")
	assert.NoError(t, err)
	assert.Equal(t, []*Range{
		NewRange(0, 0),
		NewRange(1, 2),
		NewRange(3, 3),
		NewRange(4, 4),
		NewRange(5, 6),
	}, matches)
}

func TestRegex_FindMatches_EmptyAfterMatch(t *testing.T) {
	regex, err := Compile(`b|(?=,)`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	matches, err := regex.FindMatches("ba,")
	assert.NoError(t, err)
	assert.Equal(t, []*Range{
		NewRange(0, 1),
		NewRange(2, 2),
	}, matches)
}

func TestRegex_FindMatches_ZeroLengthMatchesJumpsPastMatchLocation(t *testing.T) {
	regex, err := Compile(`\b`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	matches, err := regex.FindMatches("test string")
	assert.NoError(t, err)
	assert.Equal(t, []*Range{
		NewRange(0, 0),
		NewRange(4, 4),
		NewRange(5, 5),
		NewRange(11, 11),
	}, matches)
}

func TestRegex_Replace(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	replaced, err := regex.Replace("a12b2", "X")
	assert.NoError(t, err)
	assert.Equal(t, "aXb2", replaced)
}

func TestRegex_ReplaceAll(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	replaced, err := regex.ReplaceAll("a12b2", "X")
	assert.NoError(t, err)
	assert.Equal(t, "aXbX", replaced)
}

func TestRegex_ReplaceFunc(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	replaced, err := regex.ReplaceFunc("a12b2", func(capture *Captures) (string, error) {
		return "X", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "aXb2", replaced)

	regex, err = Compile(`[a-z]+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	replaced, err = regex.ReplaceFunc("a12b2", func(capture *Captures) (string, error) {
		pos := capture.Pos(0)
		return strings.ToUpper(capture.Text[pos.From:pos.To]), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "A12b2", replaced)
}

func TestRegex_ReplaceAllFunc(t *testing.T) {
	regex, err := Compile(`\d+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	replaced, err := regex.ReplaceAllFunc("a12b2", func(capture *Captures) (string, error) {
		return "X", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "aXbX", replaced)

	regex, err = Compile(`[a-z]+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	replaced, err = regex.ReplaceAllFunc("a12b2", func(capture *Captures) (string, error) {
		pos := capture.Pos(0)
		return strings.ToUpper(capture.Text[pos.From:pos.To]), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "A12B2", replaced)
}

func TestRegex_Split(t *testing.T) {
	regex, err := Compile(`[ \t]+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	splits, err := regex.Split("a b \t  c\td    e")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, splits)
}

func TestRegex_SplitN(t *testing.T) {
	regex, err := Compile(`\W+`)
	assert.NoError(t, err)
	assert.NotNil(t, regex)
	splits, err := regex.SplitN("Hey! How are you?", 3)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Hey", "How", "are you?"}, splits)
}
