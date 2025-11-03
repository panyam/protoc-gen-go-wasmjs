// Package wasm provides runtime support for protoc-gen-go-wasmjs generated code.
package wasm

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// VTProtoMarshaller implements ProtoMarshaller using vtprotobuf-generated methods.
// This marshaller is designed for TinyGo compatibility and improved performance.
//
// Requirements:
//   - Proto messages must be generated with vtprotobuf plugin
//   - Messages should implement MarshalJSON/UnmarshalJSON methods
//
// To generate vtprotobuf code:
//
//	protoc \
//	  --go_out=. \
//	  --go-vtproto_out=. \
//	  --go-vtproto_opt=features=marshal+unmarshal+size \
//	  your_service.proto
//
// Usage:
//
//	func main() {
//	    wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshaller())
//	    // ... rest of your code
//	}
type VTProtoMarshaller struct {
	// fallbackToStdJSON controls whether to fall back to standard encoding/json
	// if vtprotobuf methods are not available. This is useful during migration
	// or for mixed codebases.
	fallbackToStdJSON bool
}

// NewVTProtoMarshaller creates a new vtprotobuf-based marshaller.
// By default, it will return an error if vtprotobuf methods are not available.
func NewVTProtoMarshaller() *VTProtoMarshaller {
	return &VTProtoMarshaller{
		fallbackToStdJSON: false,
	}
}

// NewVTProtoMarshallerWithFallback creates a vtprotobuf marshaller that falls back
// to standard encoding/json if vtprotobuf methods are not available.
// This is useful for mixed codebases or during migration.
func NewVTProtoMarshallerWithFallback() *VTProtoMarshaller {
	return &VTProtoMarshaller{
		fallbackToStdJSON: true,
	}
}

// Marshal converts a proto message to JSON bytes using vtprotobuf.
//
// The marshaller attempts to use the following methods in order:
//  1. MarshalJSON() ([]byte, error) - vtprotobuf-generated JSON marshaler
//  2. MarshalVT() ([]byte, error) - vtprotobuf binary marshaler (requires JSON conversion)
//  3. Standard proto.Marshal + JSON conversion (if fallback enabled)
//
// Note: MarshalOptions are partially supported. vtprotobuf may not respect all options
// depending on the code generation settings.
func (v *VTProtoMarshaller) Marshal(m proto.Message, opts MarshalOptions) ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("cannot marshal nil message")
	}

	// Try vtprotobuf JSON marshaler first (if available)
	if marshaler, ok := m.(interface{ MarshalJSON() ([]byte, error) }); ok {
		data, err := marshaler.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("vtprotobuf MarshalJSON failed: %w", err)
		}
		return v.applyMarshalOptions(data, opts)
	}

	// Try vtprotobuf binary marshaler (requires conversion to JSON)
	if marshaler, ok := m.(interface{ MarshalVT() ([]byte, error) }); ok {
		// Marshal to protobuf binary
		binaryData, err := marshaler.MarshalVT()
		if err != nil {
			return nil, fmt.Errorf("vtprotobuf MarshalVT failed: %w", err)
		}

		// Create a new message and unmarshal to get JSON
		// This is less efficient but works if JSON marshaler is not generated
		newMsg := proto.Clone(m)
		if err := proto.Unmarshal(binaryData, newMsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal binary: %w", err)
		}

		// Fall through to standard JSON marshaling
		if v.fallbackToStdJSON {
			return json.Marshal(newMsg)
		}
		return nil, fmt.Errorf("message type %T does not implement MarshalJSON and fallback is disabled", m)
	}

	// Fallback to standard JSON marshaling if enabled
	if v.fallbackToStdJSON {
		return json.Marshal(m)
	}

	return nil, fmt.Errorf("message type %T does not implement vtprotobuf MarshalJSON or MarshalVT methods", m)
}

// Unmarshal parses JSON bytes into a proto message using vtprotobuf.
//
// The marshaller attempts to use the following methods in order:
//  1. UnmarshalJSON([]byte) error - vtprotobuf-generated JSON unmarshaler
//  2. UnmarshalVT([]byte) error - vtprotobuf binary unmarshaler (requires JSON conversion)
//  3. Standard json.Unmarshal (if fallback enabled)
//
// Note: UnmarshalOptions are partially supported. vtprotobuf may not respect all options
// depending on the code generation settings.
func (v *VTProtoMarshaller) Unmarshal(data []byte, m proto.Message, opts UnmarshalOptions) error {
	if m == nil {
		return fmt.Errorf("cannot unmarshal into nil message")
	}

	if len(data) == 0 {
		return fmt.Errorf("cannot unmarshal empty data")
	}

	// Try vtprotobuf JSON unmarshaler first (if available)
	if unmarshaler, ok := m.(interface{ UnmarshalJSON([]byte) error }); ok {
		if err := unmarshaler.UnmarshalJSON(data); err != nil {
			return fmt.Errorf("vtprotobuf UnmarshalJSON failed: %w", err)
		}
		return nil
	}

	// Try vtprotobuf binary unmarshaler (requires JSON->binary conversion)
	if _, ok := m.(interface{ UnmarshalVT([]byte) error }); ok {
		// First unmarshal JSON to get the structure
		// Then convert to binary and use UnmarshalVT
		// This is inefficient but works if JSON unmarshaler is not generated

		if v.fallbackToStdJSON {
			// Use standard JSON unmarshal as intermediate step
			if err := json.Unmarshal(data, m); err != nil {
				return fmt.Errorf("json.Unmarshal failed: %w", err)
			}
			return nil
		}
		return fmt.Errorf("message type %T does not implement UnmarshalJSON and fallback is disabled", m)
	}

	// Fallback to standard JSON unmarshaling if enabled
	if v.fallbackToStdJSON {
		if err := json.Unmarshal(data, m); err != nil {
			return fmt.Errorf("json.Unmarshal fallback failed: %w", err)
		}
		return nil
	}

	return fmt.Errorf("message type %T does not implement vtprotobuf UnmarshalJSON or UnmarshalVT methods", m)
}

// applyMarshalOptions applies marshal options to the JSON data.
// Note: This is a best-effort implementation. Some options may not be fully supported
// depending on how the vtprotobuf code was generated.
func (v *VTProtoMarshaller) applyMarshalOptions(data []byte, opts MarshalOptions) ([]byte, error) {
	// If no options are set that we need to handle, return as-is
	if !opts.UseProtoNames && !opts.UseEnumNumbers && opts.EmitUnpopulated {
		return data, nil
	}

	// For more complex option handling, we would need to parse and rewrite the JSON
	// This is left as a TODO for future enhancement
	// For now, we return the data as-is and rely on vtprotobuf generation options

	// TODO: Implement option handling if needed:
	// - UseProtoNames: Rename fields from camelCase to snake_case
	// - UseEnumNumbers: Convert enum strings to numbers
	// - EmitUnpopulated: Add zero-value fields if they're missing

	return data, nil
}

// Ensure VTProtoMarshaller implements ProtoMarshaller
var _ ProtoMarshaller = (*VTProtoMarshaller)(nil)
