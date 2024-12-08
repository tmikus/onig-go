package onig

/*
#cgo LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lonig

#include "regex.h"

*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

//export goOnigForeachNameCallback
func goOnigForeachNameCallback(
	name *C.UChar,
	nameEnd *C.UChar,
	nGroupNum C.int,
	groupNums *C.int,
	regex C.OnigRegex,
	arg unsafe.Pointer,
) C.int {
	names := (*[]string)(arg)
	nameLength := getPointer(nameEnd) - getPointer(name)
	byteSlice := C.GoBytes(unsafe.Pointer(name), C.int(nameLength))
	*names = append(*names, string(byteSlice))
	return 0
}

// ReplacementFunc is a function that takes the matches Captures and returns the replaced string.
type ReplacementFunc func(capture *Captures) (string, error)

// Regex represents a regular expression.
// It is a wrapper around the Oniguruma regex library.
type Regex struct {
	captureNames []string
	raw          C.OnigRegex
	syntax       *Syntax
}

// MustNewRegex creates a new Regex object.
func MustNewRegex(pattern string) *Regex {
	regex, err := NewRegex(pattern)
	if err != nil {
		panic(err)
	}
	return regex
}

// MustNewRegexWithOptions creates a new Regex object with the given options.
func MustNewRegexWithOptions(pattern string, options RegexOptions) *Regex {
	regex, err := NewRegexWithOptions(pattern, options)
	if err != nil {
		panic(err)
	}
	return regex
}

// MustNewRegexWithSyntax creates a new Regex object with the given syntax.
func MustNewRegexWithSyntax(pattern string, syntax *Syntax) *Regex {
	regex, err := NewRegexWithSyntax(pattern, syntax)
	if err != nil {
		panic(err)
	}
	return regex
}

// MustNewRegexWithOptionsAndSyntax creates a new Regex object with the given options and syntax.
func MustNewRegexWithOptionsAndSyntax(pattern string, options RegexOptions, syntax *Syntax) *Regex {
	regex, err := NewRegexWithOptionsAndSyntax(pattern, options, syntax)
	if err != nil {
		panic(err)
	}
	return regex
}

// NewRegex creates a new Regex object.
func NewRegex(pattern string) (*Regex, error) {
	return NewRegexWithOptionsAndSyntax(pattern, REGEX_OPTION_NONE, SyntaxDefault)
}

// NewRegexWithOptions creates a new Regex object with the given options.
func NewRegexWithOptions(
	pattern string,
	options RegexOptions,
) (*Regex, error) {
	return NewRegexWithOptionsAndSyntax(pattern, options, SyntaxDefault)
}

// NewRegexWithSyntax creates a new Regex object with the given syntax.
func NewRegexWithSyntax(pattern string, syntax *Syntax) (*Regex, error) {
	return NewRegexWithOptionsAndSyntax(pattern, REGEX_OPTION_NONE, syntax)
}

// NewRegexWithOptionsAndSyntax creates a new Regex object with the given options and syntax.
func NewRegexWithOptionsAndSyntax(
	pattern string,
	options RegexOptions,
	syntax *Syntax,
) (*Regex, error) {
	var onigRegex C.OnigRegex
	instance := &Regex{
		raw:    onigRegex,
		syntax: syntax,
	}
	runtime.SetFinalizer(instance, func(regex *Regex) {
		if regex.raw != nil {
			C.onig_free(regex.raw)
			regex.raw = nil
		}
	})
	patternCString := C.CString(pattern)
	patternUString := (*C.UChar)(unsafe.Pointer(patternCString))
	patternEndUString := (*C.UChar)(unsafe.Pointer(uintptr(unsafe.Pointer(patternUString)) + uintptr(len(pattern))))
	defer C.free(unsafe.Pointer(patternCString))
	var onigErrorInfo C.OnigErrorInfo
	result := C.onig_new(
		&instance.raw,
		patternUString,
		patternEndUString,
		C.OnigOptionType(options),
		C.ONIG_ENCODING_UTF8,
		syntax.raw,
		&onigErrorInfo,
	)
	if result != C.ONIG_NORMAL {
		return nil, fmt.Errorf("error creating oniguruma regex: onig_new returned %d", int(result))
	}
	return instance, nil
}

// AllCaptures returns a list of all non-overlapping capture groups matched in text.
// This is operationally the same as FindMatches, except it yields information about submatches.
func (r *Regex) AllCaptures(text string) ([]Captures, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.captures_iter
	iterator := r.AllCapturesIter(text)
	result := make([]Captures, 0)
	for captures := range iterator.All() {
		result = append(result, *captures)
	}
	if iterator.Err() != nil {
		return nil, iterator.Err()
	}
	return result, nil
}

// AllCapturesIter returns an iterator of all non-overlapping capture groups matched in text.
// This is operationally the same as FindMatches, except it yields information about submatches.
func (r *Regex) AllCapturesIter(text string) *CapturesIterator {
	return &CapturesIterator{
		r:    r,
		text: text,
	}
}

// CaptureNames returns a list of the names of all capture groups in the regular expression.
func (r *Regex) CaptureNames() []string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.foreach_name
	// and https://docs.rs/onig/latest/onig/struct.Regex.html#method.capture_names_len
	if r.captureNames == nil {
		var names []string
		C.callOnigForeachName(
			r.raw,
			unsafe.Pointer(&names),
		)
		r.captureNames = names
	}
	return r.captureNames
}

// Captures returns the capture groups corresponding to the leftmost-first match in text.
// Capture group 0 always corresponds to the entire match. If no match is found, then nil is returned.
func (r *Regex) Captures(text string) (*Captures, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.captures
	region := NewRegion()
	match, err := r.SearchWithParam(text, 0, uint(len(text)), REGEX_OPTION_NONE, region, NewMatchParam())
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, nil
	}
	return &Captures{
		Offset: *match,
		Regex:  r,
		Region: region,
		Text:   text,
	}, nil
}

// CreateReplacementFunc creates a ReplacementFunc from the given replacement string.
// The replacement func is created using the syntax's ReplacerFactory if it exists.
func (r *Regex) CreateReplacementFunc(replacement string) ReplacementFunc {
	if r.syntax != nil && r.syntax.ReplacerFactory != nil {
		return r.syntax.ReplacerFactory(replacement).Replace
	}
	return func(captures *Captures) (string, error) {
		return replacement, nil
	}
}

// FindMatches returns a list containing each non-overlapping match in text,
// returning the start and end byte indices with respect to text.
func (r *Regex) FindMatches(text string) ([]*Range, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.find_iter
	iterator := r.FindMatchesIter(text)
	result := make([]*Range, 0)
	for match := range iterator.All() {
		result = append(result, match)
	}
	if iterator.Err() != nil {
		return nil, iterator.Err()
	}
	return result, nil
}

// FindMatchesIter returns an iterator containing each non-overlapping match in text,
// returning the start and end byte indices with respect to text.
func (r *Regex) FindMatchesIter(text string) *FindMatchesIterator {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.find_iter
	return &FindMatchesIterator{
		r:    r,
		text: text,
	}
}

// MustAllCaptures returns a list of all non-overlapping capture groups matched in text.
// This is operationally the same as FindMatches, except it yields information about submatches.
// Compared to AllCaptures, this method panics on error.
func (r *Regex) MustAllCaptures(text string) []Captures {
	captures, err := r.AllCaptures(text)
	if err != nil {
		panic(err)
	}
	return captures
}

// MustCaptures returns the capture groups corresponding to the leftmost-first match in text.
// Capture group 0 always corresponds to the entire match. If no match is found, then nil is returned.
// Compared to Captures, this method panics on error.
func (r *Regex) MustCaptures(text string) *Captures {
	captures, err := r.Captures(text)
	if err != nil {
		panic(err)
	}
	return captures
}

// MustFindMatches returns a list containing each non-overlapping match in text,
// returning the start and end byte indices with respect to text.
// Compared to FindMatches, this method panics on error.
func (r *Regex) MustFindMatches(text string) []*Range {
	matches, err := r.FindMatches(text)
	if err != nil {
		panic(err)
	}
	return matches
}

// MustReplace replaces the leftmost-first match with the replacement provided.
// If no match is found, then a copy of the string is returned unchanged.
// Compared to Replace, this method panics on error.
func (r *Regex) MustReplace(text string, replacement string) string {
	result, err := r.Replace(text, replacement)
	if err != nil {
		panic(err)
	}
	return result
}

// MustReplaceAll replaces all non-overlapping matches in text with the replacement provided.
// This is the same as calling ReplaceN with limit set to 0.
// See the documentation for Replace for details on how to access submatches in the replacement string.
// Compared to ReplaceAll, this method panics on error.
func (r *Regex) MustReplaceAll(text string, replacement string) string {
	result, err := r.ReplaceAll(text, replacement)
	if err != nil {
		panic(err)
	}
	return result
}

// MustReplaceAllFunc replaces all non-overlapping matches in text with the replacement function provided.
// This is the same as calling ReplaceNFunc with limit set to 0.
// See the documentation for Replace for details on how to access submatches in the replacement string.
// Compared to ReplaceAllFunc, this method panics on error.
func (r *Regex) MustReplaceAllFunc(text string, replacement ReplacementFunc) string {
	result, err := r.ReplaceAllFunc(text, replacement)
	if err != nil {
		panic(err)
	}
	return result
}

// MustReplaceFunc replaces the leftmost-first match with the replacement provided.
// The replacement is a function that takes the matches Captures and returns the replaced string.
// If no match is found, then a copy of the string is returned unchanged.
// Compared to ReplaceFunc, this method panics on error.
func (r *Regex) MustReplaceFunc(text string, replacement ReplacementFunc) string {
	result, err := r.ReplaceFunc(text, replacement)
	if err != nil {
		panic(err)
	}
	return result
}

// MustReplaceN replaces at most limit non-overlapping matches in text with the replacement provided.
// If limit is 0, then all non-overlapping matches are replaced.
// See the documentation for Replace for details on how to access submatches in the replacement string.
// Compared to ReplaceN, this method panics on error.
func (r *Regex) MustReplaceN(text string, replacement string, limit int) string {
	result, err := r.ReplaceN(text, replacement, limit)
	if err != nil {
		panic(err)
	}
	return result
}

// MustReplaceNFunc replaces at most limit non-overlapping matches in text with the replacement provided.
// If limit is 0, then all non-overlapping matches are replaced.
// See the documentation for Replace for details on how to access submatches in the replacement string.
// Compared to ReplaceNFunc, this method panics on error.
func (r *Regex) MustReplaceNFunc(text string, replacement ReplacementFunc, limit int) string {
	result, err := r.ReplaceNFunc(text, replacement, limit)
	if err != nil {
		panic(err)
	}
	return result
}

// MustSearchWithParam searches pattern in string with match param.
//
// Search for matches the regex in a string. This method will return the index of the first match of the regex within the string,
// if there is one. If from is less than to, then search is performed in forward order, otherwise – in backward order.
//
// For more information see [Match vs Search](https://docs.rs/onig/latest/onig/index.html#match-vs-search)
//
// The encoding of the buffer passed to search in must match the encoding of the regex.
// Compared to SearchWithParam, this method panics on error.
func (r *Regex) MustSearchWithParam(
	text string,
	from uint,
	to uint,
	options RegexOptions,
	region *Region,
	matchParam *MatchParam,
) *uint {
	result, err := r.SearchWithParam(text, from, to, options, region, matchParam)
	if err != nil {
		panic(err)
	}
	return result
}

// MustSplit returns a list of substrings of text delimited by a match of the regular expression.
// Namely, each element of the iterator corresponds to text that isn’t matched by the regular expression.
// Compared to Split, this method panics on error.
func (r *Regex) MustSplit(text string) []string {
	splits, err := r.Split(text)
	if err != nil {
		panic(err)
	}
	return splits
}

// MustSplitN returns a list of at most `limit` substrings of text delimited by a match of the regular expression.
// A limit of 0 will return no substrings.
// Namely, each element of the iterator corresponds to text that isn’t matched by the regular expression.
// The remainder of the string that is not split will be the last element in the iterator.
// Compared to SplitN, this method panics on error.
func (r *Regex) MustSplitN(text string, limit int) []string {
	splits, err := r.SplitN(text, limit)
	if err != nil {
		panic(err)
	}
	return splits
}

// Replace replaces the leftmost-first match with the replacement provided.
// If no match is found, then a copy of the string is returned unchanged.
func (r *Regex) Replace(text string, replacement string) (string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace
	return r.ReplaceN(text, replacement, 1)
}

// ReplaceAll replaces all non-overlapping matches in text with the replacement provided.
// This is the same as calling ReplaceN with limit set to 0.
// See the documentation for Replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceAll(text string, replacement string) (string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace_all
	return r.ReplaceN(text, replacement, 0)
}

// ReplaceAllFunc replaces all non-overlapping matches in text with the replacement function provided.
// This is the same as calling ReplaceNFunc with limit set to 0.
// See the documentation for Replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceAllFunc(text string, replacement ReplacementFunc) (string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace_all
	return r.ReplaceNFunc(text, replacement, 0)
}

// ReplaceFunc replaces the leftmost-first match with the replacement provided.
// The replacement is a function that takes the matches Captures and returns the replaced string.
// If no match is found, then a copy of the string is returned unchanged.
func (r *Regex) ReplaceFunc(text string, replacement ReplacementFunc) (string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace
	return r.ReplaceNFunc(text, replacement, 1)
}

// ReplaceN replaces at most limit non-overlapping matches in text with the replacement provided.
// If limit is 0, then all non-overlapping matches are replaced.
// See the documentation for Replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceN(text string, replacement string, limit int) (string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replacen
	return r.ReplaceNFunc(text, r.CreateReplacementFunc(replacement), limit)
}

// ReplaceNFunc replaces at most limit non-overlapping matches in text with the replacement provided.
// If limit is 0, then all non-overlapping matches are replaced.
// See the documentation for Replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceNFunc(text string, replacement ReplacementFunc, limit int) (string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replacen
	newText := ""
	captures, err := r.AllCaptures(text)
	if err != nil {
		return "", err
	}
	lastMatch := 0
	for i, capture := range captures {
		if limit > 0 && i >= limit {
			break
		}
		pos := capture.Pos(0)
		if pos == nil {
			continue
		}
		newText += text[lastMatch:pos.From]
		replacedText, err := replacement(&capture)
		if err != nil {
			return "", fmt.Errorf("error replacing text: %w", err)
		}
		newText += replacedText
		lastMatch = pos.To
	}
	newText += text[lastMatch:]
	return newText, nil
}

// SearchWithParam searches pattern in string with match param.
//
// Search for matches the regex in a string. This method will return the index of the first match of the regex within the string,
// if there is one. If from is less than to, then search is performed in forward order, otherwise – in backward order.
//
// For more information see [Match vs Search](https://docs.rs/onig/latest/onig/index.html#match-vs-search)
//
// The encoding of the buffer passed to search in must match the encoding of the regex.
func (r *Regex) SearchWithParam(
	text string,
	from uint,
	to uint,
	options RegexOptions,
	region *Region,
	matchParam *MatchParam,
) (*uint, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.search_with_param
	region.regex = r
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	begin := getUChar(cText)
	end := getUCharEnd(cText)
	limitStart := offsetUChar(begin, int(from))
	limitRange := offsetUChar(begin, int(to))
	if getPointer(limitStart) > getPointer(end) {
		return nil, fmt.Errorf("start of match must be before end")
	}
	if getPointer(limitRange) > getPointer(end) {
		return nil, fmt.Errorf("limit of match should be before end")
	}
	var regionC *C.OnigRegion
	if region != nil {
		regionC = region.raw
	}
	result := C.onig_search_with_param(
		r.raw,
		begin,
		end,
		limitStart,
		limitRange,
		regionC,
		C.uint(options),
		matchParam.raw,
	)
	if result >= 0 {
		resultUint := uint(result)
		return &resultUint, nil
	}
	if result == C.ONIG_MISMATCH {
		return nil, nil
	}
	return nil, errorFromCode(result)
}

// Split returns a list of substrings of text delimited by a match of the regular expression.
// Namely, each element of the iterator corresponds to text that isn’t matched by the regular expression.
func (r *Regex) Split(text string) ([]string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.split
	matches, err := r.FindMatches(text)
	if err != nil {
		return nil, err
	}
	splits := make([]string, 0)
	last := 0
	for _, match := range matches {
		matched := text[last:match.From]
		last = match.To
		splits = append(splits, matched)
	}
	if last < len(text) {
		splits = append(splits, text[last:])
	}
	return splits, nil
}

// SplitN returns a list of at most `limit` substrings of text delimited by a match of the regular expression.
// A limit of 0 will return no substrings.
// Namely, each element of the iterator corresponds to text that isn’t matched by the regular expression.
// The remainder of the string that is not split will be the last element in the iterator.
func (r *Regex) SplitN(text string, limit int) ([]string, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.splitn
	matches, err := r.FindMatches(text)
	if err != nil {
		return nil, err
	}
	splits := make([]string, 0)
	last := 0
	for i, match := range matches {
		if i >= limit-1 {
			break
		}
		matched := text[last:match.From]
		last = match.To
		splits = append(splits, matched)
	}
	if last < len(text) {
		splits = append(splits, text[last:])
	}
	return splits, nil
}
