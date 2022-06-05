package errors

import "fmt"

var (
	ErrNotFound       = fmt.Sprintf("Not Found")
	ErrInvalidRequest = fmt.Sprintf("Unable to Validate Request")
	ErrInternalServer = fmt.Sprintln("Something went wrong")
)
