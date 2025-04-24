package seventveventapi

import (
	"fmt"
)

type HelloRequest struct {
	/*
		Interval in milliseconds between each heartbeat.
	*/
	HeartbeatInterval int32 `json:"heartbeat_interval"`
	/*
		Unique token for this session, used for resuming and mutating the session
	*/
	SessionID string `json:"session_id"`
	/*
		The maximum amount of subscriptions this connection can intitate
	*/
	SubscriptionLimit int32 `json:"subscription_limit"`
}

func mapToHelloRequest(data map[string]interface{}) (HelloRequest, error) {
	var helloRequest HelloRequest
	if heartbeatInterval, ok := data["heartbeat_interval"].(float64); ok {
		helloRequest.HeartbeatInterval = int32(heartbeatInterval)
	} else {
		return helloRequest, fmt.Errorf("invalid or missing heartbeat_interval")
	}

	if sessionID, ok := data["session_id"].(string); ok {
		helloRequest.SessionID = sessionID
	} else {
		return helloRequest, fmt.Errorf("invalid or missing session_id")
	}

	if subscriptionLimit, ok := data["subscription_limit"].(float64); ok {
		helloRequest.SubscriptionLimit = int32(subscriptionLimit)
	} else {
		return helloRequest, fmt.Errorf("invalid or missing subscription_limit")
	}

	return helloRequest, nil
}

type EndOfStreamRequest struct {
	Code    CloseCode `json:"code"`
	Message string    `json:"message"`
}

func mapToEndOfStreamRequest(data map[string]interface{}) (EndOfStreamRequest, error) {
	var endOfStreamRequest EndOfStreamRequest
	if code, ok := data["code"].(float64); ok {
		endOfStreamRequest.Code = CloseCode(code)
	} else {
		return endOfStreamRequest, fmt.Errorf("invalid or missing code")
	}

	if message, ok := data["message"].(string); ok {
		endOfStreamRequest.Message = message
	} else {
		return endOfStreamRequest, fmt.Errorf("invalid or missing message")
	}

	return endOfStreamRequest, nil
}

type ResumeRequest struct {
	Operation ClientOpcode `json:"op"`
	D         ResumeData   `json:"d"`
}

type ResumeData struct {
	SessionID string `json:"session_id"`
}

type SubscribeRequest struct {
	Operation ClientOpcode  `json:"op"`
	D         SubscribeData `json:"d"`
}

type SubscribeData struct {
	Type      SubscriptionType  `json:"type"`
	Condition map[string]string `json:"condition"`
}

type AckRequest struct {
	Command string                 `json:"command"`
	Data    map[string]interface{} `json:"data"`
}

func mapToAckRequest(data map[string]interface{}) (AckRequest, error) {
	var ackRequest AckRequest
	if command, ok := data["command"].(string); ok {
		ackRequest.Command = command
	} else {
		return ackRequest, fmt.Errorf("invalid or missing command")
	}

	if data, ok := data["data"].(map[string]interface{}); ok {
		ackRequest.Data = data
	} else {
		return ackRequest, fmt.Errorf("invalid or missing data")
	}

	return ackRequest, nil
}
