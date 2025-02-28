package app

type Error struct {
	Msg string
	Err error
}

func (e *Error) Ok() bool {
	return e.Err == nil
}
