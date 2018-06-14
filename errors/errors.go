package errors

// Error represents an error.
type Error struct {
	Type    Type
	Message string
	Payload map[string]interface{}
}

// Type represents the type of error.
type Type int

// Possible types of errors.
const (
	NotFound Type = iota
	NotPut
	NotDeleted
	NotMatched
	Unauthorized
	Invalid
	AlreadyExists
	FailedMarshal
	FailedHash
	Unknown
)

func (err Error) Error() string {
	return err.Message
}
