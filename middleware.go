package gotdbot

type Middleware func(client *Client, update TlObject, next func() error) error

func (c *Client) Use(mw ...Middleware) {
	if len(mw) == 0 {
		return
	}
	c.mwMu.Lock()
	c.middlewares = append(c.middlewares, mw...)
	c.mwMu.Unlock()
}

func (c *Client) runMiddlewares(update TlObject, next func() error) error {
	c.mwMu.Lock()
	chain := make([]Middleware, len(c.middlewares))
	copy(chain, c.middlewares)
	c.mwMu.Unlock()

	if len(chain) == 0 {
		return next()
	}

	var run func(int) error
	run = func(i int) error {
		if i >= len(chain) {
			return next()
		}
		return chain[i](c, update, func() error {
			return run(i + 1)
		})
	}
	return run(0)
}
