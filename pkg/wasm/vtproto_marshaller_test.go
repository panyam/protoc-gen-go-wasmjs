// Package wasm provides runtime support for protoc-gen-go-wasmjs generated code.
package wasm

import (
	"encoding/json"
	"errors"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// TestVTProtoMarshallerWithFallback tests the vtproto marshaller with standard proto messages
// using the fallback mechanism
func TestVTProtoMarshallerWithFallback(t *testing.T) {
	marshaller := NewVTProtoMarshallerWithFallback()

	t.Run("Marshal and Unmarshal Empty", func(t *testing.T) {
		msg := &emptypb.Empty{}

		// Marshal
		data, err := marshaller.Marshal(msg, MarshalOptions{
			EmitUnpopulated: true,
		})
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Should produce valid JSON
		if !json.Valid(data) {
			t.Errorf("Marshal produced invalid JSON: %s", string(data))
		}

		// Unmarshal
		result := &emptypb.Empty{}
		err = marshaller.Unmarshal(data, result, UnmarshalOptions{
			DiscardUnknown: true,
		})
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
	})

	t.Run("Marshal and Unmarshal StringValue", func(t *testing.T) {
		msg := &wrapperspb.StringValue{Value: "test string"}

		// Marshal
		data, err := marshaller.Marshal(msg, MarshalOptions{
			EmitUnpopulated: true,
		})
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Should produce valid JSON
		if !json.Valid(data) {
			t.Errorf("Marshal produced invalid JSON: %s", string(data))
		}

		// Unmarshal
		result := &wrapperspb.StringValue{}
		err = marshaller.Unmarshal(data, result, UnmarshalOptions{
			DiscardUnknown: true,
		})
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// Verify value
		if result.Value != msg.Value {
			t.Errorf("Unmarshal value mismatch: got %q, want %q", result.Value, msg.Value)
		}
	})

	t.Run("Marshal and Unmarshal Int32Value", func(t *testing.T) {
		msg := &wrapperspb.Int32Value{Value: 42}

		// Marshal
		data, err := marshaller.Marshal(msg, MarshalOptions{
			EmitUnpopulated: true,
		})
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Should produce valid JSON
		if !json.Valid(data) {
			t.Errorf("Marshal produced invalid JSON: %s", string(data))
		}

		// Unmarshal
		result := &wrapperspb.Int32Value{}
		err = marshaller.Unmarshal(data, result, UnmarshalOptions{
			DiscardUnknown: true,
		})
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// Verify value
		if result.Value != msg.Value {
			t.Errorf("Unmarshal value mismatch: got %d, want %d", result.Value, msg.Value)
		}
	})
}

// TestVTProtoMarshallerWithoutFallback tests error handling when fallback is disabled
func TestVTProtoMarshallerWithoutFallback(t *testing.T) {
	marshaller := NewVTProtoMarshaller()

	t.Run("Marshal without vtproto methods fails", func(t *testing.T) {
		msg := &emptypb.Empty{}

		// Should fail since emptypb.Empty doesn't have vtprotobuf methods
		// and fallback is disabled
		_, err := marshaller.Marshal(msg, MarshalOptions{})
		if err == nil {
			t.Error("Expected Marshal to fail without vtprotobuf methods, but it succeeded")
		}
	})

	t.Run("Unmarshal without vtproto methods fails", func(t *testing.T) {
		msg := &emptypb.Empty{}
		data := []byte("{}")

		// Should fail since emptypb.Empty doesn't have vtprotobuf methods
		// and fallback is disabled
		err := marshaller.Unmarshal(data, msg, UnmarshalOptions{})
		if err == nil {
			t.Error("Expected Unmarshal to fail without vtprotobuf methods, but it succeeded")
		}
	})
}

// TestVTProtoMarshallerErrorCases tests error handling
func TestVTProtoMarshallerErrorCases(t *testing.T) {
	marshaller := NewVTProtoMarshallerWithFallback()

	t.Run("Marshal nil message", func(t *testing.T) {
		_, err := marshaller.Marshal(nil, MarshalOptions{})
		if err == nil {
			t.Error("Expected error when marshaling nil message")
		}
	})

	t.Run("Unmarshal nil message", func(t *testing.T) {
		err := marshaller.Unmarshal([]byte("{}"), nil, UnmarshalOptions{})
		if err == nil {
			t.Error("Expected error when unmarshaling into nil message")
		}
	})

	t.Run("Unmarshal empty data", func(t *testing.T) {
		msg := &emptypb.Empty{}
		err := marshaller.Unmarshal([]byte{}, msg, UnmarshalOptions{})
		if err == nil {
			t.Error("Expected error when unmarshaling empty data")
		}
	})

	t.Run("Unmarshal invalid JSON", func(t *testing.T) {
		msg := &emptypb.Empty{}
		err := marshaller.Unmarshal([]byte("{invalid json}"), msg, UnmarshalOptions{})
		if err == nil {
			t.Error("Expected error when unmarshaling invalid JSON")
		}
	})
}

// TestVTProtoMarshallerImplementsInterface verifies the marshaller implements ProtoMarshaller
func TestVTProtoMarshallerImplementsInterface(t *testing.T) {
	var _ ProtoMarshaller = (*VTProtoMarshaller)(nil)
	var _ ProtoMarshaller = NewVTProtoMarshaller()
	var _ ProtoMarshaller = NewVTProtoMarshallerWithFallback()
}

// TestVTProtoMarshallerJSONValidity ensures marshaled output is valid JSON
func TestVTProtoMarshallerJSONValidity(t *testing.T) {
	marshaller := NewVTProtoMarshallerWithFallback()

	testCases := []struct {
		name string
		msg  proto.Message
	}{
		{"Empty", &emptypb.Empty{}},
		{"StringValue", &wrapperspb.StringValue{Value: "test"}},
		{"Int32Value", &wrapperspb.Int32Value{Value: 123}},
		{"BoolValue", &wrapperspb.BoolValue{Value: true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := marshaller.Marshal(tc.msg, MarshalOptions{
				EmitUnpopulated: true,
			})
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if !json.Valid(data) {
				t.Errorf("Marshal produced invalid JSON: %s", string(data))
			}

			// Verify we can parse it back
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Errorf("Failed to parse marshaled JSON: %v", err)
			}
		})
	}
}

