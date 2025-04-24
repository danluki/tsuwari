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

func WithShardSize(size int32) Option {
	return optionFunc(func(c *Client) {
		c.shardSize = size
	})
}

func WithOnDispatch(
	onDispatch func(channelID string, emoteName string, status string),
) Option {
	return optionFunc(func(c *Client) {
		c.onDispatch = onDispatch
	})
}
