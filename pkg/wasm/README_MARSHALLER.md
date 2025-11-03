# Marshaller Package

This directory contains the marshaller abstraction layer for `protoc-gen-go-wasmjs`, which allows you to customize how protocol buffer messages are serialized to/from JSON.

## Files

- **`marshaller.go`** - Defines the `Marshaller`, `Unmarshaller`, and `ProtoMarshaller` interfaces
- **`protojson_marshaller.go`** - Default implementation using `google.golang.org/protobuf/encoding/protojson`
- **`marshaller_config.go`** - Global marshaller configuration and management
- **`browser_channel.go`** - Uses the marshaller for browser service calls

## Why Custom Marshallers?

The default `protojson` marshaller works great for standard Go environments, but has limitations:

1. **TinyGo Compatibility**: protojson uses reflection heavily, which is not fully supported in TinyGo
2. **Performance**: Alternative marshallers like vtprotobuf can be faster
3. **Size**: Custom marshallers can reduce binary size for WASM builds
4. **Flexibility**: Different environments may need different serialization strategies

## Quick Start

### Using the Default (protojson)

Nothing to configure - it works out of the box:

```go
package main

func main() {
    // protojson is used automatically
    // No configuration needed
}
```

### Using a Custom Marshaller

1. Implement the `ProtoMarshaller` interface:

```go
type MyMarshaller struct{}

func (m *MyMarshaller) Marshal(msg proto.Message, opts wasm.MarshalOptions) ([]byte, error) {
    // Your implementation
}

func (m *MyMarshaller) Unmarshal(data []byte, msg proto.Message, opts wasm.UnmarshalOptions) error {
    // Your implementation
}
```

2. Set it globally before registering services:

```go
func main() {
    wasm.SetGlobalMarshaller(&MyMarshaller{})

    // Now register your services...
}
```

## Interfaces

### ProtoMarshaller

Combined interface for both marshaling and unmarshaling:

```go
type ProtoMarshaller interface {
    Marshaller
    Unmarshaller
}
```

### Marshaller

Converts proto messages to JSON:

```go
type Marshaller interface {
    Marshal(m proto.Message, opts MarshalOptions) ([]byte, error)
}
```

### Unmarshaller

Converts JSON to proto messages:

```go
type Unmarshaller interface {
    Unmarshal(data []byte, m proto.Message, opts UnmarshalOptions) error
}
```

## Options

### MarshalOptions

Control how proto messages are converted to JSON:

```go
type MarshalOptions struct {
    UseProtoNames   bool  // false: use JSON names (camelCase), true: use proto names
    EmitUnpopulated bool  // true: include zero values, false: omit zero values
    UseEnumNumbers  bool  // false: use enum names, true: use enum numbers
}
```

**Recommended for JavaScript compatibility:**
- `UseProtoNames: false` - JavaScript expects camelCase
- `EmitUnpopulated: true` - Avoid undefined values in JavaScript
- `UseEnumNumbers: false` - String enum names are more readable

### UnmarshalOptions

Control how JSON is converted to proto messages:

```go
type UnmarshalOptions struct {
    DiscardUnknown bool  // true: ignore unknown fields, false: error on unknown fields
    AllowPartial   bool  // true: allow missing required fields, false: error on missing fields
}
```

**Recommended for robustness:**
- `DiscardUnknown: true` - Tolerate extra fields from JavaScript
- `AllowPartial: true` - Allow incomplete messages for better compatibility

## Usage in Generated Code

The generated WASM code automatically uses the global marshaller:

```go
// In generated code (wasm_exports.go)
func (exports *ServicesExports) MyMethod(this js.Value, args []js.Value) any {
    req := &pb.MyRequest{}
    marshaller := wasm.GetGlobalMarshaller()

    // Unmarshal request
    if err := marshaller.Unmarshal([]byte(requestJSON), req, wasm.UnmarshalOptions{
        DiscardUnknown: true,
        AllowPartial:   true,
    }); err != nil {
        return createJSResponse(false, fmt.Sprintf("Failed to parse: %v", err), nil)
    }

    // ... call service method ...

    // Marshal response
    responseJSON, err := marshaller.Marshal(resp, wasm.MarshalOptions{
        UseProtoNames:   false,
        EmitUnpopulated: true,
        UseEnumNumbers:  false,
    })

    return createJSResponse(true, "Success", json.RawMessage(responseJSON))
}
```

## Thread Safety

- `SetGlobalMarshaller()` and `GetGlobalMarshaller()` are protected by a mutex
- Safe to call from multiple goroutines
- Should be called once at application startup

## Implementation Examples

### Example 1: vtprotobuf Marshaller (for TinyGo)

