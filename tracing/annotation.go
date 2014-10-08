package tracing

type Annotation struct {
	Name        string      `json:"name,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	From        From        `json:"from,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
	IsBinary    bool        `json:"isBinary, omitempty"`
}
