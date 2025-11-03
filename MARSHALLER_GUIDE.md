# Custom Marshaller Guide

This guide explains how to use custom marshallers with `protoc-gen-go-wasmjs`, including how to use vtprotobuf for TinyGo compatibility.

## Overview

By default, `protoc-gen-go-wasmjs` uses `protojson` for marshaling and unmarshaling protocol buffer messages. While this works well in most Go environments, it may not be compatible with TinyGo due to reflection requirements.

To support TinyGo and other specialized environments, we've introduced a marshaller interface that allows you to swap out the serialization implementation.

## Architecture

The marshaller system consists of three main components:

1. **Marshaller Interface** (`pkg/wasm/marshaller.go`):
   - `Marshaller` - Interface for marshaling proto messages to JSON
   - `Unmarshaller` - Interface for unmarshaling JSON to proto messages
   - `ProtoMarshaller` - Combined interface

2. **Default Implementation** (`pkg/wasm/protojson_marshaller.go`):
   - `ProtojsonMarshaller` - Default implementation using protojson

3. **Configuration** (`pkg/wasm/marshaller_config.go`):
   - `SetGlobalMarshaller()` - Set the global marshaller
   - `GetGlobalMarshaller()` - Get the current marshaller

## Using the Default (protojson)

If you're using standard Go (not TinyGo), the default protojson marshaller works out of the box:

```go
package main

import (
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
)

func main() {
    // No configuration needed - protojson is the default
    // Your generated code will automatically use it
}
```

## Using the Built-in vtprotobuf Marshaller

We provide a built-in `VTProtoMarshaller` that works with vtprotobuf-generated code:

```go
package main

import (
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
)

func main() {
    // Use the built-in vtprotobuf marshaller with fallback to standard JSON
    wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshallerWithFallback())

    // Or use strict mode (fails if vtprotobuf methods are missing)
    // wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshaller())

    // Register your services...
}
```

**With fallback enabled** (recommended during migration):
- Tries vtprotobuf methods first
- Falls back to standard `encoding/json` if vtprotobuf methods are missing
- Good for mixed codebases

**Without fallback** (strict mode):
- Only uses vtprotobuf methods
- Fails if vtprotobuf methods are missing
- Ensures all messages use the fast path

## Creating a Custom Marshaller

If you need a fully custom marshaller, implement the `ProtoMarshaller` interface:

```go
package myapp

import (
    "google.golang.org/protobuf/proto"
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
)

// VTProtoMarshaller uses vtprotobuf for TinyGo compatibility
type VTProtoMarshaller struct{}

func NewVTProtoMarshaller() *VTProtoMarshaller {
    return &VTProtoMarshaller{}
}

// Marshal converts a proto message to JSON bytes using vtprotobuf
func (v *VTProtoMarshaller) Marshal(m proto.Message, opts wasm.MarshalOptions) ([]byte, error) {
    // Option 1: If your messages implement MarshalJSON from vtprotobuf
    if marshaler, ok := m.(interface{ MarshalJSON() ([]byte, error) }); ok {
        return marshaler.MarshalJSON()
    }

    // Option 2: Fall back to proto binary + custom JSON conversion
    // (You would implement your own JSON serialization here)
    return nil, fmt.Errorf("message does not support vtprotobuf JSON marshaling")
}

// Unmarshal parses JSON bytes into a proto message using vtprotobuf
func (v *VTProtoMarshaller) Unmarshal(data []byte, m proto.Message, opts wasm.UnmarshalOptions) error {
    // Option 1: If your messages implement UnmarshalJSON from vtprotobuf
    if unmarshaler, ok := m.(interface{ UnmarshalJSON([]byte) error }); ok {
        return unmarshaler.UnmarshalJSON(data)
    }

    // Option 2: Fall back to custom JSON parsing + proto binary
    // (You would implement your own JSON deserialization here)
    return fmt.Errorf("message does not support vtprotobuf JSON unmarshaling")
}

// Ensure VTProtoMarshaller implements ProtoMarshaller
var _ wasm.ProtoMarshaller = (*VTProtoMarshaller)(nil)
```

## Setting a Custom Marshaller

Set your custom marshaller **early in your application initialization**, before any service methods are called:

```go
package main

import (
    "syscall/js"
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
    "myapp/mymarshaller"
)

func main() {
    // Set custom marshaller FIRST
    wasm.SetGlobalMarshaller(mymarshaller.NewVTProtoMarshaller())

    // Then register your services
    exports := &myapp.MyServicesExports{
        MyService: myapp.NewMyServiceImpl(),
    }
    exports.RegisterAPI()

    // Keep the program running
    select {}
}
```

## vtprotobuf Integration Example

