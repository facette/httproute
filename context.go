package httproute

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}
