package errors

// Error represents an error.
type Error struct {
	Type    Type
	Message string
}

// Type represents the type of error.
type Type int

// Possible types of errors.
const (
	NotFound Type = iota
	Unauthorized
	AlreadyExists
	FailedPut
	FailedDelete
	FailedMarshal
)

func (err Error) Error() string {
	return err.Message
}
