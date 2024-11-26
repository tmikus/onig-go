package onig

type Range struct {
	From int
	To   int
}

func NewRange(from int, to int) *Range {
	return &Range{
		From: from,
		To:   to,
	}
}
