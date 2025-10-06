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

/*
Package wasm provides runtime utilities for WebAssembly execution and browser communication.

# Overview

The wasm package implements the runtime infrastructure for WASM-to-JavaScript communication,
browser service channels, and protobuf serialization in WebAssembly environments. This package
is only compiled when GOOS=js and GOARCH=wasm.

# Build Tags

All files in this package use the build constraint:

	//go:build js && wasm

This ensures the package is only included when building for WebAssembly targets.

# Core Components

BrowserServiceChannel

Manages bidirectional communication between WASM and browser-provided services.
This enables WASM code to call browser APIs (localStorage, fetch, etc.) without deadlocking
the main thread.

Key Features:

  - FIFO call queue for ordered execution
  - Timeout support with automatic cleanup
  - Async method support (prevents browser deadlocks)
  - Error handling and propagation
  - Singleton pattern for global access

BrowserCall

Represents a single call from WASM to a browser-provided service.
Contains request data, response channel, and metadata for tracking.

CallResponse

Represents the response from a browser service call.
Contains either response data or an error.

# Browser Service Communication

The browser service channel implements a sophisticated async communication system:

	WASM Side (Go):
	1. Create BrowserCall with request data
	2. Send call to browserChannel.CallBrowserService()
	3. Wait on response channel (or use callback for async)
	4. Process response or error

	JavaScript Side:
	1. Poll __wasmGetNextBrowserCall()
	2. Execute browser service method
	3. Call __wasmHandleBrowserResponse(id, data) or __wasmHandleBrowserError(id, error)
	4. WASM receives response and unblocks

# Usage Example

Browser Service Interface (Proto):

	service BrowserAPI {
	    option (wasmjs.v1.browser_provided) = true;

	    rpc GetLocalStorage(StorageKeyRequest) returns (StorageValueResponse);
	    rpc Fetch(FetchRequest) returns (FetchResponse);
	}

Generated WASM Code:

	import (
	    "context"
	    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
	)

	// Call browser service from WASM
	func callBrowserLocalStorage(ctx context.Context, key string) (string, error) {
	    // Create request
	    req := &StorageKeyRequest{Key: key}
	    reqData, _ := proto.Marshal(req)

	    // Get browser channel
	    browserChannel := wasm.GetBrowserChannel()

	    // Create browser call
	    call := &wasm.BrowserCall{
	        Service:    "BrowserAPI",
	        Method:     "GetLocalStorage",
	        Request:    reqData,
	        ResponseCh: make(chan *wasm.CallResponse, 1),
	        Timeout:    5 * time.Second,
	    }

	    // Execute call (blocks until response)
	    respData, err := browserChannel.CallBrowserService(ctx, call)
	    if err != nil {
	        return "", err
	    }

	    // Unmarshal response
	    var resp StorageValueResponse
	    if err := proto.Unmarshal(respData, &resp); err != nil {
	        return "", err
	    }

	    return resp.Value, nil
	}

JavaScript Implementation:

	// Register browser service implementation
	const browserServiceManager = new BrowserServiceManager();
	browserServiceManager.registerService('BrowserAPI', {
	    async getLocalStorage(request) {
	        const value = localStorage.getItem(request.key);
	        return { value: value || '', exists: value !== null };
	    },
	    async fetch(request) {
	        const response = await fetch(request.url);
	        return { body: await response.text(), status: response.status };
	    }
	});

	// Browser service manager automatically polls for calls
	// and responds via __wasmHandleBrowserResponse()

# Async Method Support

Methods marked with (wasmjs.v1.async_method) = { is_async: true } use callbacks
to prevent browser main thread deadlocks:

	// Async method with callback (does not block)
	func runAsyncOperation(callback func(response []byte, err error)) {
	    call := &wasm.BrowserCall{
	        Service: "BrowserAPI",
	        Method:  "LongRunningOperation",
	        IsAsync: true,
	        // ...
	    }

	    // Execute asynchronously with callback
	    go func() {
	        resp, err := browserChannel.CallBrowserService(context.Background(), call)
	        callback(resp, err)
	    }()
	}

This prevents the WASM code from blocking the JavaScript main thread while waiting
for user input or long-running operations.

# Protobuf Serialization

The package provides helpers for protobuf serialization in WASM:

	// Serialize protobuf message to JSON (for JavaScript consumption)
	func SerializeProtoToJSON(msg proto.Message) ([]byte, error) {
	    return protojson.Marshal(msg)
	}

	// Deserialize JSON to protobuf message (from JavaScript)
	func DeserializeProtoFromJSON(data []byte, msg proto.Message) error {
	    return protojson.Unmarshal(data, msg)
	}

# Thread Safety

The BrowserServiceChannel is thread-safe:

  - Uses sync.RWMutex for concurrent access to pending calls map
  - Channel-based communication for FIFO ordering
  - Atomic operations for call ID generation
  - Safe cleanup of expired calls

# Error Handling

The package defines clear error types:

  - Context errors (timeout, cancellation)
  - Browser service errors (service not found, method error)
  - Serialization errors (protobuf marshaling/unmarshaling)

All errors are propagated back to the caller with detailed messages.

# Performance Considerations

  - Minimal allocations in hot paths
  - Pooled response channels
  - Efficient timeout handling with time.Timer
  - Early cleanup of completed calls

# Testing

The package includes comprehensive tests:

  - browser_channel_test.go: Channel operations, timeouts, cleanup
  - protobuf_deserialization_test.go: Serialization edge cases
  - All tests use mock browser environments

Run tests:

	GOOS=js GOARCH=wasm go test ./pkg/wasm/...

# Debugging

Enable debug logging:

	// In WASM code
	log.SetOutput(os.Stderr)
	log.Printf("Browser call: service=%s method=%s", call.Service, call.Method)

	// In JavaScript console
	console.log('WASM call received:', callId, serviceName, methodName);

# Links

Related packages:

  - github.com/panyam/protoc-gen-go-wasmjs/pkg/generators: Generates WASM wrappers
  - github.com/panyam/protoc-gen-go-wasmjs/runtime (NPM): JavaScript runtime utilities
*/
package wasm
