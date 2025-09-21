// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall/js"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// BrowserCall represents a call from WASM to a browser-provided service
type BrowserCall struct {
	ID         string
	Service    string // Service name (e.g., "BrowserAPI")
	Method     string // Method name (e.g., "Fetch")
	Request    []byte // Serialized proto request
	ResponseCh chan *CallResponse
	Timeout    time.Duration
	StartTime  time.Time
	IsAsync    bool // Whether this is an async browser method
}

// CallResponse represents the response from a browser service call
type CallResponse struct {
	Data  []byte
	Error error
}

// BrowserServiceChannel manages all browser service calls with FIFO ordering
type BrowserServiceChannel struct {
	callQueue    chan *BrowserCall
	pendingCalls map[string]*PendingCall
	mu           sync.RWMutex
	nextCallID   uint64
	initialized  bool
}

// PendingCall tracks an in-flight browser service call
type PendingCall struct {
	Call     *BrowserCall
	Timer    *time.Timer
	RefCount int32
}

// Global singleton browser channel instance
var (
	browserChannel     *BrowserServiceChannel
	browserChannelOnce sync.Once
)

// GetBrowserChannel returns the singleton browser channel instance
func GetBrowserChannel() *BrowserServiceChannel {
	browserChannelOnce.Do(func() {
		browserChannel = &BrowserServiceChannel{
			callQueue:    make(chan *BrowserCall, 100),
			pendingCalls: make(map[string]*PendingCall),
		}
		browserChannel.Initialize()
	})
	return browserChannel
}

// Initialize sets up the browser channel and registers JS callbacks
func (bc *BrowserServiceChannel) Initialize() {
	if bc.initialized {
		return
	}
	bc.initialized = true

	// Register JS function to get next browser call
	js.Global().Set("__wasmGetNextBrowserCall", js.FuncOf(func(this js.Value, args []js.Value) any {
		select {
		case call := <-bc.callQueue:
			bc.registerPendingCall(call)

			// Return call details to JavaScript
			return map[string]any{
				"id":      call.ID,
				"service": call.Service,
				"method":  call.Method,
				"request": string(call.Request),
			}
		case <-time.After(10 * time.Millisecond):
			// Non-blocking check, return null if no calls pending
			return js.Null()
		}
	}))

	// Register JS function to deliver browser call response
	js.Global().Set("__wasmDeliverBrowserResponse", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 3 {
			return false
		}

		callID := args[0].String()
		responseData := args[1]
		errorMsg := args[2]

		bc.mu.RLock()
		pending, exists := bc.pendingCalls[callID]
		bc.mu.RUnlock()

		if !exists {
			return false
		}

		// Prepare response
		var response CallResponse
		if !errorMsg.IsNull() && !errorMsg.IsUndefined() {
			response.Error = errors.New(errorMsg.String())
		} else if !responseData.IsNull() && !responseData.IsUndefined() {
			response.Data = []byte(responseData.String())
		}

		// Send response to waiting goroutine
		select {
		case pending.Call.ResponseCh <- &response:
		default:
			// Channel might be closed if timeout occurred
		}

		// Cleanup
		bc.cleanupCall(callID)
		return true
	}))

	// Start background processor for timeouts
	go bc.processTimeouts()
}

// NextCallID generates a unique call ID
func (bc *BrowserServiceChannel) NextCallID() string {
	id := atomic.AddUint64(&bc.nextCallID, 1)
	return fmt.Sprintf("call_%d_%d", time.Now().UnixNano(), id)
}

// QueueCall queues a synchronous browser service call and waits for response
func (bc *BrowserServiceChannel) QueueCall(ctx context.Context, service, method string, request []byte, timeout time.Duration) ([]byte, error) {
	return bc.queueCallInternal(ctx, service, method, request, timeout, false)
}

