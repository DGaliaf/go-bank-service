package custom_error

import "encoding/json"

var (
	ErrNotFound        = NewError(nil, "not found", "")
	ErrNegativeBalance = NewError(nil, "negative balance", "")
	ErrScanFail        = NewError(nil, "failed to scan", "")
	ErrTransaction     = NewError(nil, "transaction failed", "")
	ErrExecuteSQL      = NewError(nil, "failed to execute sql", "")
)

type Error struct {
	Err        error  `json:"-"`
	Msg        string `json:"msg,omitempty"`
	DevMessage string `json:"dev_message,omitempty"`
}

func NewError(err error, msg string, devMessage string) *Error {
	return &Error{Err: err, Msg: msg, DevMessage: devMessage}
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}

	return marshal
}

func (e Error) Unwrap() error {
	return e.Err
}

func systemError(err error) *Error {
	return NewError(err, "internal server error", err.Error())
}
