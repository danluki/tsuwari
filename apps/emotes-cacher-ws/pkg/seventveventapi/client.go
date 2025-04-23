package seventveventapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type SubscriptionStatus string

const (
	baseUrl         = "wss://events.7tv.io/v3"
	heartbeatCycles = 3
)

type Subscription struct {
	// Subscription is temporary instance so its requies context inside
	ctx context.Context

	subscriptionLimit atomic.Int32
	curSubscriptions  atomic.Int32
	// Using atomic.Int32 instead of time.Duration for performance
	hearbeatInterval  atomic.Int32
	curHearbeatCycles atomic.Int32
	active            atomic.Bool

	sessionID       string
	heartbeatTicker *time.Ticker
	wsConn          *websocket.Conn
	mu              sync.Mutex
}

func (s *Subscription) Connect(ctx context.Context) error {
	var err error
	s.wsConn, _, err = websocket.DefaultDialer.DialContext(ctx, baseUrl, nil)
	if err != nil {
		return err
	}
	s.active.Store(true)

	return nil
}

func (s *Subscription) Reconnect(ctx context.Context) error {
	if s.wsConn == nil {
		return s.Connect(ctx)
	}

	_ = s.wsConn.Close()
	s.wsConn = nil

	return s.Connect(ctx)
}

func (s *Subscription) Resume(ctx context.Context) error {
	err := s.Reconnect(ctx)
	if err != nil {
		s.active.Store(false)
		return err
	}

	/*
		Success or no will be inside ACK from 7TV, but what data will comes from,
		TODO: probably later we need on ACK find session with this ID and set it to DEAD but how,
		for now ack is not work and probably 7TV will say about this later.
	*/
	s.mu.Lock()
	err = s.wsConn.WriteJSON(ResumeRequest{
		Operation: Resume,
		D: ResumeData{
			SessionID: s.sessionID,
		},
	})
	s.mu.Unlock()
	if err != nil {
		s.active.Store(false)
		return err
	}

	return nil
}

type Client struct {
	logger Logger

	debugMode atomic.Bool

	shard         int32
	subscriptions []*Subscription
	mu            sync.RWMutex
}

func NewFx() *Client {
	return NewClient(WithDebugMode(true))
}

