// Package codec is an interface for encoding messages
// From https://github.com/micro/go-micro/blob/master/codec/codec.go
package codec

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
