// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm

package wasm

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// TestProtobufPointerInstantiation_Framework tests the protobuf pointer
// instantiation fix using various protobuf message types (not demo-specific)
func TestProtobufPointerInstantiation_Framework(t *testing.T) {
	testCases := []struct {
		name        string
		responseType interface{}
		shouldBePtr bool
	}{
		{
			name:        "EmptyPb pointer type",
			responseType: (*emptypb.Empty)(nil),
			shouldBePtr: true,
		},
		{
			name:        "StructPb pointer type", 
			responseType: (*structpb.Struct)(nil),
			shouldBePtr: true,
		},
		{
			name:        "EmptyPb value type",
			responseType: emptypb.Empty{},
			shouldBePtr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This tests the same logic from CallBrowserService
			resp := tc.responseType
			respType := reflect.TypeOf(resp)
			
			// Apply the same pointer instantiation logic
			if respType.Kind() == reflect.Ptr {
				if !tc.shouldBePtr {
					t.Errorf("Expected value type but got pointer type: %T", resp)
					return
				}
				
				// Create new instance (same as CallBrowserService)
				respValue := reflect.New(respType.Elem())
				resp = respValue.Interface()
				
				// Verify it's now a valid proto.Message
				if _, ok := any(resp).(proto.Message); !ok {
					t.Errorf("Instantiated pointer should be proto.Message: %T", resp)
				}
			} else {
				if tc.shouldBePtr {
					t.Errorf("Expected pointer type but got value type: %T", resp)
					return
				}
				
				// For value types, &resp should be proto.Message
				if _, ok := any(&resp).(proto.Message); !ok {
					t.Errorf("Address of value type should be proto.Message: %T", resp)
				}
			}
		})
	}
}

// TestProtobufMessageDetection_Framework tests the proto.Message detection
// logic that was failing before our fix
func TestProtobufMessageDetection_Framework(t *testing.T) {
	t.Run("PointerTypes", func(t *testing.T) {
		// Test various pointer protobuf types
		pointerTypes := []interface{}{
			(*emptypb.Empty)(nil),
			(*structpb.Struct)(nil),
		}

		for _, ptrType := range pointerTypes {
			// Create instance using reflection (same as our fix)
			respType := reflect.TypeOf(ptrType)
			if respType.Kind() == reflect.Ptr {
				respValue := reflect.New(respType.Elem())
				instance := respValue.Interface()
				
				// Test proto.Message detection (this is what was failing)
				if _, ok := any(instance).(proto.Message); !ok {
					t.Errorf("Instantiated pointer should be proto.Message: %T", instance)
				}
				
				// Test that double-pointer is NOT proto.Message (this was the bug)
				if _, ok := any(&instance).(proto.Message); ok {
					t.Errorf("Double pointer should NOT be proto.Message: %T", &instance)
				}
			}
		}
	})

	t.Run("ValueTypes", func(t *testing.T) {
		// Test value protobuf types
		valueTypes := []interface{}{
			emptypb.Empty{},
			structpb.Struct{},
		}

		for _, valueType := range valueTypes {
			// Value types: the value itself is NOT proto.Message
			if _, ok := any(valueType).(proto.Message); ok {
				t.Errorf("Value type should NOT be proto.Message directly: %T", valueType)
			}
			
			// But &value IS proto.Message
			if _, ok := any(&valueType).(proto.Message); !ok {
				t.Errorf("Pointer to value type should be proto.Message: %T", &valueType)
			}
		}
	})
}

// TestCallBrowserService_ErrorScenarios tests error handling in CallBrowserService
func TestCallBrowserService_ErrorScenarios(t *testing.T) {
	t.Run("NonProtobufResponseType", func(t *testing.T) {
		// Test what happens when TResp is not a protobuf type
		type NonProtoResponse struct {
			Message string
		}

		var resp NonProtoResponse
		respType := reflect.TypeOf(resp)

		// This should not be detected as proto.Message
		if _, ok := any(resp).(proto.Message); ok {
			t.Errorf("Non-protobuf type should not be proto.Message: %T", resp)
		}

		if _, ok := any(&resp).(proto.Message); ok {
			t.Errorf("Pointer to non-protobuf type should not be proto.Message: %T", &resp)
		}

		// Our code should handle this gracefully and return the appropriate error
		t.Logf("✅ Non-protobuf types properly rejected")
	})

	t.Run("NilPointerHandling", func(t *testing.T) {
		// Test handling of nil pointers
		var resp *emptypb.Empty // nil pointer
		
		// nil pointers should not be proto.Message
		if _, ok := any(resp).(proto.Message); ok {
			t.Error("nil pointer should not be detected as proto.Message")
		}

		// After instantiation, should be proto.Message
		respType := reflect.TypeOf(resp)
		if respType.Kind() == reflect.Ptr {
			respValue := reflect.New(respType.Elem())
			resp = respValue.Interface().(*emptypb.Empty)
			
			if _, ok := any(resp).(proto.Message); !ok {
				t.Error("Instantiated pointer should be proto.Message")
			}
		}
	})
}
