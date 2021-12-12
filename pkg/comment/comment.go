package comment

type Comment struct {
	StreamingPlatform string
	Msg               string
}

// return pretty string
func (c *Comment) ToPString() string {
	if c == nil {
		return ""
	}

	return c.Msg
}
