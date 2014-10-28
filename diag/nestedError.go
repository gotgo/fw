package diag

import "fmt"

type NestedError struct {
	Message    string
	InnerError error
}

func NewError(m string, e error) *NestedError {
	return &NestedError{
		Message:    m,
		InnerError: e,
	}
}

func (ne *NestedError) Error() string {
	err := ""
	if ne.InnerError != nil {
		err = "inner error : " + ne.InnerError.Error()
	}

	return fmt.Sprintf("%s %s", ne.Message, err)
}

//Private type
type emptyError struct {
}

func (ee *emptyError) Error() string {
	return ""
}