// QueueCallAsync queues an async browser service call and waits for response
func (bc *BrowserServiceChannel) QueueCallAsync(ctx context.Context, service, method string, request []byte, timeout time.Duration) ([]byte, error) {
	return bc.queueCallInternal(ctx, service, method, request, timeout, true)
}

// queueCallInternal is the internal implementation for queuing calls
func (bc *BrowserServiceChannel) queueCallInternal(ctx context.Context, service, method string, request []byte, timeout time.Duration, isAsync bool) ([]byte, error) {
	callID := bc.NextCallID()
	responseCh := make(chan *CallResponse, 1)

	call := &BrowserCall{
		ID:         callID,
		Service:    service,
		Method:     method,
		Request:    request,
		ResponseCh: responseCh,
		Timeout:    timeout,
		StartTime:  time.Now(),
		IsAsync:    isAsync,
	}

	// Queue the call
	select {
	case bc.callQueue <- call:
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout queuing browser call")
	}

	// Wait for response
	select {
	case resp := <-responseCh:
		if resp.Error != nil {
			return nil, resp.Error
		}
		return resp.Data, nil
	case <-ctx.Done():
		bc.cleanupCall(callID)
		return nil, ctx.Err()
	case <-time.After(timeout):
		bc.cleanupCall(callID)
		return nil, fmt.Errorf("browser call timeout after %v", timeout)
	}
}

// registerPendingCall registers a call as pending with timeout
func (bc *BrowserServiceChannel) registerPendingCall(call *BrowserCall) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	timer := time.AfterFunc(call.Timeout, func() {
		bc.handleTimeout(call.ID)
	})

	bc.pendingCalls[call.ID] = &PendingCall{
		Call:     call,
		Timer:    timer,
		RefCount: 1,
	}
}

// handleTimeout handles a call timeout
func (bc *BrowserServiceChannel) handleTimeout(callID string) {
	bc.mu.Lock()
	pending, exists := bc.pendingCalls[callID]
	if !exists {
		bc.mu.Unlock()
		return
	}
	delete(bc.pendingCalls, callID)
	bc.mu.Unlock()

	// Send timeout error
	select {
	case pending.Call.ResponseCh <- &CallResponse{
		Error: fmt.Errorf("browser call timeout"),
	}:
	default:
	}

	// Cleanup
	close(pending.Call.ResponseCh)
}

// cleanupCall cleans up a completed or cancelled call
func (bc *BrowserServiceChannel) cleanupCall(callID string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	pending, exists := bc.pendingCalls[callID]
	if !exists {
		return
	}

	// Stop timeout timer
	if pending.Timer != nil {
		pending.Timer.Stop()
	}

	// Decrement reference count
	if atomic.AddInt32(&pending.RefCount, -1) <= 0 {
		delete(bc.pendingCalls, callID)
		close(pending.Call.ResponseCh)
	}
}

// processTimeouts periodically checks for timed-out calls
func (bc *BrowserServiceChannel) processTimeouts() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		bc.mu.Lock()

		var timedOut []string
		for id, pending := range bc.pendingCalls {
			if now.Sub(pending.Call.StartTime) > pending.Call.Timeout {
				timedOut = append(timedOut, id)
			}
		}
		bc.mu.Unlock()

		// Handle timeouts outside of lock
		for _, id := range timedOut {
			bc.handleTimeout(id)
		}
	}
}

// GetPendingCallCount returns the number of pending calls (for debugging)
func (bc *BrowserServiceChannel) GetPendingCallCount() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.pendingCalls)
}

