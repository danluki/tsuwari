package seventveventapi

type WsMessage struct {
	Operation uint8                  `json:"op"`
	Timestamp int64                  `json:"t,omitempty"`
	D         map[string]interface{} `json:"d"`
}
