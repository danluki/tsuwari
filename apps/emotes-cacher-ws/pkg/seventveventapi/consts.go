package seventveventapi

type ServerOpcode int

const (
	/*
		Sent by server. A standard event message, sent when a subscribed event is emitted
	*/
	Dispatch ServerOpcode = 0
	/*
		Sent by server. Received upon connecting, presents info about the session
	*/
	Hello ServerOpcode = 1
	/*
		Sent by server. Ensures the connection is still alive
	*/
	Heartbeat ServerOpcode = 2
	/*
		Sent by server. Server wants the client to reconnect
	*/
	Reconnect ServerOpcode = 4
	/*
		Sent by server. Server acknowledges an action by the client
	*/
	Ack ServerOpcode = 5
	/*
		Sent by server. An error occured, you should log this
	*/
	Error ServerOpcode = 6
	/*
		Sent by server. The server will send no further data and imminently end the connection
	*/
	EndOfStream ServerOpcode = 7
)

type ClientOpcode int

const (
	/*
		Sent by client. Authenticate with an account
	*/
	Identify ClientOpcode = 33
	/*
		Sent by client. Try to resume a previous session
	*/
	Resume ClientOpcode = 34
	/*
		Sent by client. Watch for changes on specific objects or sources. Don't smash it!
	*/
	Subscribe ClientOpcode = 35
	/*
		Sent by client. Stop listening for changes
	*/
	Unsubscribe ClientOpcode = 36
	/*
		Sent by client.
	*/
	Signal ClientOpcode = 37
)

type CloseCode int

const (
	/*
		An error occured on the server's end.
		Can reconnect.
	*/
	ServerError CloseCode = 4000
	/*
		The client sent an unexpected ServerOpcode. Can't reconnect.
		This code indicate a bad client implementation. you must log such error and fix the issue before reconnecting.
	*/
	UnknownOperation CloseCode = 4001
	/*
		The client sent a payload that couldn't be decoded. Can't reconnect.
		This code indicate a bad client implementation. you must log such error and fix the issue before reconnecting.
	*/
	InvalidPayload CloseCode = 4002
	/*
		The client unsucessfully tried to identify. Can't reconnect.
		This code indicate a bad client implementation. you must log such error and fix the issue before reconnecting.
	*/
	AuthFailure CloseCode = 4003
	/*
		the client wanted to identify again. Can't reconnect.
		This code indicate a bad client implementation. you must log such error and fix the issue before reconnecting.
	*/
	AlreadyIdentified CloseCode = 4004
	/*
		The client is being rate-limited. Only reconnect if this was initiated by action of the end-user.
	*/
	RateLimited CloseCode = 4005
	/*
		The server is restarting and the client should reconnect. Can reconnect.
	*/
	Restart CloseCode = 4006
	/*
		The server is in maintenance mode and not accepting connections.
		Can reconnect. Reconnect with significantly greater delay, i.e at least 5 minutes, including jitter.
	*/
	Maintenance CloseCode = 4007
	/*
		The client was idle for too long. Can reconnect.
	*/
	Timeout CloseCode = 4008
	/*
		The client tried to subscribe to an event twice. Can't reconnect.
		This code indicate a bad client implementation. you must log such error and fix the issue before reconnecting.
	*/
	AlreadySubscribed CloseCode = 4009
	/*
		The client tried to unsubscribe from an event they weren't subscribing to. Can't reconnect.
		This code indicate a bad client implementation. you must log such error and fix the issue before reconnecting.
	*/
	NotSubscribed CloseCode = 4010
	/*
		The client did something that they did not have permission for. Only reconnect if this was initiated by action of the end-user.
	*/
	InsufficientPrivilege CloseCode = 4011
)

type SubscriptionType string

const (
	SystemAnnouncement   SubscriptionType = "system.announcement"
	CreateEmote          SubscriptionType = "emote.create"
	UpdateEmote          SubscriptionType = "emote.update"
	DeleteEmote          SubscriptionType = "emote.delete"
	PatchEmote           SubscriptionType = "emote.*"
	CreateEmoteSet       SubscriptionType = "emote_set.create"
	UpdateEmoteSet       SubscriptionType = "emote_set.update"
	DeleteEmoteSet       SubscriptionType = "emote_set.delete"
	PatchEmoteSet        SubscriptionType = "emote_set.*"
	CreateUser           SubscriptionType = "user.create"
	DeleteUser           SubscriptionType = "user.delete"
	AddUserConnection    SubscriptionType = "user.add_connection"
	UpdateUserConnection SubscriptionType = "user.update_connection"
	DeleteUserConnection SubscriptionType = "user.delete_connection"
	CreateCosmetic       SubscriptionType = "cosmetic.create"
	UpdateCosmetic       SubscriptionType = "cosmetic.update"
	DeleteCosmetic       SubscriptionType = "cosmetic.delete"
	PatchCosmetic        SubscriptionType = "cosmetic.*"
	CreateEntitlement    SubscriptionType = "entitlement.create"
	UpdateEntitlement    SubscriptionType = "entitlement.update"
	DeleteEntitlement    SubscriptionType = "entitlement.delete"
	PatchEntitlement     SubscriptionType = "entitlement.delete"
)
