package k8s

import (
	"fmt"

	"github.com/portworx/torpedo/drivers/scheduler"
)

// ErrFailedToScheduleApp error type for failing to schedule an app
type ErrFailedToScheduleApp struct {
	// App is the app that failed to schedule
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToScheduleApp) Error() string {
	return fmt.Sprintf("Failed to schedule app: %v due to err: %v", e.App.Name, e.Cause)
}

// ErrFailedToDestroyApp error type for failing to destroy an app
type ErrFailedToDestroyApp struct {
	// App is the app that failed to destroy
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToDestroyApp) Error() string {
	return fmt.Sprintf("Failed to destory app: %v due to err: %v", e.App.Name, e.Cause)
}

// ErrFailedToDestroyStorage error type for failing to destroy an app's storage
type ErrFailedToDestroyStorage struct {
	// App is the app that failed to destroy
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToDestroyStorage) Error() string {
	return fmt.Sprintf("Failed to destory storage for app: %v due to err: %v", e.App.Name, e.Cause)
}

// ErrFailedToValidateStorage error type for failing to validate an app's storage
type ErrFailedToValidateStorage struct {
	// App is the app that failed to destroy
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToValidateStorage) Error() string {
	return fmt.Sprintf("Failed to validate storage for app: %v due to err: %v", e.App.Name, e.Cause)
}

// ErrFailedToValidateApp error type for failing to validate an app
type ErrFailedToValidateApp struct {
	// App is the app that failed to destroy
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToValidateApp) Error() string {
	return fmt.Sprintf("Failed to validate app: %v due to err: %v", e.App.Name, e.Cause)
}

// ErrFailedToGetVolumesForApp error type for failing to get an app's volumes
type ErrFailedToGetVolumesForApp struct {
	// App is the app that failed to destroy
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToGetVolumesForApp) Error() string {
	return fmt.Sprintf("Failed to get volumes for app: %v due to err: %v", e.App.Name, e.Cause)
}

// ErrFailedToGetVolumesParameters error type for failing to get an app's volume paramters
type ErrFailedToGetVolumesParameters struct {
	// App is the app for which we failed to get volume parameters
	App scheduler.App
	// Cause is the underlying cause of the error
	Cause string
}

func (e *ErrFailedToGetVolumesParameters) Error() string {
	return fmt.Sprintf("Failed to get volume parameters for app: %v due to err: %v", e.App.Name, e.Cause)
}