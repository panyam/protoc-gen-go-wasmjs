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
	"context"
	"reflect"
	"testing"
	"time"
	
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestCallBrowserService_ProtobufPointerInstantiation tests that the CallBrowserService
// function properly instantiates protobuf pointer types before deserialization
func TestCallBrowserService_ProtobufPointerInstantiation(t *testing.T) {
	// Test that pointer types are properly instantiated
	t.Run("PointerTypeInstantiation", func(t *testing.T) {
		// This test validates the reflection-based pointer instantiation logic
		// We can't easily test the full CallBrowserService without a running browser,
		// but we can test the pointer instantiation logic directly
		
		// Test with emptypb.Empty (a common protobuf pointer type)
		type TestResponse = *emptypb.Empty
		
		var resp TestResponse
		respType := reflect.TypeOf(resp)
		
		// This is the same logic from CallBrowserService
		if respType.Kind() == reflect.Ptr {
			respValue := reflect.New(respType.Elem())
			resp = respValue.Interface().(TestResponse)
		}
		
		// Verify that resp is now a non-nil pointer
		if resp == nil {
			t.Error("Expected non-nil pointer after instantiation")
		}
		
		// Verify it's a valid proto.Message
		if _, ok := any(resp).(proto.Message); !ok {
			t.Error("Instantiated pointer should be a proto.Message")
		}
		
		// Verify we can use it as a protobuf message
		if resp.ProtoReflect() == nil {
			t.Error("Expected valid protobuf reflection")
		}
	})
	
	t.Run("ValueTypeHandling", func(t *testing.T) {
		// Test that value types (non-pointers) are handled correctly
		type TestResponse = emptypb.Empty // Value type, not pointer
		
		var resp TestResponse
		respType := reflect.TypeOf(resp)
		
		// Value types shouldn't trigger pointer instantiation
		if respType.Kind() == reflect.Ptr {
			t.Error("Test setup error: expected value type, got pointer")
		}
		
		// Verify that &resp is a valid proto.Message
		if _, ok := any(&resp).(proto.Message); !ok {
			t.Error("Address of value type should be a proto.Message")
		}
	})
}

// TestBrowserServiceChannel_CallTimeout tests that browser service calls
// properly handle timeouts and context cancellation
func TestBrowserServiceChannel_CallTimeout(t *testing.T) {
	// Create a mock browser service channel for testing
	channel := NewBrowserServiceChannel()
	
	// Test context cancellation
	t.Run("ContextCancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel the context immediately
		cancel()
		
		// This should return immediately with context error
		_, err := channel.QueueCall(ctx, "TestService", "testMethod", []byte("{}"), 1*time.Second)
		if err == nil {
			t.Error("Expected context cancellation error")
		}
		
		if err.Error() != "context canceled" {
			t.Errorf("Expected 'context canceled', got: %v", err)
		}
	})
	
	t.Run("TimeoutHandling", func(t *testing.T) {
		ctx := context.Background()
		
		// Use very short timeout to trigger timeout quickly
		start := time.Now()
		_, err := channel.QueueCall(ctx, "TestService", "testMethod", []byte("{}"), 1*time.Millisecond)
		elapsed := time.Since(start)
		
		if err == nil {
			t.Error("Expected timeout error")
		}
		
		// Should timeout quickly
		if elapsed > 100*time.Millisecond {
			t.Errorf("Timeout took too long: %v", elapsed)
		}
		
		if !contains(err.Error(), "timeout") {
			t.Errorf("Expected timeout error, got: %v", err)
		}
	})
}

// TestBrowserServiceChannel_CallIDGeneration tests that unique call IDs are generated
func TestBrowserServiceChannel_CallIDGeneration(t *testing.T) {
	channel := NewBrowserServiceChannel()
	
	// Generate multiple call IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := channel.generateCallID()
		
		if ids[id] {
			t.Errorf("Duplicate call ID generated: %s", id)
		}
		ids[id] = true
		
		if len(id) == 0 {
			t.Error("Empty call ID generated")
		}
	}
	
	// Should have 100 unique IDs
	if len(ids) != 100 {
		t.Errorf("Expected 100 unique IDs, got %d", len(ids))
	}
}

// TestBrowserServiceChannel_PendingCallManagement tests pending call tracking
func TestBrowserServiceChannel_PendingCallManagement(t *testing.T) {
	channel := NewBrowserServiceChannel()
	
	// Test pending call count starts at zero
	if count := channel.GetPendingCallCount(); count != 0 {
		t.Errorf("Expected 0 pending calls initially, got %d", count)
	}
	
	// Test call cleanup
	callID := "test-call-123"
	channel.cleanupCall(callID)
	
	// Should still be zero after cleanup of non-existent call
	if count := channel.GetPendingCallCount(); count != 0 {
		t.Errorf("Expected 0 pending calls after cleanup, got %d", count)
	}
}

// Helper function for string containment check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
		 (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		  someSubstring(s, substr))))
}

func someSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
