# VTProto Marshaller

The `VTProtoMarshaller` is a built-in marshaller that works with [vtprotobuf](https://github.com/planetscale/vtprotobuf)-generated code, providing TinyGo compatibility and improved performance.

## Quick Start

```go
package main

import "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"

func main() {
    // Use vtprotobuf marshaller with fallback
    wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshallerWithFallback())

    // Register your services...
}
```

## Two Modes

### With Fallback (Recommended)

```go
wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshallerWithFallback())
```

- **Tries vtprotobuf methods first** (MarshalJSON/UnmarshalJSON)
- **Falls back to standard encoding/json** if vtprotobuf methods are missing
- **Best for**: Mixed codebases, gradual migration, development

### Strict Mode

```go
wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshaller())
```

- **Only uses vtprotobuf methods** (MarshalJSON/UnmarshalJSON)
- **Fails if methods are missing**
- **Best for**: Production, TinyGo builds, maximum performance

## How It Works

The marshaller attempts these methods in order:

### For Marshaling:
1. `MarshalJSON() ([]byte, error)` - vtprotobuf JSON marshaler (if available)
2. `MarshalVT() ([]byte, error)` - vtprotobuf binary marshaler (if available)
3. `json.Marshal()` - Standard Go JSON (if fallback enabled)

### For Unmarshaling:
1. `UnmarshalJSON([]byte) error` - vtprotobuf JSON unmarshaler (if available)
2. `UnmarshalVT([]byte) error` - vtprotobuf binary unmarshaler (if available)
3. `json.Unmarshal()` - Standard Go JSON (if fallback enabled)

## Generating vtprotobuf Code

Generate your proto files with vtprotobuf:

```bash
# Install vtprotobuf plugin
go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest

# Generate with both standard proto and vtprotobuf
protoc \
  --go_out=. \
  --go-vtproto_out=. \
  --go-vtproto_opt=features=marshal+unmarshal+size \
  --go-wasmjs_out=. \
  your_service.proto
```

This generates methods like:
- `MarshalJSON() ([]byte, error)` - Fast JSON marshaling
- `UnmarshalJSON([]byte) error` - Fast JSON unmarshaling
- `MarshalVT() ([]byte, error)` - Fast binary marshaling
- `UnmarshalVT([]byte) error` - Fast binary unmarshaling

## Benefits

### 1. TinyGo Compatibility
- vtprotobuf uses code generation instead of reflection
- Works perfectly with TinyGo's limited reflection support

### 2. Performance
- Up to **10x faster** than protojson for large messages
- No reflection overhead at runtime
- Optimized generated code

### 3. Smaller Binary Size
- Code generation produces smaller binaries
- No reflection tables needed
- Better for WASM deployment

### 4. Backward Compatibility
- Works alongside standard protobuf messages
- Fallback mode supports gradual migration
- No changes to existing API

## Example: TinyGo Build

```go
// main.go
package main

import (
    "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
    pb "myapp/gen"
)

func main() {
    // Use vtprotobuf for TinyGo compatibility
    wasm.SetGlobalMarshaller(wasm.NewVTProtoMarshallerWithFallback())

    // Register services as usual
    exports := &pb.MyServicesExports{
        MyService: &MyServiceImpl{},
    }
    exports.RegisterAPI()

    select {} // Keep running
}
```

Build with TinyGo:

```bash
tinygo build -o myapp.wasm -target wasm ./main.go
```

## Testing

The marshaller includes comprehensive tests:

```bash
cd pkg/wasm
go test -v -run TestVTProto
```

Tests cover:
- ✅ Marshaling with vtprotobuf methods
- ✅ Unmarshaling with vtprotobuf methods
- ✅ Fallback to standard JSON
- ✅ Error handling
- ✅ JSON validity
- ✅ Roundtrip consistency

## When to Use

### Use VTProtoMarshaller when:
- ✅ Building with TinyGo
- ✅ Performance is critical
- ✅ Binary size matters (WASM)
- ✅ You want reflection-free code

### Stick with ProtojsonMarshaller when:
- ✅ Using standard Go (not TinyGo)
- ✅ Don't need maximum performance
- ✅ Want standard protojson compatibility
- ✅ Simplicity is more important

## Implementation Details

Located in `pkg/wasm/vtproto_marshaller.go`:

```go
type VTProtoMarshaller struct {
    fallbackToStdJSON bool
}

func NewVTProtoMarshaller() *VTProtoMarshaller
func NewVTProtoMarshallerWithFallback() *VTProtoMarshaller

func (v *VTProtoMarshaller) Marshal(m proto.Message, opts MarshalOptions) ([]byte, error)
func (v *VTProtoMarshaller) Unmarshal(data []byte, m proto.Message, opts UnmarshalOptions) error
```

## Limitations

1. **MarshalOptions support**: Some options may not be respected if vtprotobuf code was generated without corresponding flags
2. **Binary-only vtprotobuf**: If your vtprotobuf code only has `MarshalVT/UnmarshalVT` (binary), JSON conversion will be less efficient
3. **Mixed codebases**: Use fallback mode if you have both vtprotobuf and standard messages

## See Also

- [MARSHALLER_GUIDE.md](./MARSHALLER_GUIDE.md) - Complete marshaller guide
- [pkg/wasm/README_MARSHALLER.md](./pkg/wasm/README_MARSHALLER.md) - API documentation
- [vtprotobuf GitHub](https://github.com/planetscale/vtprotobuf) - vtprotobuf project
- [TinyGo WASM Guide](https://tinygo.org/docs/guides/webassembly/) - TinyGo WASM documentation

## Support

The `VTProtoMarshaller` is included in `protoc-gen-go-wasmjs` and requires no additional dependencies beyond:
- `google.golang.org/protobuf/proto` (already required)
- Your vtprotobuf-generated code

For issues or questions, please file an issue on the project repository.
