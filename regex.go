package onig

/*
#cgo CFLAGS: -I/opt/homebrew/include -I/usr/local/include
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

type Regex struct {
	raw C.OnigRegex
}

func NewRegex(pattern string) (*Regex, error) {
	return NewRegexWithOptions(pattern, REGEX_OPTION_NONE, SyntaxDefault)
}

func NewRegexWithOptions(
	pattern string,
	options RegexOptions,
	syntax *Syntax,
) (*Regex, error) {
	var onigRegex C.OnigRegex
	instance := &Regex{
		raw: onigRegex,
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

// AllCaptures returns a list of all non-overlapping capture groups matched in text. This is operationally the same as find_iter (except it yields information about submatches).
func (r *Regex) AllCaptures(text string) ([]Captures, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.captures_iter
	textLength := uint(len(text))
	var lastEnd uint = 0
	var lastMatchEnd *uint
	captures := make([]Captures, 0)
	for {
		if lastEnd > textLength {
			break
		}
		region := NewRegion()
		result, err := r.SearchWithParam(
			text,
			lastEnd,
			uint(len(text)),
			REGEX_OPTION_NONE,
			region,
			NewMatchParam(),
		)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}
		pos := region.Pos(0)
		if pos == nil {
			break
		}
		// Don't accept empty matches immediately following the last match.
		// i.e., no infinite loops please.
		if pos.To == pos.From && lastMatchEnd != nil && *lastMatchEnd == uint(pos.To) {
			offset := 1
			if lastEnd < textLength-1 {
				offset = len(text[lastEnd : lastEnd+1])
			}
			lastEnd += uint(offset)
			continue
		} else {
			toUint := uint(pos.To)
			lastEnd = toUint
			lastMatchEnd = &toUint
		}
		captures = append(captures, Captures{
			Offset: *result,
			Region: region,
			Text:   text,
		})
	}
	return captures, nil
}

// CaptureNames returns a list of the names of all capture groups in the regular expression.
func (r *Regex) CaptureNames() []string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.foreach_name
	// and https://docs.rs/onig/latest/onig/struct.Regex.html#method.capture_names_len
	var names []string
	C.callOnigForeachName(
		r.raw,
		unsafe.Pointer(&names),
	)
	return names
}

// Captures returns the capture groups corresponding to the leftmost-first match in text. Capture group 0 always corresponds to the entire match. If no match is found, then None is returned.
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
		Region: region,
		Text:   text,
	}, nil
}

// FindMatches returns a list containing each non-overlapping match in text, returning the start and end byte indices with respect to text.
func (r *Regex) FindMatches(text string) ([]*Range, error) {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.find_iter
	textLength := uint(len(text))
	region := NewRegion()
	lastEnd := uint(0)
	var lastMatchEnd *uint
	matches := make([]*Range, 0)
	for {
		if lastEnd > textLength {
			break
		}
		region.Clear()
		_, err := r.SearchWithParam(
			text,
			lastEnd,
			textLength,
			REGEX_OPTION_NONE,
			region,
			NewMatchParam(),
		)
		if err != nil {
			return nil, err
		}
		pos := region.Pos(0)
		if pos == nil {
			break
		}
		// Don't accept empty matches immediately following the last match.
		// i.e., no infinite loops please.
		if pos.To == pos.From && lastMatchEnd != nil && *lastMatchEnd == uint(pos.To) {
			offset := 1
			if lastEnd < textLength-1 {
				offset = len(text[lastEnd : lastEnd+1])
			}
			lastEnd += uint(offset)
			continue
		} else {
			toUint := uint(pos.To)
			lastEnd = toUint
			lastMatchEnd = &toUint
		}
		matches = append(matches, pos)
	}
	return matches, nil
}

// Replace replaces the leftmost-first match with the replacement provided. If no match is found, then a copy of the string is returned unchanged.
func (r *Regex) Replace(text string, replacement string) string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace
	panic("not implemented")
}

// ReplaceAll replaces all non-overlapping matches in text with the replacement provided. This is the same as calling replacen with limit set to 0.
// See the documentation for replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceAll(text string, replacement string) string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace_all
	panic("not implemented")
}

// ReplaceAllFunc replaces all non-overlapping matches in text with the replacement function provided. This is the same as calling replacen with limit set to 0.
// See the documentation for replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceAllFunc(text string, replacement string) string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace_all
	panic("not implemented")
}

// ReplaceFunc replaces the leftmost-first match with the replacement provided. The replacement is a function that takes the matches Captures and returns the replaced string.
// If no match is found, then a copy of the string is returned unchanged.
func (r *Regex) ReplaceFunc(text string, replacement func(capture *Captures) string) string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replace
	panic("not implemented")
}

// ReplaceN replaces at most limit non-overlapping matches in text with the replacement provided. If limit is 0, then all non-overlapping matches are replaced.
// See the documentation for replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceN(text string, replacement string, limit int) string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replacen
	panic("not implemented")
}

// ReplaceNFunc replaces at most limit non-overlapping matches in text with the replacement provided. If limit is 0, then all non-overlapping matches are replaced.
// See the documentation for replace for details on how to access submatches in the replacement string.
func (r *Regex) ReplaceNFunc(text string, replacement func(capture *Captures) string, limit int) string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.replacen
	panic("not implemented")
}

// SearchWithParam searches pattern in string with match param
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

// Split returns a list of substrings of text delimited by a match of the regular expression. Namely, each element of the iterator corresponds to text that isn’t matched by the regular expression.
func (r *Regex) Split(text string) []string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.split
	panic("not implemented")
}

// SplitN returns a list of at most `limit` substrings of text delimited by a match of the regular expression. (A limit of 0 will return no substrings.) Namely, each element of the iterator corresponds to text that isn’t matched by the regular expression. The remainder of the string that is not split will be the last element in the iterator.
func (r *Regex) SplitN(text string, limit int) []string {
	// Based on https://docs.rs/onig/latest/onig/struct.Regex.html#method.splitn
	panic("not implemented")
}