The built-in `VTProtoMarshaller` handles vtprotobuf integration automatically. Here's how to use it:

### 1. Generate Code with vtprotobuf

Add vtprotobuf to your proto generation:

```bash
# Install vtprotobuf plugin
go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest

# Generate with both protoc-gen-go-wasmjs and vtprotobuf
protoc \
  --go_out=. \
  --go-vtproto_out=. \
  --go-vtproto_opt=features=marshal+unmarshal+size \
  --go-wasmjs_out=. \
  your_service.proto
```

### 2. Configure Your WASM Application

Use the built-in `VTProtoMarshaller`:

```go
package main

import (
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
    pb "myapp/gen/myservice"
)

func main() {
    // Use built-in vtprotobuf marshaller for TinyGo compatibility
    // With fallback for any messages that don't have vtprotobuf methods
    wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshallerWithFallback())

    // Register services
    exports := &pb.MyServicesExports{
        MyService: &MyServiceImpl{},
    }
    exports.RegisterAPI()

    // Keep running
    select {}
}
```

### 3. Build with TinyGo

```bash
tinygo build -o myapp.wasm -target wasm ./main.go
```

## Testing Your Custom Marshaller

Create tests to ensure your marshaller works correctly:

```go
package mymarshaller_test

import (
    "testing"
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
    "myapp/mymarshaller"
    pb "myapp/gen/myservice"
)

func TestCustomMarshaller(t *testing.T) {
    marshaller := mymarshaller.NewVTProtoMarshaller()

    // Test marshaling
    msg := &pb.MyRequest{
        Name: "test",
        Value: 42,
    }

    data, err := marshaller.Marshal(msg, wasm.MarshalOptions{
        EmitUnpopulated: true,
    })
    if err != nil {
        t.Fatalf("Marshal failed: %v", err)
    }

    // Test unmarshaling
    msg2 := &pb.MyRequest{}
    err = marshaller.Unmarshal(data, msg2, wasm.UnmarshalOptions{
        DiscardUnknown: true,
    })
    if err != nil {
        t.Fatalf("Unmarshal failed: %v", err)
    }

    if msg2.Name != msg.Name || msg2.Value != msg.Value {
        t.Errorf("Roundtrip failed: got %+v, want %+v", msg2, msg)
    }
}
```

## Marshaller Options

Both `Marshal` and `Unmarshal` accept options to control serialization behavior:

### MarshalOptions

```go
type MarshalOptions struct {
    UseProtoNames   bool  // Use proto field names instead of camelCase
    EmitUnpopulated bool  // Include zero/default values
    UseEnumNumbers  bool  // Use enum integers instead of names
}
```

### UnmarshalOptions

```go
type UnmarshalOptions struct {
    DiscardUnknown bool  // Ignore unknown fields
    AllowPartial   bool  // Allow missing required fields
}
```

## Thread Safety

The global marshaller is protected by a mutex and is safe to use from multiple goroutines. However, you should set the marshaller only once during initialization:

```go
func main() {
    // ✅ Good: Set once at startup
    wasm.SetGlobalMarshaller(mymarshaller.New())

    // Start application...
}

// ❌ Bad: Don't change marshaller after services are registered
func someHandler() {
    wasm.SetGlobalMarshaller(differentMarshaller) // Dangerous!
}
```

## Troubleshooting

### Issue: "message does not support JSON marshaling"

**Solution**: Ensure your custom marshaller correctly implements the `Marshal` and `Unmarshal` methods, and that your proto messages are generated with the appropriate plugins (e.g., vtprotobuf with JSON support).

### Issue: TinyGo build fails with reflection errors

**Solution**: Make sure you're using vtprotobuf or another reflection-free marshaller. The default protojson uses reflection which is not fully supported in TinyGo.

### Issue: Different JSON format from JavaScript

**Solution**: Check your `MarshalOptions` - particularly `UseProtoNames` and `EmitUnpopulated`. JavaScript typically expects camelCase field names (`UseProtoNames: false`) and may need zero values included (`EmitUnpopulated: true`).

## Performance Considerations

- **protojson** (default): Good for standard Go, uses reflection
- **vtprotobuf**: Faster, no reflection, TinyGo-compatible, but requires code generation
- Custom marshallers: Performance depends on your implementation

## Summary

1. By default, `protoc-gen-go-wasmjs` uses protojson
2. For TinyGo, implement a custom `ProtoMarshaller` using vtprotobuf or similar
3. Set the marshaller early in your `main()` function using `wasm.SetGlobalMarshaller()`
4. All generated code will automatically use your custom marshaller
5. Test your marshaller thoroughly to ensure correct serialization

For more examples, see the `pkg/wasm/` directory in the repository.
