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

	"google.golang.org/protobuf/proto"
)

// CallResponse represents the response from a browser service call.
type CallResponse struct {
	// Data contains the serialized protobuf response data (JSON-encoded).
	Data []byte

	// Error contains any error that occurred during the call.
	Error error
}

// Future represents an async operation that will complete in the future.
// The result can be retrieved using Get(), which blocks on first call
// and returns cached result on subsequent calls.
type Future[TResp proto.Message] struct {
	responseCh chan *CallResponse
	marshaller *ProtoJSONMarshaller

	// Cached result after first Get()
	result TResp
	err    error
	done   bool
	mu     sync.Mutex
}

// Get blocks until the async operation completes and returns the result.
// On first call, it waits for the response, unmarshals it, and caches the result.
// On subsequent calls, it returns the cached result immediately.
// The response channel is closed after the first Get() completes.
func (f *Future[TResp]) Get() (TResp, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Already resolved - return cached result
	if f.done {
		return f.result, f.err
	}

	// First call - block and cache
	resp := <-f.responseCh
	f.done = true

	// Close channel since we're done with it
	close(f.responseCh)

	if resp.Error != nil {
		f.err = resp.Error
		return f.result, f.err
	}

	// Create new instance if TResp is a pointer type
	respType := reflect.TypeOf(f.result)
	if respType.Kind() == reflect.Ptr {
		respValue := reflect.New(respType.Elem())
		f.result = respValue.Interface().(TResp)
	}

	// Unmarshal and cache
	f.err = f.marshaller.Unmarshal(resp.Data, f.result, UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	})

	return f.result, f.err
}

// GetWithTimeout is like Get but with a timeout.
func (f *Future[TResp]) GetWithTimeout(timeout time.Duration) (TResp, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Already resolved - return cached result
	if f.done {
		return f.result, f.err
	}

	// First call - block with timeout
	select {
	case resp := <-f.responseCh:
		f.done = true
		close(f.responseCh)

		if resp.Error != nil {
			f.err = resp.Error
			return f.result, f.err
		}

		// Create new instance if TResp is a pointer type
		respType := reflect.TypeOf(f.result)
		if respType.Kind() == reflect.Ptr {
			respValue := reflect.New(respType.Elem())
			f.result = respValue.Interface().(TResp)
		}

		// Unmarshal and cache
		f.err = f.marshaller.Unmarshal(resp.Data, f.result, UnmarshalOptions{
			DiscardUnknown: true,
			AllowPartial:   true,
		})

		return f.result, f.err

	case <-time.After(timeout):
		f.done = true
		f.err = fmt.Errorf("timeout waiting for async response after %v", timeout)
		close(f.responseCh)
		return f.result, f.err
	}
}

// IsDone returns true if the Future has been resolved (either success or error).
func (f *Future[TResp]) IsDone() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.done
}

// BrowserServiceChannel manages browser service calls from WASM.
// It provides both synchronous and asynchronous calling patterns.
//
// Synchronous calls directly invoke JavaScript functions and return immediately.
// Asynchronous calls return a Future that can be used to retrieve the result later.
type BrowserServiceChannel struct {
	// pendingCalls tracks async calls awaiting response.
	pendingCalls map[string]chan *CallResponse

	// mu protects concurrent access to pendingCalls map.
	mu sync.RWMutex

	// nextCallID is atomically incremented to generate unique call IDs.
	nextCallID uint64

	// initialized indicates if Initialize() has been called.
	initialized bool
}

// Global singleton browser channel instance.
var (
	browserChannel     *BrowserServiceChannel
	browserChannelOnce sync.Once
)

// GetBrowserChannel returns the singleton BrowserServiceChannel instance.
func GetBrowserChannel() *BrowserServiceChannel {
	browserChannelOnce.Do(func() {
		browserChannel = &BrowserServiceChannel{
			pendingCalls: make(map[string]chan *CallResponse),
		}
		browserChannel.Initialize()
	})
	return browserChannel
}

// Initialize sets up the browser channel and registers JS callbacks.
func (bc *BrowserServiceChannel) Initialize() {
	if bc.initialized {
		return
	}
	bc.initialized = true

	// Register JS function to deliver async browser call responses
	js.Global().Set("__wasmDeliverBrowserResponse", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 3 {
			return false
		}

		callID := args[0].String()
		responseData := args[1]
		errorMsg := args[2]

		bc.mu.RLock()
		responseCh, exists := bc.pendingCalls[callID]
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

		// Send response to waiting Future
		select {
		case responseCh <- &response:
		default:
			// Channel might be closed if timeout occurred
		}

		// Cleanup
		bc.mu.Lock()
		delete(bc.pendingCalls, callID)
		bc.mu.Unlock()

		return true
	}))
}

// NextCallID generates a unique call ID for async operations.
func (bc *BrowserServiceChannel) NextCallID() string {
	id := atomic.AddUint64(&bc.nextCallID, 1)
	return fmt.Sprintf("call_%d_%d", time.Now().UnixNano(), id)
}

