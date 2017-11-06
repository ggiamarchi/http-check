package common

type Error struct {
	Msg  string
	Code string
}

func (e *Error) Error() string {
	return e.Msg
}
