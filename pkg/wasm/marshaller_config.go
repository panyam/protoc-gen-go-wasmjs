// Package wasm provides runtime support for protoc-gen-go-wasmjs generated code.
package wasm

import (
	"sync"
)

var (
	// globalMarshaller is the default marshaller used by all generated code
	globalMarshaller ProtoMarshaller
	// marshallerMutex protects access to globalMarshaller
	marshallerMutex sync.RWMutex
)

func init() {
	// Default to protojson marshaller for backward compatibility
	globalMarshaller = NewProtojsonMarshaller()
}

// SetGlobalMarshaller sets the marshaller to be used by all generated WASM code.
// This should be called early in your application initialization, before any
// service methods are invoked.
//
// Example usage with protojson (default):
//
//	wasm.SetGlobalMarshaller(wasm.NewProtojsonMarshaller())
//
// Example usage with a custom marshaller (e.g., vtprotobuf):
//
//	wasm.SetGlobalMarshaller(myapp.NewVTProtoMarshaller())
//
// This function is safe to call from multiple goroutines, but should typically
// be called only once during application initialization.
func SetGlobalMarshaller(marshaller ProtoMarshaller) {
	marshallerMutex.Lock()
	defer marshallerMutex.Unlock()
	globalMarshaller = marshaller
}

// GetGlobalMarshaller returns the currently configured global marshaller.
// This is used internally by generated code and should rarely need to be
// called directly by application code.
func GetGlobalMarshaller() ProtoMarshaller {
	marshallerMutex.RLock()
	defer marshallerMutex.RUnlock()
	return globalMarshaller
}
