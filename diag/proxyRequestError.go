package diag

import "fmt"

type ProxyRequestError struct {
	RequestFor string
	Provider   string
	Message    string
	InnerError error
}

func NewProxyRequestError(f, p, m string, err error) *ProxyRequestError {
	pre := &ProxyRequestError{
		Provider:   p,
		RequestFor: f,
		Message:    m,
		InnerError: err,
	}
	return pre
}

func (pre *ProxyRequestError) Error() string {
	err := ""
	if pre.InnerError != nil {
		err = pre.InnerError.Error()
	}

	return fmt.Sprintf("%s Request for %s using %s failed due to %s", pre.Message, pre.RequestFor, pre.Provider, err)
}
