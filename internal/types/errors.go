package types

// ErrorCode represents the type of an error.
type ErrorCode int

const (
	ErrBadRequest ErrorCode = iota + 1
	ErrUnauthorized
	ErrForbidden
	ErrNotFound
	ErrInternal
)

// Error represents an error object that can be returned by the API.
type Error struct {
	Code   ErrorCode `json:"code"`
	Reason string    `json:"message"`
}

// NewError creates a new Error object with the given reason with default code ErrInternal.
// ErrInternal coded errors will not expose the reason to the client.
// Use WithCode to set a different error code and expose the reason to the client.
func NewError(reason string) *Error {
	return &Error{
		Code:   ErrInternal,
		Reason: reason,
	}
}

// NewFromError creates a new Error object from the given error.
func NewFromError(err error) *Error {
	if err == nil {
		return nil
	}
	return NewError(err.Error())
}

// WithCode sets the error code of the Error object.
func (e *Error) WithCode(code ErrorCode) *Error {
	e.Code = code
	return e
}

// Error wraps the reason of the Error object.
// If the error code is ErrInternal, the error message will be "internal error".
func (e Error) Error() string {
	switch e.Code {
	case ErrInternal:
		return "internal error"
	default:
		return e.Reason
	}
}