// MockVTProtoMessage is a mock message that implements vtprotobuf methods
type MockVTProtoMessage struct {
	Value       string
	marshalErr  error
	unmarshalErr error
}

func (m *MockVTProtoMessage) Reset()         {}
func (m *MockVTProtoMessage) String() string { return m.Value }
func (m *MockVTProtoMessage) ProtoMessage()  {}
func (m *MockVTProtoMessage) ProtoReflect() protoreflect.Message {
	return nil // Minimal implementation for testing
}

func (m *MockVTProtoMessage) MarshalJSON() ([]byte, error) {
	if m.marshalErr != nil {
		return nil, m.marshalErr
	}
	return json.Marshal(map[string]string{"value": m.Value})
}

func (m *MockVTProtoMessage) UnmarshalJSON(data []byte) error {
	if m.unmarshalErr != nil {
		return m.unmarshalErr
	}
	var obj map[string]string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	m.Value = obj["value"]
	return nil
}

// TestVTProtoMarshallerWithMockVTProto tests the marshaller with a mock vtprotobuf message
func TestVTProtoMarshallerWithMockVTProto(t *testing.T) {
	marshaller := NewVTProtoMarshaller() // No fallback needed

	t.Run("Marshal and Unmarshal with vtproto methods", func(t *testing.T) {
		msg := &MockVTProtoMessage{Value: "test value"}

		// Marshal
		data, err := marshaller.Marshal(msg, MarshalOptions{})
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Unmarshal
		result := &MockVTProtoMessage{}
		err = marshaller.Unmarshal(data, result, UnmarshalOptions{})
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// Verify
		if result.Value != msg.Value {
			t.Errorf("Value mismatch: got %q, want %q", result.Value, msg.Value)
		}
	})

	t.Run("Marshal error propagation", func(t *testing.T) {
		msg := &MockVTProtoMessage{
			Value:      "test",
			marshalErr: errors.New("mock marshal error"),
		}

		_, err := marshaller.Marshal(msg, MarshalOptions{})
		if err == nil {
			t.Error("Expected marshal error to be propagated")
		}
	})

	t.Run("Unmarshal error propagation", func(t *testing.T) {
		msg := &MockVTProtoMessage{
			unmarshalErr: errors.New("mock unmarshal error"),
		}

		err := marshaller.Unmarshal([]byte(`{"value":"test"}`), msg, UnmarshalOptions{})
		if err == nil {
			t.Error("Expected unmarshal error to be propagated")
		}
	})
}
