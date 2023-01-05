package twitch

import (
	"time"

	helix "github.com/satont/go-helix/v2"
)

type Token struct {
	tokenExpiresIn *int
	tokenCreatedAt *int64
}

type Twitch struct {
	Token
	Client *helix.Client
}

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	var reset64 int64
	reset64 = int64(lastResponse.GetRateLimitReset())

	currentTime := time.Now().UTC().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64 - currentTime)
		if timeDiff > 0 {
			time.Sleep(timeDiff * time.Second)
		}
	}

	return nil
}

func NewClient(options *helix.Options) *Twitch {
	opts := options
	opts.RateLimitFunc = rateLimitCallback

	client, err := helix.NewClient(opts)
	if err != nil {
		panic(err)
	}

	twitch := Twitch{
		Client: client,
		Token:  Token{},
	}

	go func() {
		for {
			if !twitch.isTokenValid() {
				exp := twitch.Refresh()
				twitch.setExpiresAndCreated(exp)
			}
			time.Sleep(time.Duration(*twitch.Token.tokenExpiresIn) * time.Second)
		}
	}()

	return &twitch
}

func (c *Twitch) setExpiresAndCreated(expiresIn int) {
	exp := expiresIn
	c.Token.tokenExpiresIn = &exp
	t := time.Now().UTC().UnixMilli()
	c.Token.tokenCreatedAt = &t
}

func (c *Twitch) isTokenValid() bool {
	if c.tokenCreatedAt == nil || c.tokenExpiresIn == nil {
		return false
	}

	curr := time.Now().UTC().UnixMilli()
	isExpired := curr > (*c.Token.tokenCreatedAt + int64(*c.Token.tokenExpiresIn))

	return isExpired
}

func (c *Twitch) Refresh() int {
	token, err := c.Client.RequestAppAccessToken([]string{})
	if err != nil {
		panic(err)
	}

	c.Client.SetAppAccessToken(token.Data.AccessToken)
	return token.Data.ExpiresIn
}
