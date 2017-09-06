package errors

import "fmt"

// ErrNotFound error type for objects not found
type ErrNotFound struct {
	// ID unique object identifier.
	ID string
	// Type of the object which wasn't found
	Type string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%v with ID/Name: %v not found", e.Type, e.ID)
}

// ErrValidateVol is error type when a volume fails validation
type ErrValidateVol struct {
	// ID unique object identifier.
	ID string
	// Error is the underlying error
	Cause string
}

func (e *ErrValidateVol) Error() string {
	return fmt.Sprintf("Failed to validate volume: %v Err: %v", e.ID, e.Cause)
}
