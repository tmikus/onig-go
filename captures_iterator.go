package onig

import "iter"

// CapturesIterator is an iterator over non-overlapping capture groups matched in text.
type CapturesIterator struct {
	err  error
	r    *Regex
	text string
}

// All iterates over all non-overlapping capture groups matched in text.
func (c *CapturesIterator) All() iter.Seq[*Captures] {
	return func(yield func(*Captures) bool) {
		for _, captures := range c.AllWithIndex() {
			if !yield(captures) {
				return
			}
		}
	}
}

// AllWithIndex iterates over all non-overlapping capture groups matched in text.
func (c *CapturesIterator) AllWithIndex() iter.Seq2[int, *Captures] {
	return func(yield func(int, *Captures) bool) {
		textLength := uint(len(c.text))
		var lastEnd uint = 0
		var lastMatchEnd *uint
		index := 0
		for {
			if lastEnd > textLength {
				break
			}
			region := NewRegion()
			result, err := c.r.SearchWithParam(
				c.text,
				lastEnd,
				textLength,
				REGEX_OPTION_NONE,
				region,
				NewMatchParam(),
			)
			if err != nil {
				c.err = err
				return
			}
			if result == nil {
				return
			}
			pos := region.Pos(0)
			if pos == nil {
				return
			}
			// Don't accept empty matches immediately following the last match.
			// i.e., no infinite loops please.
			if pos.To == pos.From && lastMatchEnd != nil && *lastMatchEnd == uint(pos.To) {
				offset := 1
				if lastEnd < textLength-1 {
					offset = len(c.text[lastEnd : lastEnd+1])
				}
				lastEnd += uint(offset)
				continue
			} else {
				toUint := uint(pos.To)
				lastEnd = toUint
				lastMatchEnd = &toUint
			}
			if !yield(index, &Captures{
				Offset: *result,
				Regex:  c.r,
				Region: region,
				Text:   c.text,
			}) {
				return
			}
			index++
		}
	}
}

// Err returns the error, if any, that occurred during iteration.
func (c *CapturesIterator) Err() error {
	return c.err
}
