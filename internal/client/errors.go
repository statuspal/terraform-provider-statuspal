package statuspal

import (
	"fmt"
	"net/http"
)

// Error represents a custom error that cares the HTTP status, and the body of the failure request.
type Error struct {
	Status int
	Body   string
}

func (e Error) Error() string {
	return fmt.Sprintf("status: %d, body: %s", e.Status, e.Body)
}

func NewError(status int, body []byte) error {
	return &Error{
		Status: status,
		Body:   string(body),
	}
}

// ErrorStatusIs checks if a error received is a HTTP error, and its status.
func ErrorStatusIs(err error, status int) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.Status == status
}

// ErrorNotFound checks if the error returned is a HTTP error, and if the status is [http.StatusNotFound].
func ErrorNotFound(err error) bool {
	return ErrorStatusIs(err, http.StatusNotFound)
}