func NewClient(opts ...Option) *Client {
	client := &Client{}

	for _, opt := range opts {
		opt.apply(client)
	}

	if client.debugMode.Load() && client.logger == nil {
		client.logger = NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	return client
}

func (c *Client) ReconnectSubscription(ctx context.Context, subscription *Subscription) error {
	/*
		TODO: more verbose logs
	*/
	if c.debugMode.Load() {
		c.logger.Debug("Trying to reconnect to 7TV WS")
	}

	err := subscription.Reconnect(ctx)
	if err != nil {
		if c.debugMode.Load() {
			c.logger.Error(
				"Error when trying to reconnect to 7TV WS",
				slog.Any("err", err),
			)
		}

		return err
	}

	if c.debugMode.Load() {
		c.logger.Debug("Reconnected to 7TV WS")
	}

	return nil
}

func (c *Client) AddListener(ctx context.Context) error {
	subscription := &Subscription{
		ctx: ctx,
	}

	err := subscription.Connect(context.WithoutCancel(ctx))
	if err != nil {
		if c.debugMode.Load() {
			c.logger.Error("Error when trying to connect to 7TV WS", slog.Any("err", err))
		}
		return err
	}
	c.subscriptions = append(c.subscriptions, subscription)

	go func() {
		defer func() {
			subscription.heartbeatTicker.Stop()
		}()

		select {
		case <-subscription.ctx.Done():
			if c.debugMode.Load() {
				c.logger.Debug("Listener context is done")
			}
			return
		default:
			for {
				_, rawMsg, err := subscription.wsConn.ReadMessage()
				if err != nil {
					if c.debugMode.Load() {
						c.logger.Error("Error when reading message from WS", slog.Any("err", err))
					}
					subscription.active.Store(false)
					continue
				}

				var eventApiWsMsg WsMessage
				if err = json.Unmarshal(rawMsg, &eventApiWsMsg); err != nil {
					if c.debugMode.Load() {
						c.logger.Error(
							"Error when marshalling message from WS",
							slog.Any("err", err),
						)
					}
					subscription.active.Store(false)
					continue
				}

				if c.debugMode.Load() {
					c.logger.Info("New message from 7TV WS", slog.Any("msg", eventApiWsMsg))
				}

				switch ServerOpcode(eventApiWsMsg.Operation) {
				case Dispatch:
					/*
						TODO: Looks for emotes-sets update, delete, create, and update information about this in redis
					*/
				case Hello:
					var body HelloRequest
					body, err = mapToHelloRequest(eventApiWsMsg.D)
					if err != nil {
						if c.debugMode.Load() {
							c.logger.Error(
								"Wrong body comes for HelloRequest",
								slog.Any("err", err),
							)
						}
						subscription.active.Store(false)
					}

					subscription.subscriptionLimit.Store(body.SubscriptionLimit)
					subscription.hearbeatInterval.Store(body.HeartbeatInterval)

					subscription.mu.Lock()
					subscription.sessionID = body.SessionID
					subscription.mu.Unlock()

					subscription.heartbeatTicker = time.NewTicker(
						time.Duration(subscription.hearbeatInterval.Load()) * time.Millisecond,
					)
					subscription.active.Store(true)

					go func() {
						for range subscription.heartbeatTicker.C {
							if c.debugMode.Load() {
								c.logger.Debug(
									"Heartbeat timeout cycle",
									slog.Any("currCycles", subscription.curHearbeatCycles.Load()),
								)
							}

							subscription.curHearbeatCycles.Add(1)

							if subscription.curHearbeatCycles.Load() > heartbeatCycles {
								subscription.curHearbeatCycles.Store(0)

								if c.debugMode.Load() {
									c.logger.Debug("Heartbeat timeout, reconnecting")
								}

								subscription.active.Store(false)
							}
						}
					}()
				case Heartbeat:
					if c.debugMode.Load() {
						c.logger.Debug(
							"Heartbeat from 7TV",
							slog.Any("currCycles", subscription.curHearbeatCycles.Load()),
						)
					}
					subscription.curHearbeatCycles.Store(0)
				case Reconnect:
					err = subscription.Resume(context.TODO())
					if err != nil {
						if c.debugMode.Load() {
							c.logger.Debug(
								"Cannot reconnect to 7TV",
								slog.Any("currCycles", subscription.curHearbeatCycles.Load()),
							)
						}
					}
					subscription.active.Store(true)
				case Ack:
					var body AckRequest
					body, err := mapToAckRequest(eventApiWsMsg.D)
					if err != nil {
						if c.debugMode.Load() {
							c.logger.Error(
								"Wrong body comes for AckRequest",
								slog.Any("err", err),
							)
						}
					}

					if body.Command == "RESUME" {
						successRaw, ok := body.Data["success"]
						success := false
						if ok {
							if val, valid := successRaw.(bool); valid {
								success = val
							}
						}
						if success {
							subscription.active.Store(true)

							if c.debugMode.Load() {
								c.logger.Debug(
									"RESUMED session",
									slog.String("session_id", subscription.sessionID),
									slog.Bool("status", success),
								)
							}
						} else {
							subscription.active.Store(false)
							c.logger.Debug(
								"RESUMED session error",
								slog.String("session_id", subscription.sessionID),
								slog.Bool("status", success),
							)
						}
					}
				case Error:
					if c.debugMode.Load() {
						c.logger.Debug(
							"Error occured on 7TV",
						)
					}
					subscription.active.Store(false)
				case EndOfStream:
					var body EndOfStreamRequest
					body, err = mapToEndOfStreamRequest(eventApiWsMsg.D)
					if err != nil {
						if c.debugMode.Load() {
							c.logger.Error(
								"Wrong body comes for EndOfStream",
								slog.Any("err", err),
							)
						}
						subscription.active.Store(false)
					}

					if c.debugMode.Load() {
						c.logger.Debug(
							"7TV error message: ",
							slog.Any("code", body.Code),
							slog.String("msg", body.Message),
						)
					}

					switch body.Code {
					case ServerError:
						err = subscription.Resume(context.TODO())
						if err != nil {
							if c.debugMode.Load() {
								c.logger.Debug(
									"Cannot reconnect to 7TV",
									slog.Any("currCycles", subscription.curHearbeatCycles.Load()),
								)
							}
						}
						subscription.active.Store(true)
					case UnknownOperation:
						panic(ErrLibraryBadImplementation)
					case AuthFailure:
						panic(ErrLibraryBadImplementation)
					case AlreadyIdentified:
						panic(ErrLibraryBadImplementation)
					case InvalidPayload:
					case RateLimited:
						subscription.active.Store(false)
						c.mu.Lock()
						_ = subscription.wsConn.Close()
						c.mu.Unlock()
					case Restart:
						err = subscription.Resume(context.TODO())
						if err != nil {
							if c.debugMode.Load() {
								c.logger.Debug(
									"Cannot reconnect to 7TV",
									slog.Any("currCycles", subscription.curHearbeatCycles.Load()),
								)
							}
						}
						subscription.active.Store(true)
					case Maintenance:
						/*
							Need's to implement some hard logic here, for now just sleep for some time and panic
						*/
						if c.debugMode.Load() {
							c.logger.Debug(
								"7TV maintenance, waiting for 3 minues and reconnect",
							)
						}
						time.Sleep(3 * time.Minute)
						err = subscription.Resume(context.TODO())
						if err != nil {
							if c.debugMode.Load() {
								c.logger.Debug(
									"Cannot reconnect to 7TV",
								)
							}
						}
						subscription.active.Store(false)
					case Timeout:
						if c.debugMode.Load() {
							c.logger.Debug(
								"Client timeout",
							)
						}
						err = subscription.Resume(context.TODO())
						if err != nil {
							if c.debugMode.Load() {
								c.logger.Debug(
									"Cannot reconnect to 7TV",
									slog.Any("currCycles", subscription.curHearbeatCycles.Load()),
								)
							}
						}
						subscription.active.Store(true)
					case AlreadySubscribed:
						panic(ErrAlreadySubscribed)
					case NotSubscribed:
						panic(ErrNotSubscribed)
					case InsufficientPrivilege:
						panic(InsufficientPrivilege)
					}
				default:
					panic(ErrUnsupportedOperation)
				}
			}
		}
	}()

	return nil
}

func (c *Client) SubscribeToEmoteSetPatch(
	ctx context.Context,
	id string,
) error {
	subscription, err := c.GetCurrentSubscriptionOrCreateNew(ctx)
	if err != nil {
		return err
	}

	subscription.mu.Lock()
	err = subscription.wsConn.WriteJSON(SubscribeRequest{
		Operation: Subscribe,
		D: SubscribeData{
			Type: PatchEmoteSet,
			Condition: map[string]string{
				"object_id": id,
			},
		},
	})
	subscription.mu.Unlock()
	subscription.curSubscriptions.Add(1)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCurrentSubscriptionOrCreateNew(ctx context.Context) (*Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.subscriptions) == 0 {
		err := c.AddListener(ctx)
		if err != nil {
			return nil, err
		}

		return c.subscriptions[c.shard], nil
	}

	subscription := c.subscriptions[c.shard]
	if c.debugMode.Load() {
		c.logger.Debug(
			"Information output",
			slog.Any("subscription", subscription.curSubscriptions.Load()),
		)
	}

	if subscription.curSubscriptions.Load() > subscription.subscriptionLimit.Load()-20 {
		err := c.AddListener(ctx)
		if err != nil {
			return nil, err
		}

		c.shard++
		if c.debugMode.Load() {
			c.logger.Debug(
				"Subscription limited, must add new subscription",
				slog.Any("subscriptions", c.subscriptions[c.shard]),
			)
		}
		return c.subscriptions[c.shard], nil
	}

	return c.subscriptions[c.shard], nil
}

func (c *Client) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, subscription := range c.subscriptions {
		_, cancel := context.WithCancel(subscription.ctx)
		cancel()
	}
}

func (c *Client) CollectMetrics() {
	metrics := c.getMetrics()
	totalShards.Set(float64(metrics.TotalShards))
	aliveShards.Set(float64(metrics.AliveShards))
	deadShards.Set(float64(metrics.DeadShards))
}

func (c *Client) getMetrics() ClientMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	total := len(c.subscriptions)
	alive := 0

	for _, sub := range c.subscriptions {
		if sub.active.Load() {
			alive++
		}
	}

	return ClientMetrics{
		TotalShards: total,
		AliveShards: alive,
		DeadShards:  total - alive,
	}
}
