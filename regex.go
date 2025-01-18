package onig

/*
#cgo CXXFLAGS: -std=c++11
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
	captureNames    []string
	groupIndicesMap map[string][]int
	raw             C.OnigRegex
	syntax          *Syntax
}

// MustCompile creates a new Regex object.
func MustCompile(pattern string) *Regex {
	regex, err := Compile(pattern)
	if err != nil {
		panic(err)
	}
	return regex
}

// MustCompileWithOptions creates a new Regex object with the given options.
func MustCompileWithOptions(pattern string, options RegexOptions) *Regex {
	regex, err := CompileWithOptions(pattern, options)
	if err != nil {
		panic(err)
	}
	return regex
}

// MustCompileWithSyntax creates a new Regex object with the given syntax.
func MustCompileWithSyntax(pattern string, syntax *Syntax) *Regex {
	regex, err := CompileWithSyntax(pattern, syntax)
	if err != nil {
		panic(err)
	}
	return regex
}

// MustCompileWithOptionsAndSyntax creates a new Regex object with the given options and syntax.
func MustCompileWithOptionsAndSyntax(pattern string, options RegexOptions, syntax *Syntax) *Regex {
	regex, err := CompileWithOptionsAndSyntax(pattern, options, syntax)
	if err != nil {
		panic(err)
	}
	return regex
}

// Compile creates a new Regex object.
func Compile(pattern string) (*Regex, error) {
	return CompileWithOptionsAndSyntax(pattern, REGEX_OPTION_NONE, SyntaxDefault)
}

// CompileWithOptions creates a new Regex object with the given options.
func CompileWithOptions(
	pattern string,
	options RegexOptions,
) (*Regex, error) {
	return CompileWithOptionsAndSyntax(pattern, options, SyntaxDefault)
}

// CompileWithSyntax creates a new Regex object with the given syntax.
func CompileWithSyntax(pattern string, syntax *Syntax) (*Regex, error) {
	return CompileWithOptionsAndSyntax(pattern, REGEX_OPTION_NONE, syntax)
}

// CompileWithOptionsAndSyntax creates a new Regex object with the given options and syntax.
func CompileWithOptionsAndSyntax(
	pattern string,
	options RegexOptions,
	syntax *Syntax,
) (*Regex, error) {
	var onigRegex C.OnigRegex
	instance := &Regex{
		groupIndicesMap: map[string][]int{},
		raw:             onigRegex,
		syntax:          syntax,
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
	cText := C.CString(text)
	result := C.searchAllWithParam(
		r.raw,
		cText,
		C.uint(len(text)),
		C.uint(0),
		C.uint(len(text)),
		C.uint(REGEX_OPTION_NONE),
		C.uint(0),
		C.uint(0),
	)
	if result.result == C.ONIG_MISMATCH || result.array == nil {
		return nil, nil
	}
	if result.result < 0 {
		return nil, errorFromCode(result.result)
	}
	length := int(result.array.count)
	regions := make([]*Region, length)
	rawRegions := (*[1 << 30]*C.region)(unsafe.Pointer(result.array.regions))[:length:length]
	for i, rawRegion := range rawRegions {
		regions[i] = newRegion(r, rawRegion)
	}
	C.freeRegionsArray(result.array)
	captures := make([]Captures, len(regions))
	for i, region := range regions {
		captures[i] = Captures{
			Regex:  r,
			Region: region,
			Text:   text,
		}
	}
	return captures, nil
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
	region, err := r.SearchFirstWithParam(
		text,
		0,
		uint(len(text)),
		REGEX_OPTION_NONE,
		0,
		0,
	)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, nil
	}
	return &Captures{
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

// FindMatch returns the first match of the regex in the given text.
// If no match is found, then a nil Range is returned.
func (r *Regex) FindMatch(text string) (*Range, error) {
	region, err := r.SearchFirstWithParam(
		text,
		0,
		uint(len(text)),
		REGEX_OPTION_NONE,
		0,
		0,
	)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, nil
	}
	return region.Pos(0), nil
}

// FindMatches returns a list containing each non-overlapping match in text,
// returning the start and end byte indices with respect to text.
func (r *Regex) FindMatches(text string) ([]*Range, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.find_iter
	cText := C.CString(text)
	result := C.searchAllWithParam(
		r.raw,
		cText,
		C.uint(len(text)),
		C.uint(0),
		C.uint(len(text)),
		C.uint(REGEX_OPTION_NONE),
		C.uint(0),
		C.uint(0),
	)
	if result.result == C.ONIG_MISMATCH || result.array == nil {
		return nil, nil
	}
	if result.result < 0 {
		return nil, errorFromCode(result.result)
	}
	length := int(result.array.count)
	regions := make([]*Region, length)
	rawRegions := (*[1 << 30]*C.region)(unsafe.Pointer(result.array.regions))[:length:length]
	for i, rawRegion := range rawRegions {
		regions[i] = newRegion(r, rawRegion)
	}
	matches := make([]*Range, len(regions))
	for i, region := range regions {
		matches[i] = region.Pos(0)
	}
	C.freeRegionsArrayWithRegions(result.array)
	return matches, nil
}

// GetGroupNumbersForGroupName returns the group numbers for the given group name.
// Returns nil if the group name is not found.
func (r *Regex) GetGroupNumbersForGroupName(groupName string) []int {
	if indices, ok := r.groupIndicesMap[groupName]; ok {
		return indices
	}
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))
	begin := getUChar(cGroupName)
	end := getUCharEnd(cGroupName)
	var nums *C.int
	groupCount := C.onig_name_to_group_numbers(
		r.raw,
		begin,
		end,
		&nums,
	)
	if groupCount <= 0 {
		return nil
	}
	cNums := (*[1 << 30]C.int)(unsafe.Pointer(nums))[:groupCount:groupCount]
	result := make([]int, groupCount)
	for i, num := range cNums {
		result[i] = int(num)
	}
	r.groupIndicesMap[groupName] = result
	return result
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

// MustFindMatch returns the first match of the regex in the given text.
// If no match is found, then a nil Range is returned.
// Compared to FindMatch, this method panics on error.
func (r *Regex) MustFindMatch(text string) *Range {
	match, err := r.FindMatch(text)
	if err != nil {
		panic(err)
	}
	return match
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

// MustSearchFirstWithParam searches pattern in string with match param.
//
// Search for matches the regex in a string. This method will return the index of the first match of the regex within the string,
// if there is one. If from is less than to, then search is performed in forward order, otherwise – in backward order.
//
// For more information see [Match vs Search](https://docs.rs/onig/latest/onig/index.html#match-vs-search)
//
// The encoding of the buffer passed to search in must match the encoding of the regex.
// Compared to SearchWithParam, this method panics on error.
func (r *Regex) MustSearchFirstWithParam(
	text string,
	from uint,
	to uint,
	options RegexOptions,
	maxStackSize uint,
	retryLimitInMatch uint,
) *Region {
	result, err := r.SearchFirstWithParam(text, from, to, options, maxStackSize, retryLimitInMatch)
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

// SearchFirstWithParam searches pattern in string with match param.
//
// Search for matches the regex in a string. This method will return the index of the first match of the regex within the string,
// if there is one. If from is less than to, then search is performed in forward order, otherwise – in backward order.
//
// For more information see [Match vs Search](https://docs.rs/onig/latest/onig/index.html#match-vs-search)
//
// The encoding of the buffer passed to search in must match the encoding of the regex.
func (r *Regex) SearchFirstWithParam(
	text string,
	from uint,
	to uint,
	options RegexOptions,
	maxStackSize uint,
	retryLimitInMatch uint,
) (*Region, error) {
	cText := C.CString(text)
	result := C.searchFirstWithParam(
		r.raw,
		cText,
		C.uint(len(text)),
		C.uint(from),
		C.uint(to),
		C.uint(options),
		C.uint(maxStackSize),
		C.uint(retryLimitInMatch),
	)
	if result.result == C.ONIG_MISMATCH {
		return nil, nil
	}
	if result.result < 0 {
		return nil, errorFromCode(result.result)
	}
	return newRegion(r, result.region), nil
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
