# protoc-gen-go-wasmjs

**protoc-gen-go-wasmjs** is a [Protocol Buffers](https://protobuf.dev) compiler plugin that generates WASM bindings and TypeScript clients for your gRPC services, enabling local-first applications that can run the same service logic in both server and browser environments.

It generates flexible WASM exports and TypeScript clients from your protobuf services, allowing you to deploy identical business logic as WebAssembly modules in the browser or as traditional gRPC servers with full dependency injection control.

## Features

- **Dual-Target Architecture**: Generate WASM and TypeScript artifacts separately for flexible deployment
- **Smart Import Detection**: Automatically analyzes proto files to generate accurate TypeScript imports
- **Auto-Extension Detection**: Automatically detects `.ts` vs `.js` extensions based on protoc-gen-es configuration
- **Multi-Target Generation**: Generate optimized WASM bundles per page/use case (user page, admin page, etc.)
- **Dependency Injection**: Full control over service initialization with database, auth, config injection
- **Optimized Bundles**: Each target includes only the services it needs for smaller bundle sizes
- **Local-First Architecture**: Same service interface runs on server (full database) or browser (local storage)
- **Export Pattern**: Generates reusable exports instead of fixed main() for maximum flexibility
- **TypeScript Integration**: Works with existing protobuf TypeScript generators (protoc-gen-es, protoc-gen-ts)
- **Flexible Deployment**: TypeScript clients can be placed directly in frontend source directories
- **Extensive Customization**: Method filtering, renaming, and service targeting
- **Build Pipeline Integration**: Seamless integration with buf and modern protobuf toolchains

## Quick Start

### Installation

**Option 1: Local Installation**
```bash
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest
```

**Option 2: Use from buf.build (Recommended)**
No installation required - use the remote plugin directly in your `buf.gen.yaml`.

**Option 3: Split Generators (Latest)**
Install language-specific generators for focused generation:
```bash
# Install Go generator only
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-go@latest

# Install TypeScript generator only  
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-ts@latest
```

### Runtime Package Installation

Generated TypeScript code requires the runtime utilities package:

```bash
npm install @protoc-gen-go-wasmjs/runtime
# or
pnpm add @protoc-gen-go-wasmjs/runtime
# or  
yarn add @protoc-gen-go-wasmjs/runtime
```

## Architecture Patterns

### Dual-Target Architecture (Most Flexible)

Generate WASM and TypeScript artifacts separately for maximum deployment flexibility:

```yaml
plugins:
  # Standard protobuf generation...
  
  # WASM wrapper only - optimized for server-side deployment
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm/user-services
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=../../../gen/ts  # Relative to out directory
      - services=UsersService
      - generate_typescript=false  # Only generate WASM
      
  # TypeScript client only - deploy directly to frontend
  - local: protoc-gen-go-wasmjs
    out: ./web/frontend/src/wasm-clients
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=../../../gen/ts  # Relative to out directory
      - services=UsersService
      - generate_wasm=false  # Only generate TypeScript
```

**Benefits:**
- **Flexible placement**: TypeScript clients can go directly into frontend source directories
- **Clean separation**: WASM and TypeScript artifacts in completely different locations
- **Independent generation**: Generate just WASM, just TypeScript, or both as needed
- **Standard buf patterns**: Each target uses native protoc `out` directories

### Multi-Target Usage (Co-located)

Generate optimized WASM bundles per page/use case:

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

  # Generate TypeScript protobuf types (shared)
  - remote: buf.build/bufbuild/es
    out: ./web/frontend/gen
    opt: target=ts

  # User page target (UsersService only)
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm/user-page
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=web/frontend/gen
      - services=UsersService
      - module_name=user_page_services
      - js_namespace=userPage

  # Game page target (GamesService + WorldsService)
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm/game-page
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=web/frontend/gen
      - services=GamesService,WorldsService
      - module_name=game_page_services
      - js_namespace=gamePage

  # Admin page target (all services)
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm/admin-page
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=web/frontend/gen
      - module_name=admin_services
      - js_namespace=admin
```

### Single Target Usage (Simple)

For simple projects with one WASM module:

```yaml
plugins:
  # Standard protobuf generation...
  
  # Single WASM target
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm
    opt:
      - ts_generator=protoc-gen-es
      - ts_import_path=./gen/ts
      - js_structure=namespaced
      - js_namespace=myapp
```

### Using Generated Exports (Dependency Injection)

After running `buf generate`, each target generates:
- `{module_name}.wasm.go` - Importable WASM package with export struct
- `main.go.example` - Template showing how to use the exports
- `{module_name}Client.client.ts` - TypeScript client

**Step 1**: Copy and customize the `main.go.example`:

```go
// cmd/user-page-wasm/main.go
package main

import (
    "your-project/gen/wasm/user-page/user_page_services"
    libraryv1 "your-project/gen/go/library/v1"
)

func main() {
    // Initialize with your service implementations
    exports := &user_page_services.User_page_servicesServicesExports{
        UsersService: &myUserService{
            db: database,
            auth: authService,
            cache: redis,
        },
    }
    
    // Register JavaScript API
    exports.RegisterAPI()
    
    // Keep WASM running
    select {}
}
```

**Step 2**: Build the WASM binary:

```bash
cd cmd/user-page-wasm
GOOS=js GOARCH=wasm go build -o user_page.wasm
```

**Step 3**: Use in browser with TypeScript client:

```typescript
import { User_page_servicesClient } from './gen/wasm/user-page/user_page_servicesClient';

const client = new User_page_servicesClient();
await client.loadWasm('./user_page.wasm');

// Clean API - only UsersService methods available
const user = await client.usersService.getUser({ id: "123" });
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

## Proto to JSON Conversion

The plugin includes a flexible proto to JSON conversion system to handle differences between Go's protojson and TypeScript protobuf libraries. See [PROTO_CONVERSION.md](PROTO_CONVERSION.md) for detailed documentation.

### Quick Example
```typescript
// Create client with custom conversion options
const client = new MyServicesClient({
    handleOneofs: true,      // Flatten oneof fields for Go compatibility
    emitDefaults: false,     // Don't send default values
    fieldTransformer: (field) => {
        // Convert camelCase to snake_case if needed
        return field.replace(/([A-Z])/g, '_$1').toLowerCase();
    }
});
```

## Configuration Options

### Core Integration

| Option | Description | Default |
|--------|-------------|---------|
| `ts_generator` | TypeScript generator used | `protoc-gen-es` |
| `ts_import_path` | Path to generated TS types (relative to out dir) | `./gen/ts` |
| `ts_import_extension` | Extension for TS imports (`js`, `ts`, `none`, or empty for auto-detect) | auto-detect |
| `wasm_export_path` | Where to generate WASM wrapper | `.` |
| `generate_wasm` | Generate WASM wrapper | `true` |
| `generate_typescript` | Generate TypeScript client | `true` |

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
  
  // Async method with callback support (prevents browser deadlocks)
  rpc LoadData(LoadDataRequest) returns (LoadDataResponse) {
    option (wasmjs.v1.async_method) = { is_async: true };
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
// Import from runtime package for shared utilities
import { WASMResponse, WasmError } from '@protoc-gen-go-wasmjs/runtime';

// Can switch between local WASM or remote HTTP seamlessly
const client = new LibraryServicesClient();
await client.loadWasm('./library.wasm');

// Same API regardless of backend
const books = await client.libraryService.searchBooks({ 
  query: "golang", 
  limit: 10 
});
```

## Runtime Package (@protoc-gen-go-wasmjs/runtime)

Generated TypeScript code imports shared utilities from the runtime package, reducing bundle size and improving maintainability:

### **Key Components**

- **`WASMServiceClient`**: Base class for all generated WASM clients with streaming support
- **`BrowserServiceManager`**: Handles browser-provided service calls from WASM  
- **`BaseDeserializer`**: Schema-aware deserialization with cross-package support
- **`BaseSchemaRegistry`**: Utility methods for protobuf schema operations

### **Benefits**

- **🚀 Smaller bundles**: Shared utilities eliminate code duplication
- **🔧 Better maintenance**: Runtime fixes benefit all projects immediately
- **📦 Tree-shakeable**: Import only the utilities you need
- **🎯 Type safety**: Full TypeScript support with complete definitions

### **Usage**

```typescript
// Generated clients automatically import runtime utilities
import { MyServiceClient } from './generated/my_service_client';

// Manual usage (advanced scenarios)
import { 
  WASMServiceClient, 
  BaseDeserializer,
  FieldType 
} from '@protoc-gen-go-wasmjs/runtime';
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

Works with popular TypeScript protobuf generators and automatically detects optimal import settings:

**protoc-gen-es** (Recommended):
```yaml
# buf.gen.yaml - TypeScript generation
- remote: buf.build/bufbuild/es
  out: ./gen/ts
  opt: target=ts  # Generates .ts files

# WASM client generation
- local: protoc-gen-go-wasmjs
  out: ./web/frontend/src/clients
  opt:
    - ts_generator=protoc-gen-es
    - ts_import_path=../../../gen/ts
    # ts_import_extension auto-detected as "none" for .ts files
```

```typescript
// Generated client with smart imports:
import { CreateGameRequest, CreateGameResponse } from './games_pb';
import { CreateUserRequest, CreateUserResponse } from './users_pb';

// Auto-detects .toJson() and .fromJson() methods
const response = await client.gamesService.createGame(request);
```

**protoc-gen-ts**:
```typescript  
// Auto-detects .toJSON() and fromJSON() functions
const response = await client.method(request);
```

**Manual Extension Control**:
```yaml
# Force specific extension behavior
- local: protoc-gen-go-wasmjs
  opt:
    - ts_import_extension=js    # Force .js imports
    - ts_import_extension=ts    # Force .ts imports  
    - ts_import_extension=none  # No extension (TypeScript default)
```

## Smart Import Detection

The plugin automatically analyzes your proto files to generate accurate TypeScript imports:

**Before (Hardcoded)**:
```typescript
// Everything imported from models_pb regardless of actual source
import { CreateGameRequest, CreateUserRequest, CreateWorldRequest } from './models_pb';
```

**After (Smart Detection)**:
```typescript  
// Types imported from their actual proto file sources
import { CreateGameRequest, UpdateGameRequest } from './games_pb';
import { CreateUserRequest, UpdateUserRequest } from './users_pb';
import { CreateWorldRequest, UpdateWorldRequest } from './worlds_pb';
```

**How it works:**
1. **Proto File Analysis**: For each gRPC method, analyzes `method.Input.Desc.ParentFile().Path()` to determine which `.proto` file defines each type
2. **Automatic Grouping**: Groups types by source proto file for clean, organized imports
3. **Extension Detection**: Automatically detects whether protoc-gen-es generated `.ts` or `.js` files and adjusts import paths accordingly
4. **Zero Configuration**: Works out of the box with any proto file structure - no manual configuration needed

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
├── examples/                    # Comprehensive examples
│   ├── library/                 # LibraryService example
│   └── connect4/                # Real-time multiplayer Connect4 with stateful proxies
└── PLAN.md                      # Development progress tracking
```

## Development

For detailed development instructions, testing guidelines, and contribution workflows, see [DEVELOPMENT.md](DEVELOPMENT.md).

### Quick Start

```bash
# Run the test suite
./test.sh

# Build the plugin  
make tool

# Test with examples
cd examples/library && buf generate
cd examples/connect4 && make all

# Run core utility tests
go test ./pkg/core/... -v
```

## Contributing

1. Fork the repository and create a feature branch
2. Make changes and add comprehensive tests
3. Run the test suite: `./test.sh`
4. Update documentation as needed
5. Submit a pull request

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed contribution guidelines, testing requirements, and code quality standards.

## License

Licensed under the Apache License, Version 2.0.