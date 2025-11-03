// Package wasm provides runtime support for protoc-gen-go-wasmjs generated code.
package wasm

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtojsonMarshaller implements ProtoMarshaller using protojson encoding.
// This is the default marshaller and is compatible with most Go environments.
// For TinyGo compatibility, consider using vtprotobuf or a custom marshaller.
type ProtojsonMarshaller struct{}

// NewProtojsonMarshaller creates a new protojson-based marshaller.
func NewProtojsonMarshaller() *ProtojsonMarshaller {
	return &ProtojsonMarshaller{}
}

// Marshal converts a proto message to JSON bytes using protojson.
func (p *ProtojsonMarshaller) Marshal(m proto.Message, opts MarshalOptions) ([]byte, error) {
	marshalOpts := protojson.MarshalOptions{
		UseProtoNames:   opts.UseProtoNames,
		EmitUnpopulated: opts.EmitUnpopulated,
		UseEnumNumbers:  opts.UseEnumNumbers,
	}
	return marshalOpts.Marshal(m)
}

// Unmarshal parses JSON bytes into a proto message using protojson.
func (p *ProtojsonMarshaller) Unmarshal(data []byte, m proto.Message, opts UnmarshalOptions) error {
	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: opts.DiscardUnknown,
		AllowPartial:   opts.AllowPartial,
	}
	return unmarshalOpts.Unmarshal(data, m)
}

// Ensure ProtojsonMarshaller implements ProtoMarshaller
var _ ProtoMarshaller = (*ProtojsonMarshaller)(nil)
