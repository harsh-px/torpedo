package spec

// AppSpec defines a k8s application specification
type AppSpec interface {
	// UID
	ID() string
	// Core
	Core(name string) []interface{}
	// Storage
	Storage() []interface{}
}
