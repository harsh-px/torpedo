package k8sutils

import (
	"errors"
	"fmt"
)

// ErrK8SApiAccountNotSet is returned when the account used to talk to k8s api is not setup
var ErrK8SApiAccountNotSet = errors.New("k8s api account is not setup")

// ErrFailedToParseYAML error type for objects not found
type ErrFailedToParseYAML struct {
	// Path is the path of the yaml file that was to be parsed
	Path string
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToParseYAML) Error() string {
	return fmt.Sprintf("Failed to parse file: %v due to err: %v", e.Path, e.Cause)
}

// ErrFailedToApplySpec error type for failing to apply a spec file
type ErrFailedToApplySpec struct {
	// Path is the path of the yaml file that was to be applied
	Path string
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToApplySpec) Error() string {
	return fmt.Sprintf("Failed to apply spec file: %v due to err: %v", e.Path, e.Cause)
}