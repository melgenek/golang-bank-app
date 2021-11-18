package errors

type InternalServerError struct {
	Err error
}

func (err *InternalServerError) Error() string {
	return err.Err.Error()
}
