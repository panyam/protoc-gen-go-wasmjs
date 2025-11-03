// Package wasm provides runtime support for protoc-gen-go-wasmjs generated code.
package wasm

import (
	"google.golang.org/protobuf/proto"
)

// MarshalOptions contains options for marshaling proto messages to JSON.
type MarshalOptions struct {
	// UseProtoNames uses proto field names instead of lowerCamelCase names
	UseProtoNames bool
	// EmitUnpopulated includes zero values in the output
	EmitUnpopulated bool
	// UseEnumNumbers emits enum values as integers instead of strings
	UseEnumNumbers bool
}

// UnmarshalOptions contains options for unmarshaling JSON to proto messages.
type UnmarshalOptions struct {
	// DiscardUnknown ignores unknown fields in the JSON
	DiscardUnknown bool
	// AllowPartial allows partial messages (missing required fields)
	AllowPartial bool
}

// Marshaller is an interface for marshaling proto messages to JSON bytes.
// Implementations can use protojson, vtprotobuf, or any other serialization library.
type Marshaller interface {
	// Marshal converts a proto message to JSON bytes using the provided options
	Marshal(m proto.Message, opts MarshalOptions) ([]byte, error)
}

// Unmarshaller is an interface for unmarshaling JSON bytes to proto messages.
// Implementations can use protojson, vtprotobuf, or any other serialization library.
type Unmarshaller interface {
	// Unmarshal parses JSON bytes into a proto message using the provided options
	Unmarshal(data []byte, m proto.Message, opts UnmarshalOptions) error
}

// ProtoMarshaller combines both Marshaller and Unmarshaller interfaces.
type ProtoMarshaller interface {
	Marshaller
	Unmarshaller
}
