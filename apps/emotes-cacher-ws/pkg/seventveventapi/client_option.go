package seventveventapi

type Option interface {
	apply(*Client)
}

type optionFunc func(*Client)

func (f optionFunc) apply(c *Client) {
	f(c)
}

func WithLogger(logger Logger) Option {
	return optionFunc(func(c *Client) {
		c.logger = logger
	})
}

func WithDebugMode(enabled bool) Option {
	return optionFunc(func(c *Client) {
		c.debugMode.Store(enabled)
	})
}