```go
package vtmarshaller

import (
    "google.golang.org/protobuf/proto"
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
)

type VTProtoMarshaller struct{}

func (v *VTProtoMarshaller) Marshal(m proto.Message, opts wasm.MarshalOptions) ([]byte, error) {
    // Use vtprotobuf's generated MarshalJSON if available
    if marshaler, ok := m.(interface{ MarshalJSON() ([]byte, error) }); ok {
        return marshaler.MarshalJSON()
    }
    return nil, fmt.Errorf("message doesn't support vtprotobuf marshaling")
}

func (v *VTProtoMarshaller) Unmarshal(data []byte, m proto.Message, opts wasm.UnmarshalOptions) error {
    if unmarshaler, ok := m.(interface{ UnmarshalJSON([]byte) error }); ok {
        return unmarshaler.UnmarshalJSON(data)
    }
    return fmt.Errorf("message doesn't support vtprotobuf unmarshaling")
}
```

### Example 2: Logging Marshaller (for debugging)

```go
package debugmarshaller

import (
    "log"
    "google.golang.org/protobuf/proto"
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
)

type LoggingMarshaller struct {
    underlying wasm.ProtoMarshaller
}

func NewLoggingMarshaller(underlying wasm.ProtoMarshaller) *LoggingMarshaller {
    return &LoggingMarshaller{underlying: underlying}
}

func (l *LoggingMarshaller) Marshal(m proto.Message, opts wasm.MarshalOptions) ([]byte, error) {
    log.Printf("Marshaling %T with opts %+v", m, opts)
    data, err := l.underlying.Marshal(m, opts)
    if err != nil {
        log.Printf("Marshal error: %v", err)
    } else {
        log.Printf("Marshaled to %d bytes", len(data))
    }
    return data, err
}

func (l *LoggingMarshaller) Unmarshal(data []byte, m proto.Message, opts wasm.UnmarshalOptions) error {
    log.Printf("Unmarshaling %d bytes to %T with opts %+v", len(data), m, opts)
    err := l.underlying.Unmarshal(data, m, opts)
    if err != nil {
        log.Printf("Unmarshal error: %v", err)
    }
    return err
}
```

Usage:

```go
func main() {
    // Wrap default marshaller with logging
    base := wasm.NewProtojsonMarshaller()
    logging := debugmarshaller.NewLoggingMarshaller(base)
    wasm.SetGlobalMarshaller(logging)
}
```

## Testing Your Marshaller

```go
package mymarshaller_test

import (
    "testing"
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
    pb "myapp/gen/proto"
)

func TestMarshaller(t *testing.T) {
    marshaller := NewMyMarshaller()

    // Create test message
    original := &pb.MyMessage{
        Name:  "test",
        Value: 42,
    }

    // Marshal
    data, err := marshaller.Marshal(original, wasm.MarshalOptions{
        UseProtoNames:   false,
        EmitUnpopulated: true,
        UseEnumNumbers:  false,
    })
    if err != nil {
        t.Fatalf("Marshal failed: %v", err)
    }

    // Unmarshal
    result := &pb.MyMessage{}
    err = marshaller.Unmarshal(data, result, wasm.UnmarshalOptions{
        DiscardUnknown: true,
        AllowPartial:   true,
    })
    if err != nil {
        t.Fatalf("Unmarshal failed: %v", err)
    }

    // Verify
    if result.Name != original.Name || result.Value != original.Value {
        t.Errorf("Roundtrip mismatch: got %+v, want %+v", result, original)
    }
}
```

## Migration from Direct protojson Usage

If you have existing code that uses protojson directly, migration is straightforward:

### Before:

```go
import "google.golang.org/protobuf/encoding/protojson"

opts := protojson.MarshalOptions{
    UseProtoNames:   false,
    EmitUnpopulated: true,
}
data, err := opts.Marshal(msg)
```

### After:

```go
import "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"

marshaller := wasm.GetGlobalMarshaller()
data, err := marshaller.Marshal(msg, wasm.MarshalOptions{
    UseProtoNames:   false,
    EmitUnpopulated: true,
})
```

The generated code already uses the new interface, so you only need to update custom code.

## Performance Tips

1. **Reuse marshaller instances** - Don't create new marshallers for each operation
2. **Use appropriate options** - `EmitUnpopulated: false` can reduce JSON size for streaming
3. **Consider vtprotobuf** - Up to 10x faster than protojson for large messages
4. **Profile your marshaller** - Use Go profiling tools to identify bottlenecks

## See Also

- [MARSHALLER_GUIDE.md](../../MARSHALLER_GUIDE.md) - Complete guide with examples
- [vtprotobuf](https://github.com/planetscale/vtprotobuf) - Fast proto marshaller for TinyGo
- [TinyGo documentation](https://tinygo.org/docs/guides/webassembly/) - WASM with TinyGo
