package diag

type ExternalServiceError struct {
	NestedError
}

func NewExternalServiceError(m string, e error) *ExternalServiceError {
	err := &ExternalServiceError{}
	err.Message = m
	err.InnerError = e
	return err
}
