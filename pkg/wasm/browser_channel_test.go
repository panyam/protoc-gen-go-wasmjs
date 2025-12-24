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
	"sync"
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

// TestBrowserServiceChannel_CallIDGeneration tests that unique call IDs are generated
func TestBrowserServiceChannel_CallIDGeneration(t *testing.T) {
	channel := &BrowserServiceChannel{
		pendingCalls: make(map[string]chan *CallResponse),
	}

	// Generate multiple call IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := channel.NextCallID()

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
	channel := &BrowserServiceChannel{
		pendingCalls: make(map[string]chan *CallResponse),
	}

	// Test pending call count starts at zero
	if count := channel.GetPendingCallCount(); count != 0 {
		t.Errorf("Expected 0 pending calls initially, got %d", count)
	}

	// Add a pending call
	callID := "test-call-123"
	responseCh := make(chan *CallResponse, 1)
	channel.mu.Lock()
	channel.pendingCalls[callID] = responseCh
	channel.mu.Unlock()

	// Should have one pending call
	if count := channel.GetPendingCallCount(); count != 1 {
		t.Errorf("Expected 1 pending call, got %d", count)
	}

	// Cleanup the call
	channel.mu.Lock()
	delete(channel.pendingCalls, callID)
	channel.mu.Unlock()

	// Should be back to zero
	if count := channel.GetPendingCallCount(); count != 0 {
		t.Errorf("Expected 0 pending calls after cleanup, got %d", count)
	}
}

// TestFuture_BasicGet tests the basic Get functionality of Future
func TestFuture_BasicGet(t *testing.T) {
	marshaller := &ProtoJSONMarshaller{}
	responseCh := make(chan *CallResponse, 1)

	future := &Future[*emptypb.Empty]{
		responseCh: responseCh,
		marshaller: marshaller,
	}

	// Send a response
	responseCh <- &CallResponse{
		Data: []byte("{}"),
	}

	// Get should return the result
	result, err := future.Get()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Future should be done
	if !future.IsDone() {
		t.Error("Expected future to be done after Get()")
	}
}

// TestFuture_CachedResult tests that subsequent Get calls return cached result
func TestFuture_CachedResult(t *testing.T) {
	marshaller := &ProtoJSONMarshaller{}
	responseCh := make(chan *CallResponse, 1)

	future := &Future[*emptypb.Empty]{
		responseCh: responseCh,
		marshaller: marshaller,
	}

	// Send a response
	responseCh <- &CallResponse{
		Data: []byte("{}"),
	}

	// First Get
	result1, err1 := future.Get()
	if err1 != nil {
		t.Errorf("Expected no error on first Get, got: %v", err1)
	}

	// Second Get should return same cached result
	result2, err2 := future.Get()
	if err2 != nil {
		t.Errorf("Expected no error on second Get, got: %v", err2)
	}

	// Results should be the same instance
	if result1 != result2 {
		t.Error("Expected same cached result on subsequent Get calls")
	}
}

// TestFuture_ErrorCaching tests that errors are also cached
func TestFuture_ErrorCaching(t *testing.T) {
	marshaller := &ProtoJSONMarshaller{}
	responseCh := make(chan *CallResponse, 1)

	future := &Future[*emptypb.Empty]{
		responseCh: responseCh,
		marshaller: marshaller,
	}

	// Send an error response
	responseCh <- &CallResponse{
		Error: proto.Error("test error"),
	}

	// First Get should return error
	_, err1 := future.Get()
	if err1 == nil {
		t.Error("Expected error on first Get")
	}

	// Second Get should return same cached error
	_, err2 := future.Get()
	if err2 == nil {
		t.Error("Expected error on second Get")
	}

	if err1.Error() != err2.Error() {
		t.Error("Expected same cached error on subsequent Get calls")
	}
}

// TestFuture_GetWithTimeout tests the timeout functionality
func TestFuture_GetWithTimeout(t *testing.T) {
	marshaller := &ProtoJSONMarshaller{}
	responseCh := make(chan *CallResponse, 1)

	future := &Future[*emptypb.Empty]{
		responseCh: responseCh,
		marshaller: marshaller,
	}

	// Don't send any response - should timeout
	start := time.Now()
	_, err := future.GetWithTimeout(10 * time.Millisecond)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected timeout error")
	}

	if elapsed < 10*time.Millisecond {
		t.Error("Timeout should have waited at least 10ms")
	}

	if elapsed > 100*time.Millisecond {
		t.Errorf("Timeout took too long: %v", elapsed)
	}

	// Future should still be marked as done after timeout
	if !future.IsDone() {
		t.Error("Expected future to be done after timeout")
	}
}

// TestFuture_ConcurrentGet tests that concurrent Get calls are safe
func TestFuture_ConcurrentGet(t *testing.T) {
	marshaller := &ProtoJSONMarshaller{}
	responseCh := make(chan *CallResponse, 1)

	future := &Future[*emptypb.Empty]{
		responseCh: responseCh,
		marshaller: marshaller,
	}

	// Send a response
	responseCh <- &CallResponse{
		Data: []byte("{}"),
	}

	// Launch multiple concurrent Get calls
	var wg sync.WaitGroup
	results := make([]*emptypb.Empty, 10)
	errors := make([]error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errors[idx] = future.Get()
		}(i)
	}

	wg.Wait()

	// All results should be the same (cached) and no errors
	for i := 0; i < 10; i++ {
		if errors[i] != nil {
			t.Errorf("Goroutine %d got error: %v", i, errors[i])
		}
		if results[i] != results[0] {
			t.Errorf("Goroutine %d got different result", i)
		}
	}
}
