package errs

type ErrNotFound struct {
	Msg string
}

func (err *ErrNotFound) Error() string {
	return err.Msg
}
