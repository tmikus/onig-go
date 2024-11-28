package onig

import "iter"

// FindMatchesIterator is an iterator over each non-overlapping match in text.
type FindMatchesIterator struct {
	err  error
	r    *Regex
	text string
}

// All iterates over each non-overlapping match in text.
func (c *FindMatchesIterator) All() iter.Seq[*Range] {
	return func(yield func(*Range) bool) {
		for _, match := range c.AllWithIndex() {
			if !yield(match) {
				return
			}
		}
	}
}

// AllWithIndex iterates over each non-overlapping match in text.
func (c *FindMatchesIterator) AllWithIndex() iter.Seq2[int, *Range] {
	return func(yield func(int, *Range) bool) {
		textLength := uint(len(c.text))
		region := NewRegion()
		lastEnd := uint(0)
		index := 0
		var lastMatchEnd *uint
		for {
			if lastEnd > textLength {
				return
			}
			region.Clear()
			_, err := c.r.SearchWithParam(
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
			if !yield(index, pos) {
				return
			}
			index++
		}
	}
}

// Err returns the error, if any, that occurred during iteration.
func (c *FindMatchesIterator) Err() error {
	return c.err
}
