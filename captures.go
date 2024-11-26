package onig

type Captures struct {
	Offset uint
	Region *Region
	Text   string
}

func (c *Captures) All() []string {
	result := make([]string, c.Len())
	for i := 0; i < c.Len(); i++ {
		result[i] = c.At(i)
	}
	return result
}

func (c *Captures) AllPos() []*Range {
	result := make([]*Range, c.Len())
	for i := 0; i < c.Len(); i++ {
		result[i] = c.Pos(i)
	}
	return result
}

func (c *Captures) At(i int) string {
	r := c.Pos(i)
	if r == nil {
		return ""
	}
	return c.Text[r.From:r.To]
}

func (c *Captures) IsEmpty() bool {
	return c.Len() == 0
}

func (c *Captures) Len() int {
	return c.Region.Len()
}

func (c *Captures) Pos(i int) *Range {
	return c.Region.Pos(i)
}
