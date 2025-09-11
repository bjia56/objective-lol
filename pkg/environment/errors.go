package environment

type NotFound struct {
	Message string
}

func (e *NotFound) Error() string {
	return e.Message
}

type AlreadyExists struct {
	Message string
}

func (e *AlreadyExists) Error() string {
	return e.Message
}
