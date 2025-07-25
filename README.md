# protoc-gen-go-wasmjs

**protoc-gen-go-wasmjs** is a [Protocol Buffers](https://protobuf.dev) compiler plugin that generates WASM bindings and TypeScript clients for your gRPC services, enabling local-first applications that can run the same service logic in both server and browser environments.

It generates WASM wrapper Go code and TypeScript clients from your protobuf services, allowing you to deploy identical business logic as WebAssembly modules in the browser or as traditional gRPC servers.

## Features

- **Local-First Architecture**: Same service interface runs on server (full database) or browser (local storage)
- **Multi-Service WASM Modules**: Bundle related services in single WASM binaries
- **Flexible JavaScript APIs**: Choose between namespaced, flat, or service-based API structures
- **TypeScript Integration**: Works with existing protobuf TypeScript generators (protoc-gen-es, protoc-gen-ts)
- **Extensive Customization**: Method filtering, renaming, and template overrides
- **Build Pipeline Integration**: Seamless integration with buf and modern protobuf toolchains

## Quick Start

### Installation

**Option 1: Local Installation**
```bash
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest
```

**Option 2: Use from buf.build (Recommended)**
No installation required - use the remote plugin directly in your `buf.gen.yaml`.

### Basic Usage

**With Local Plugin:**
```yaml
plugins:
  # Generate standard Go protobuf types
  - remote: buf.build/protocolbuffers/go
    out: ./gen/go
    opt: paths=source_relative

  # Generate gRPC service interfaces  
  - remote: buf.build/grpc/go
    out: ./gen/go
    opt: paths=source_relative

  # Generate TypeScript protobuf types
  - remote: buf.build/bufbuild/es
    out: ./gen/ts
    opt: target=ts

  # Generate WASM wrapper and TypeScript client (local)
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=./gen/ts
      - js_structure=namespaced
      - js_namespace=library
```

**With Remote Plugin from buf.build:**
```yaml
plugins:
  # Generate standard Go protobuf types
  - remote: buf.build/protocolbuffers/go
    out: ./gen/go
    opt: paths=source_relative

  # Generate gRPC service interfaces  
  - remote: buf.build/grpc/go
    out: ./gen/go
    opt: paths=source_relative

  # Generate TypeScript protobuf types
  - remote: buf.build/bufbuild/es
    out: ./gen/ts
    opt: target=ts

  # Generate WASM wrapper and TypeScript client (remote)
  - remote: buf.build/panyam/protoc-gen-go-wasmjs
    out: ./gen/wasm
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=./gen/ts
      - js_structure=namespaced
      - js_namespace=library
```

### Example Service

```protobuf
syntax = "proto3";
package library.v1;

import "wasmjs/v1/annotations.proto";

service LibraryService {
  // Custom method name for cleaner JavaScript API
  rpc FindBooks(FindBooksRequest) returns (FindBooksResponse) {
    option (wasmjs.v1.wasm_method_name) = "searchBooks";
  }
  
  rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse);
}
```

This generates:

**Go WASM wrapper** (`library_v1_services.wasm.go`):
```go
//go:build js && wasm

// Namespaced JavaScript API: library.libraryService.searchBooks()
js.Global().Set("library", js.ValueOf(map[string]interface{}{
  "libraryService": map[string]interface{}{
    "searchBooks": js.FuncOf(libraryServiceFindBooks),
    "checkoutBook": js.FuncOf(libraryServiceCheckoutBook),
  },
}))
```

**TypeScript client** (`Library_v1_servicesClient.ts`):
```typescript
export class Library_v1_servicesClient {
  public readonly libraryService: LibraryServiceClientImpl;
  
  async loadWasm(wasmPath?: string): Promise<void> { /* ... */ }
}

class LibraryServiceClientImpl {
  async searchBooks(request: FindBooksRequest): Promise<FindBooksResponse> {
    return this.parent.callMethod('libraryService.searchBooks', request);
  }
}
```

## Configuration Options

### Core Integration

| Option | Description | Default |
|--------|-------------|---------|
| `ts_generator` | TypeScript generator used | `protoc-gen-es` |
| `ts_import_path` | Path to generated TS types | `./gen/ts` |
| `ts_export_path` | Where to generate TS client | `.` |
| `wasm_export_path` | Where to generate WASM wrapper | `.` |

### Service & Method Selection

| Option | Description | Example |
|--------|-------------|---------|
| `services` | Specific services to generate | `LibraryService,UserService` |
| `method_include` | Include methods by glob pattern | `Find*,Get*,Create*` |
| `method_exclude` | Exclude methods by glob pattern | `*Internal,*Debug` |
| `method_rename` | Rename methods | `FindBooks:searchBooks,GetUser:fetchUser` |

### JavaScript API Structure

| Option | Description | Result |
|--------|-------------|--------|
| `js_structure=namespaced` | Clean namespaced API | `myapp.service.method()` |
| `js_structure=flat` | Flat function names | `myappServiceMethod()` |
| `js_structure=service_based` | Service grouping | `services.library.findBooks()` |
| `js_namespace` | Global namespace name | Custom namespace |
| `module_name` | WASM module name | Custom module name |

### Advanced Customization

| Option | Description |
|--------|-------------|
| `template_dir` | Override default templates |
| `wasm_template` | Custom WASM template file |
| `ts_template` | Custom TypeScript template file |
| `generate_build_script` | Generate build.sh script |

## WASM Annotations

Customize generation behavior with protobuf annotations:

```protobuf
import "wasmjs/v1/annotations.proto";

service LibraryService {
  // Custom JavaScript method name
  rpc FindBooks(FindBooksRequest) returns (FindBooksResponse) {
    option (wasmjs.v1.wasm_method_name) = "searchBooks";
  }
  
  // Exclude from WASM generation
  rpc AdminMethod(AdminRequest) returns (AdminResponse) {
    option (wasmjs.v1.wasm_method_exclude) = true;
  }
}

// Exclude entire service
service InternalService {
  option (wasmjs.v1.wasm_service_exclude) = true;
  // ...methods
}

// Custom service name in JavaScript
service LibraryService {
  option (wasmjs.v1.wasm_service_name) = "books";
  // Results in: namespace.books.method() instead of namespace.libraryService.method()
}
```

## Local-First Use Case

The primary use case is enabling local-first applications where the same business logic runs in both environments:

**Server Environment** (Full Dataset):
```go
type LibraryService struct {
  db *sql.DB // Access to millions of books
}

func (s *LibraryService) FindBooks(ctx context.Context, req *FindBooksRequest) (*FindBooksResponse, error) {
  // Query full database
  return s.searchDatabase(req.Query)
}
```

**Browser Environment** (Local Subset):
```go  
type LibraryService struct {
  books []Book // Local subset from localStorage
}

func (s *LibraryService) FindBooks(ctx context.Context, req *FindBooksRequest) (*FindBooksResponse, error) {
  // Search local books only
  return s.searchLocalBooks(req.Query)
}
```

**Frontend Code** (Same Interface):
```typescript
// Can switch between local WASM or remote HTTP seamlessly
const client = new LibraryServicesClient();
await client.loadWasm('./library.wasm');

// Same API regardless of backend
const books = await client.libraryService.searchBooks({ 
  query: "golang", 
  limit: 10 
});
```

## Build Process

Generated files include a build script:

```bash
# Generated build.sh
#!/bin/bash
export GOOS=js GOARCH=wasm
go build -o library_v1_services.wasm library_v1_services.wasm.go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

Integration in web applications:

```html
<script src="wasm_exec.js"></script>
<script>
  // Initialize WASM
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("library.wasm"), go.importObject)
    .then(result => {
      go.run(result.instance);
      
      // Use generated TypeScript client
      const client = new LibraryServicesClient();
      // WASM is already loaded, so this returns immediately
      await client.waitUntilReady();
      
      const books = await client.libraryService.searchBooks({
        query: "javascript",
        limit: 5
      });
    });
