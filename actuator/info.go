package actuator

// Info is a model represents for service info.
// Includes Name is the service name
// and a map of Info with key and value
type Info struct {
	Name string                 `json:"service_name"`
	Info map[string]interface{} `json:"info,omitempty"`
}

// Informer is an interface that implemented by
// component to provide their information
type Informer interface {

	// Key returns the key of this information
	Key() string

	// Value returns the value of this information
	Value() interface{}
}
