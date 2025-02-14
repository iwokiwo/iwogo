package paginator

type EntryWithID interface {
	GetID() int
}

type Cursor struct {
	PaginateReq
	Start      int `json:"start"`
	End        int `json:"end"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
}

func (c *Cursor) MapCursor(entries []EntryWithID) {
	if len(entries) == 0 {
		c.Start = 0
		return
	}
	c.Start = entries[0].GetID()
	c.End = entries[len(entries)-1].GetID()
	return
}

func (c *Cursor) GetOffset() int {
	return (c.GetPage() - 1) * c.GetLimit()
}
func (c *Cursor) GetLimit() int {
	if c.Limit == 0 {
		c.Limit = 10
	}
	return c.Limit
}
func (c *Cursor) GetPage() int {
	if c.Page == 0 {
		c.Page = 1
	}
	return c.Page
}
func (c *Cursor) GetSort() string {
	if c.Sort == "" {
		c.Sort = "id desc"
	}
	return c.Sort
}
