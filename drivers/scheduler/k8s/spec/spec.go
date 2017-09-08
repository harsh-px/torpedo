package spec

// AppSpec defines a k8s application specification
type AppSpec interface {
	// UID
	ID() string
	// Core
	Core(replicas int32, name string) []interface{}
	// Storage
	Storage() []interface{}
}
