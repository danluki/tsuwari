package seventveventapi

type WsMessage struct {
	Operation uint8                  `json:"op"`
	Timestamp int64                  `json:"t,omitempty"`
	D         map[string]interface{} `json:"d"`
}

type ChangeMap struct {
	ID         string        `json:"id"`
	Kind       int8          `json:"kind"`
	Contextual *bool         `json:"contextual,omitempty"`
	Actor      string        `json:"actor"`
	Added      []ChangeField `json:"added,omitempty"`
	Updated    []ChangeField `json:"updated,omitempty"`
	Removed    []ChangeField `json:"removed,omitempty"`
	Pushed     []ChangeField `json:"pushed,omitempty"`
	Pulled     []ChangeField `json:"pulled,omitempty"`
}

type ChangeField struct {
	Key    string `json:"key"`
	Index  *int   `json:"index,omitempty"`
	Nested bool   `json:"nested"`
	/*
		Can be object, nil, or ChangeField
	*/
	OldValue interface{} `json:"old_value"`
	Value    interface{} `json:"value"`
}
