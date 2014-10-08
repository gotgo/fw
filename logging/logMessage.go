package logging

type LogMessage struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Key     string      `json:"key,omitempty"`
	Value   interface{} `json:"value,omitempty"`
	Kind    Kind        `json:"kind,omitempty"`
}