// CallBrowserService is a generic helper for calling synchronous browser services
// The browser method should return a value directly (not a Promise)
func CallBrowserService[TReq any, TResp any](channel *BrowserServiceChannel, ctx context.Context, serviceName, methodName string, req TReq) (TResp, error) {
var resp TResp

// If TResp is a pointer type, we need to create a new instance
// This is necessary for protobuf message types which are pointers
respType := reflect.TypeOf(resp)
fmt.Printf("DEBUG: CallBrowserService - Initial resp type=%T, kind=%v\n", resp, respType.Kind())
if respType.Kind() == reflect.Ptr {
// Create a new instance of the underlying type
respValue := reflect.New(respType.Elem())
 resp = respValue.Interface().(TResp)
		fmt.Printf("DEBUG: CallBrowserService - Created new instance, resp type=%T\n", resp)
	}

	// Marshal the request using protojson
	opts := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
		UseEnumNumbers:  false,
	}

	// Use reflection to get the proto message interface
	reqMsg, ok := any(req).(proto.Message)
	if !ok {
		return resp, fmt.Errorf("request is not a proto message")
	}

	requestData, err := opts.Marshal(reqMsg)
	if err != nil {
		return resp, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call browser service through the channel
	fmt.Printf("DEBUG: About to queue browser call: %s.%s\n", serviceName, methodName)
	responseData, err := channel.QueueCall(ctx, serviceName, methodName, requestData, 30*time.Second)
	if err != nil {
		fmt.Printf("DEBUG: QueueCall failed: %v\n", err)
		return resp, err
	}
	fmt.Printf("DEBUG: QueueCall succeeded, got response data (len=%d): %s\n", len(responseData), string(responseData))

	// Unmarshal the response
	// Check if resp is already a proto.Message (if it's a pointer type)
	// or if we need to take its address (if it's a value type)
	var respMsg proto.Message
	var isProtoMsg bool
	
	// Try resp directly first (for pointer types like *PromptResponse)
	if respMsg, isProtoMsg = any(resp).(proto.Message); !isProtoMsg {
		// Try &resp (for value types)
		respMsg, isProtoMsg = any(&resp).(proto.Message)
	}
	
	if !isProtoMsg {
		fmt.Printf("DEBUG: respType=%T, resp=%+v\n", resp, resp)
		fmt.Printf("DEBUG: responseData=%s\n", string(responseData))
		return resp, fmt.Errorf("response is not a proto message (type: %T)", resp)
	}

	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := unmarshalOpts.Unmarshal(responseData, respMsg); err != nil {
		return resp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

// CallBrowserServiceAsync is a generic helper for calling async browser services
// The browser method returns a Promise and we need to handle it with a callback
// This is necessary for browser APIs that are inherently async (fetch, IndexedDB, etc.)
func CallBrowserServiceAsync[TReq any, TResp any](channel *BrowserServiceChannel, ctx context.Context, serviceName, methodName string, req TReq) (TResp, error) {
	var resp TResp

	// If TResp is a pointer type, we need to create a new instance
	// This is necessary for protobuf message types which are pointers
	respType := reflect.TypeOf(resp)
	if respType.Kind() == reflect.Ptr {
		// Create a new instance of the underlying type
		respValue := reflect.New(respType.Elem())
		resp = respValue.Interface().(TResp)
	}

	// For async methods, we need to tell the browser side to handle it as a Promise
	// We'll add a special flag in the call to indicate async handling

	// Marshal the request using protojson
	opts := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
		UseEnumNumbers:  false,
	}

	reqMsg, ok := any(req).(proto.Message)
	if !ok {
		return resp, fmt.Errorf("request is not a proto message")
	}

	requestData, err := opts.Marshal(reqMsg)
	if err != nil {
		return resp, fmt.Errorf("failed to marshal request: %w", err)
	}

	// For async calls, we use a longer timeout since they may involve network operations
	responseData, err := channel.QueueCallAsync(ctx, serviceName, methodName, requestData, 60*time.Second)
	if err != nil {
		return resp, err
	}

	// Unmarshal the response
	respMsg, ok := any(&resp).(proto.Message)
	if !ok {
		return resp, fmt.Errorf("response is not a proto message")
	}

	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := unmarshalOpts.Unmarshal(responseData, respMsg); err != nil {
		return resp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}