</script>
```

## TypeScript Generator Compatibility

Works with popular TypeScript protobuf generators:

**protoc-gen-es** (Recommended):
```typescript
// Auto-detects .toJson() and .fromJson() methods
const response = await client.method(request);
```

**protoc-gen-ts**:
```typescript  
// Auto-detects .toJSON() and fromJSON() functions
const response = await client.method(request);
```

**Generic Fallback**:
```typescript
// Falls back to JSON.stringify/parse
const response = await client.method(request);
```

## Installation & Usage Notes

### Using the Local Plugin

1. **Install the plugin:**
   ```bash
   go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest
   ```

2. **Ensure the plugin is in your PATH:**
   ```bash
   which protoc-gen-go-wasmjs
   # Should output: /path/to/go/bin/protoc-gen-go-wasmjs
   ```

3. **Use `local:` in your buf.gen.yaml:**
   ```yaml
   - local: protoc-gen-go-wasmjs
   ```

### Using the Remote Plugin (buf.build)

1. **No installation required** - buf automatically downloads and runs the plugin

2. **Use `remote:` in your buf.gen.yaml:**
   ```yaml
   - remote: buf.build/panyam/protoc-gen-go-wasmjs
   ```

3. **Benefits of remote plugins:**
   - No local installation required
   - Always uses the latest version
   - Consistent across team members
   - Works in CI/CD without additional setup

### Publishing to buf.build

To publish this plugin to buf.build (for maintainers):

1. **Create a buf.plugin.yaml:**
   ```yaml
   version: v1
   name: buf.build/panyam/protoc-gen-go-wasmjs
   plugin_version: v1.0.0
   description: Generate WASM bindings and TypeScript clients for gRPC services
   ```

2. **Push to buf.build:**
   ```bash
   buf plugin push
   ```

## Project Structure

```
├── cmd/protoc-gen-go-wasmjs/     # Plugin entry point
├── pkg/generator/                # Code generation logic
│   ├── templates/                # Embedded template files
│   │   ├── wasm.go.tmpl         # Go WASM wrapper template
│   │   ├── client.ts.tmpl       # TypeScript client template
│   │   └── build.sh.tmpl        # Build script template
│   ├── config.go                # Configuration parsing
│   ├── generator.go             # Main generation logic
│   └── types.go                 # Template data structures
├── proto/wasmjs/v1/             # WASM annotation definitions
├── example/                     # LibraryService example
└── PLAN.md                      # Development progress tracking
```

## Contributing

1. Fork the repository and create a feature branch
2. Make changes and add tests
3. Ensure all tests pass: `go test ./...`
4. Update documentation as needed
5. Submit a pull request

## Development

```bash
# Build the plugin
go build ./cmd/protoc-gen-go-wasmjs

# Test with example
cd example && buf generate

# Run tests
go test ./...
```

## License

Licensed under the Apache License, Version 2.0.