// CallSync makes a synchronous call to a browser service.
// The JavaScript function must return a value directly (not a Promise).
// This call blocks until JavaScript returns.
func (bc *BrowserServiceChannel) CallSync(service, method string, request []byte) ([]byte, error) {
	// Get the browser services object
	browserServices := js.Global().Get("__browserServices")
	if browserServices.IsUndefined() || browserServices.IsNull() {
		return nil, errors.New("__browserServices not defined in JavaScript")
	}

	// Get the service
	svc := browserServices.Get(service)
	if svc.IsUndefined() || svc.IsNull() {
		return nil, fmt.Errorf("browser service '%s' not registered", service)
	}

	// Get the method
	methodFunc := svc.Get(method)
	if methodFunc.IsUndefined() || methodFunc.IsNull() {
		return nil, fmt.Errorf("method '%s' not found on browser service '%s'", method, service)
	}

	// Call the method directly - this blocks until JS returns
	result := methodFunc.Invoke(string(request))

	if result.IsNull() || result.IsUndefined() {
		return nil, nil
	}

	// Check if result is an error object
	if result.Type() == js.TypeObject {
		errProp := result.Get("error")
		if !errProp.IsUndefined() && !errProp.IsNull() {
			return nil, errors.New(errProp.String())
		}
	}

	return []byte(result.String()), nil
}

// CallAsync makes an asynchronous call to a browser service.
// The JavaScript function should handle the callback pattern.
// Returns a Future that can be used to retrieve the result.
func (bc *BrowserServiceChannel) CallAsync(service, method string, request []byte) (chan *CallResponse, string) {
	callID := bc.NextCallID()
	responseCh := make(chan *CallResponse, 1)

	// Register pending call
	bc.mu.Lock()
	bc.pendingCalls[callID] = responseCh
	bc.mu.Unlock()

	// Call JavaScript async method
	js.Global().Call("__wasmCallBrowserMethodAsync", service, method, string(request), callID)

	return responseCh, callID
}

// GetPendingCallCount returns the number of pending async calls (for debugging).
func (bc *BrowserServiceChannel) GetPendingCallCount() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.pendingCalls)
}

// CallBrowserService is a generic helper for calling synchronous browser services.
// The browser method should return a value directly (not a Promise).
func CallBrowserService[TReq any, TResp any](channel *BrowserServiceChannel, ctx context.Context, serviceName, methodName string, req TReq) (TResp, error) {
	var resp TResp

	// If TResp is a pointer type, create a new instance
	respType := reflect.TypeOf(resp)
	if respType.Kind() == reflect.Ptr {
		respValue := reflect.New(respType.Elem())
		resp = respValue.Interface().(TResp)
	}

	// Marshal the request
	marshaller := GetGlobalMarshaller()
	reqMsg, ok := any(req).(proto.Message)
	if !ok {
		return resp, fmt.Errorf("request is not a proto message")
	}

	requestData, err := marshaller.Marshal(reqMsg, MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
		UseEnumNumbers:  false,
	})
	if err != nil {
		return resp, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make synchronous call to browser
	responseData, err := channel.CallSync(serviceName, methodName, requestData)
	if err != nil {
		return resp, err
	}

	// Unmarshal the response
	var respMsg proto.Message
	var isProtoMsg bool

	if respMsg, isProtoMsg = any(resp).(proto.Message); !isProtoMsg {
		respMsg, isProtoMsg = any(&resp).(proto.Message)
	}

	if !isProtoMsg {
		return resp, fmt.Errorf("response is not a proto message (type: %T)", resp)
	}

	if err := marshaller.Unmarshal(responseData, respMsg, UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}); err != nil {
		return resp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

// CallBrowserServiceAsync is a generic helper for calling async browser services.
// Returns a Future that can be used to retrieve the result when ready.
func CallBrowserServiceAsync[TReq any, TResp proto.Message](channel *BrowserServiceChannel, ctx context.Context, serviceName, methodName string, req TReq) *Future[TResp] {
	marshaller := GetGlobalMarshaller()

	// Marshal the request
	reqMsg, ok := any(req).(proto.Message)
	if !ok {
		// Return a Future that immediately errors
		responseCh := make(chan *CallResponse, 1)
		responseCh <- &CallResponse{Error: fmt.Errorf("request is not a proto message")}
		return &Future[TResp]{
			responseCh: responseCh,
			marshaller: marshaller,
		}
	}

	requestData, err := marshaller.Marshal(reqMsg, MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
		UseEnumNumbers:  false,
	})
	if err != nil {
		responseCh := make(chan *CallResponse, 1)
		responseCh <- &CallResponse{Error: fmt.Errorf("failed to marshal request: %w", err)}
		return &Future[TResp]{
			responseCh: responseCh,
			marshaller: marshaller,
		}
	}

	// Make async call to browser - returns the response channel
	responseCh, _ := channel.CallAsync(serviceName, methodName, requestData)

	return &Future[TResp]{
		responseCh: responseCh,
		marshaller: marshaller,
	}
}
