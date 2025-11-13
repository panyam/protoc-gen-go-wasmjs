# protoc-gen-go-wasmjs

**protoc-gen-go-wasmjs** is a [Protocol Buffers](https://protobuf.dev) compiler plugin that generates WASM bindings and TypeScript clients for your gRPC services, enabling local-first applications that can run the same service logic in both server and browser environments.

It generates flexible WASM exports and TypeScript clients from your protobuf services, allowing you to deploy identical business logic as WebAssembly modules in the browser or as traditional gRPC servers with full dependency injection control.

## Features

- **BaseGenerator Architecture**: 4-step artifact processing approach separating collection from file generation
- **Composition-Based Bundles**: Simple base bundle classes with user-controlled service composition
- **Per-Service Client Generation**: Individual service clients following proto directory structure
- **Cross-Package Artifact Visibility**: Complete artifact catalog regardless of protoc Generate flags
- **Flexible File Mapping**: Generator-specific logic for N:1 and 1:N artifact-to-file mapping
- **Dual-Target Architecture**: Generate WASM and TypeScript artifacts with shared BaseGenerator foundation
- **Full TypeScript Type Safety**: Automatically generates typed interfaces with proper import resolution
- **Runtime Package Integration**: Clean inheritance-based architecture with shared utilities
- **Browser Service Communication**: Seamless WASM to browser API integration with async support
- **Dependency Injection**: Full control over service initialization with database, auth, config injection
- **Local-First Architecture**: Same service interface runs on server (full database) or browser (local storage)
- **Build Pipeline Integration**: Seamless integration with buf and modern protobuf toolchains

## Quick Start

### Installation

The plugin consists of two focused generators that can be used together or independently:

**Install Generators:**
```bash
# Install Go generator for WASM wrapper generation
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-go@latest

# Install TypeScript generator for client generation
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-ts@latest
```

**Verify Installation:**
```bash
which protoc-gen-go-wasmjs-go
which protoc-gen-go-wasmjs-ts
```

**Alternative: Use from buf.build**
No installation required - use the remote plugin directly in your `buf.gen.yaml` (when published).

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

### Modern Split-Generator Architecture

The plugin uses a split architecture with dedicated generators for each target language:

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

  # Generate Go WASM wrappers
  - local: protoc-gen-go-wasmjs-go
    out: ./gen/wasm/go
    opt:
      - js_structure=namespaced
      - js_namespace=myApp
      - module_name=my_services

  # Generate TypeScript clients (per-service)
  - local: protoc-gen-go-wasmjs-ts
    out: ./web/src/generated
    opt:
      - js_structure=namespaced
      - js_namespace=myApp
```

**Benefits:**
- **Bundle-based loading**: Single WASM load per module, shared across all services in that module
- **Individual service clients**: Each service gets its own typed client interface within the bundle
- **Full type safety**: Generated clients use proper TypeScript types with IntelliSense support
- **Efficient resource usage**: No duplicate WASM loading for services in the same module
- **Clean organization**: Bundle manages WASM lifecycle, service clients handle business logic
- **Runtime integration**: Uses WASMBundle and ServiceClient base classes from @protoc-gen-go-wasmjs/runtime package

### Browser Service Integration

For services that need to call browser APIs (localStorage, fetch, etc.):

```protobuf
// browser/v1/browser.proto
service BrowserAPI {
    option (wasmjs.v1.browser_provided) = true;
    
    rpc GetLocalStorage(StorageKeyRequest) returns (StorageValueResponse);
    rpc SetLocalStorage(StorageSetRequest) returns (StorageSetResponse); 
    rpc Alert(AlertRequest) returns (AlertResponse);
}
```

```typescript
// New composition-based architecture
import { ExampleBundle } from './generated';
import { PresenterServiceClient } from './generated/presenter/v1/presenterServiceClient';
import { BrowserAPIServiceClient } from './generated/browser/v1/browserAPIClient';

// Create base bundle with module configuration
const wasmBundle = new ExampleBundle();

// Create service clients using composition
const presenterService = new PresenterServiceClient(wasmBundle);
const browserAPI = new BrowserAPIServiceClient(wasmBundle);

// Register browser service implementations
wasmBundle.registerBrowserService('BrowserAPI', {
  async getLocalStorage(request) {
    return { value: localStorage.getItem(request.key) || '', exists: true };
  },
  async setLocalStorage(request) {
    localStorage.setItem(request.key, request.value);
    return { success: true };
  },
  async alert(request) {
    alert(request.message);
    return { shown: true };
  }
});

// Load WASM once for all services in this module
await wasmBundle.loadWasm('./my_module.wasm');

// Use individual service clients (all share the same WASM)
await presenterService.loadUserData({ userId: '123' });
```

## Generated File Structure

The generators create clean, organized file structures following proto package hierarchy:

```
web/src/generated/
├── index.ts                           # Base bundle class (module-level)
├── presenter/v1/
│   ├── presenterServiceClient.ts      # Service client (package-level)
│   └── interfaces.ts                  # TypeScript interfaces
└── browser/v1/
    ├── browserAPIClient.ts            # Browser service client
    └── interfaces.ts                  # TypeScript interfaces

gen/wasm/go/
├── presenter/v1/
│   ├── presenter_v1.wasm.go          # WASM wrapper
│   └── main.go.example               # Usage example
└── browser/v1/
    ├── browser_v1.wasm.go
    └── main.go.example
```

**Key Architecture:**
- **Base Bundle**: Simple class extending WASMBundle with module configuration
- **Service Clients**: Individual clients per service following proto structure
- **User Composition**: Users choose which services to include
- **Package Organization**: Mirrors proto package structure for clarity

### Using Generated Exports (Dependency Injection)

After running `buf generate`, each target generates:
- `{module_name}.wasm.go` - Importable WASM package with export struct
- `main.go.example` - Template showing how to use the exports
- `{serviceName}Client.ts` - TypeScript bundle and service clients

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

**Step 3**: Use in browser with TypeScript bundle client:

```typescript
import { User_page_servicesBundle } from './gen/wasm/user-page/user_page_servicesClient';

// Create bundle - manages WASM loading for all services in this module
const bundle = new User_page_servicesBundle();
await bundle.loadWasm('./user_page.wasm');

// Access individual service clients - all share the same WASM instance
const user = await bundle.usersService.getUser({ id: "123" });
const profile = await bundle.profileService.getProfile({ userId: "123" });
```

### Example Service with Full Type Safety (from example)

```protobuf
syntax = "proto3";
package presenter.v1;

import "wasmjs/v1/annotations.proto";

service PresenterService {
    // Regular sync method
    rpc LoadUserData(LoadUserRequest) returns (LoadUserResponse);

    // Streaming method for real-time updates
    rpc UpdateUIState(StateUpdateRequest) returns (stream UIUpdate);

    // Async method for long-running operations  
    rpc RunCallbackDemo(CallbackDemoRequest) returns (CallbackDemoResponse) {
        option (wasmjs.v1.async_method) = { is_async: true };
    };
}
```

This generates:

**Per-service TypeScript client** (`presenter/v1/presenterServiceClient.ts`):
```typescript
import { WASMBundle, WASMBundleConfig, ServiceClient } from '@protoc-gen-go-wasmjs/runtime';
import {
  LoadUserRequest,
  LoadUserResponse,
  StateUpdateRequest,
  UIUpdate,
  CallbackDemoRequest,
  CallbackDemoResponse,
} from './interfaces';

export class ExampleBundle {
  private wasmBundle: WASMBundle;
  public readonly presenterService: PresenterServiceClient;
  
  constructor() {
    const config: WASMBundleConfig = {
      moduleName: 'example',
      apiStructure: 'namespaced',
      jsNamespace: 'example'
    };
    this.wasmBundle = new WASMBundle(config);
    this.presenterService = new PresenterServiceClient(this.wasmBundle);
  }
  
  async loadWasm(wasmPath: string): Promise<void> { /* ... */ }
}

export class PresenterServiceClient extends ServiceClient implements PresenterServiceMethods {
  // Fully typed sync method
  async loadUserData(request: LoadUserRequest): Promise<LoadUserResponse> {
    return this.callMethod('presenterService.loadUserData', request);
  }
  
  // Fully typed streaming method
  updateUIState(
    request: StateUpdateRequest,
    callback: (response: UIUpdate | null, error: string | null, done: boolean) => boolean
  ): void {
    return this.callStreamingMethod('presenterService.updateUIState', request, callback);
  }
  
  // Fully typed async method with callback
  async runCallbackDemo(
    request: CallbackDemoRequest, 
    callback: (response: CallbackDemoResponse, error?: string) => void
  ): Promise<void> {
    return this.callMethodWithCallback('presenterService.runCallbackDemo', request, callback);
  }
}
```

## TypeScript Generation Model

The TypeScript generator creates a complete set of files per proto package, following a clean separation between interfaces and implementations:

### Generated Files Per Package

For each proto package (e.g., `utils.v1`), the generator creates:

**`interfaces.ts`** - Pure TypeScript interfaces for type safety:
```typescript
export interface NestedUtilType {
  topLevelCount: number;
  topLevelValue: string;
}
```

**`models.ts`** - Concrete implementations with default values:
```typescript
export class NestedUtilType implements NestedUtilTypeInterface {
  topLevelCount: number = 0;
  topLevelValue: string = "";
}
```

**`factory.ts`** - Object construction (when `generate_factories=true`):
```typescript
export class UtilsV1Factory {
  newNestedUtilType = (parent?: any, attributeName?: string, attributeKey?: string | number, data?: any): FactoryResult<NestedUtilTypeInterface> => {
    const instance = new ConcreteNestedUtilType();
    return { instance, fullyLoaded: false };
  }
}
```

**`schemas.ts`** - Field metadata for runtime processing:
```typescript
export const NestedUtilTypeSchema: MessageSchema = {
  name: "NestedUtilType",
  fields: [
    { name: "topLevelCount", type: FieldType.INT32, id: 1 },
    { name: "topLevelValue", type: FieldType.STRING, id: 2 },
  ]
};
```

**`deserializer.ts`** - Schema-driven deserialization:
```typescript
export class UtilsV1Deserializer extends BaseDeserializer {
  static from<T>(messageType: string, data: any): T {
    const deserializer = new UtilsV1Deserializer();
    return deserializer.createAndDeserialize<T>(messageType, data);
  }
}
```

### Usage Patterns

**Working with Interfaces (Type-safe but flexible):**
```typescript
import { NestedUtilType } from './utils/v1/interfaces';

const data: NestedUtilType = {
  topLevelCount: 1,
  topLevelValue: 'hello'
};
```

**Using Deserializer for Proper Defaults:**
```typescript
import { UtilsV1Deserializer } from './utils/v1/deserializer';

// Creates object with proper defaults for missing fields
const obj = UtilsV1Deserializer.from<NestedUtilType>(
  "utils.v1.NestedUtilType",
  { topLevelCount: 1 }  // topLevelValue will be "" (default)
);
```

**Using Model Classes:**
```typescript
import { NestedUtilType as ConcreteNestedUtilType } from './utils/v1/models';

const obj = new ConcreteNestedUtilType();
// obj.topLevelCount is 0 (default)
// obj.topLevelValue is "" (default)
```

### Architecture Benefits

- **Interfaces** - Lightweight type definitions, zero runtime cost
- **Models** - Concrete classes when you need instantiation
- **Factories** - Handles complex object graphs with proper defaults
- **Schemas** - Runtime type introspection for advanced use cases
- **Deserializers** - Schema-aware data population with type safety

This model eliminates the need for manual default handling and provides proper protobuf semantics in TypeScript.

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

### Core Generation

| Option | Description | Default |
|--------|-------------|---------|
| `wasm_export_path` | Where to generate WASM wrapper (Go generator) | `.` |
| `ts_export_path` | Where to generate TypeScript files (TS generator) | `.` |
| `module_name` | WASM module name | `{package}_services` |
| `generate_build_script` | Generate build.sh script (Go generator) | `true` |
| `generate_clients` | Generate TypeScript clients (TS generator) | `true` |
| `generate_types` | Generate TypeScript interfaces/models (TS generator) | `true` |
| `generate_factories` | Generate TypeScript factories (TS generator) | `true` |

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

### Build Integration

| Option | Description | Applies To |
|--------|-------------|------------|
| `wasm_package_suffix` | Package suffix for WASM wrapper | Go generator |
| `generate_build_script` | Generate build.sh script | Go generator |

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

**Frontend Code** (Fully Typed):
```typescript
// Import the generated per-service bundle
import { ExampleBundle } from './generated/presenter/v1/presenterServiceClient';
import type { LoadUserRequest, CallbackDemoRequest } from './generated/presenter/v1/interfaces';

// Create and load WASM bundle
const bundle = new ExampleBundle();
await bundle.loadWasm('./example.wasm');

// Fully typed method calls with IntelliSense support
const loadRequest: LoadUserRequest = { 
  userId: "user123" 
};
const userData = await presenterService.loadUserData(loadRequest);

// Async method with typed callback
const demoRequest: CallbackDemoRequest = {
  demoName: 'User Input Collection'
};
await presenterService.runCallbackDemo(demoRequest, (response, error) => {
  if (error) {
    console.error('Demo failed:', error);
    return;
  }
  console.log('Demo completed:', response.completed);
  console.log('Collected inputs:', response.collectedInputs.join(', '));
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

- **Smaller bundles**: Shared utilities eliminate code duplication
- **Better maintenance**: Runtime fixes benefit all projects immediately
- **Tree-shakeable**: Import only the utilities you need
- **Type safety**: Full TypeScript support with complete definitions
- **Inheritance-based**: Clean architecture with base class functionality

### **Usage**

```typescript
// Generated per-service bundles automatically use WASMBundle and ServiceClient
import { ExampleBundle } from './generated/presenter/v1/presenterServiceClient';

// Manual usage (advanced scenarios)
import { 
  WASMBundle, 
  ServiceClient,
  BrowserServiceManager,
  WasmError 
} from '@protoc-gen-go-wasmjs/runtime';
```

## Build Process

Generated files include a build script:

```bash
# Generated build.sh
#!/bin/bash
export GOOS=js GOARCH=wasm
go build -o example.wasm example.wasm.go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

Integration in web applications:

```html
<script src="wasm_exec.js"></script>
<script type="module">
  import { ExampleBundle } from './generated/presenter/v1/presenterServiceClient.js';
  
  // Initialize and load WASM
  const bundle = new ExampleBundle();
  await bundle.loadWasm('./example.wasm');
  
  // Use with full type safety (in TypeScript)
  const userData = await presenterService.loadUserData({
    userId: "user123"
  });
  
  console.log('User loaded:', userData.username);
</script>
```

## Self-Contained TypeScript Generation

The TypeScript generator produces complete, self-contained TypeScript code without requiring external TypeScript protobuf generators:

**What Gets Generated:**

For each proto package, the generator creates:
- **interfaces.ts** - Pure TypeScript type definitions
- **models.ts** - Concrete class implementations with defaults
- **schemas.ts** - Field metadata for runtime introspection
- **deserializer.ts** - Schema-driven data population
- **factory.ts** - Object construction (when `generate_factories=true`)
- **{service}Client.ts** - Per-service typed clients

**Benefits:**
- **No external dependencies** - Self-contained generation
- **Perfect Go protojson compatibility** - Direct JSON serialization
- **Full type safety** - Complete TypeScript types with IntelliSense
- **Cross-package imports** - Automatic import resolution
- **Well-known types** - Built-in support for google.protobuf.*

## Installation & Usage Notes

### Using the Local Generators

The plugin uses dedicated generators for Go and TypeScript:

1. **Install both generators:**
   ```bash
   go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-go@latest
   go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-ts@latest
   ```

2. **Verify installation:**
   ```bash
   which protoc-gen-go-wasmjs-go
   which protoc-gen-go-wasmjs-ts
   # Both should be in your PATH
   ```

3. **Use `local:` in your buf.gen.yaml:**
   ```yaml
   plugins:
     - local: protoc-gen-go-wasmjs-go
       out: ./gen/wasm/go
     - local: protoc-gen-go-wasmjs-ts
       out: ./web/src/generated
   ```

### Using Remote Plugins (buf.build)

When published to buf.build, you can use the generators without local installation:

```yaml
plugins:
  - remote: buf.build/panyam/protoc-gen-go-wasmjs-go
    out: ./gen/wasm/go
  - remote: buf.build/panyam/protoc-gen-go-wasmjs-ts
    out: ./web/src/generated
```

**Benefits of remote plugins:**
- No local installation required
- Always uses the latest version
- Consistent across team members
- Works in CI/CD without additional setup

## Project Structure

```
├── cmd/
│   ├── protoc-gen-go-wasmjs-go/     # Go WASM generator
│   └── protoc-gen-go-wasmjs-ts/     # TypeScript client generator
├── pkg/
│   ├── core/                        # Pure utility functions (30+ tests)
│   ├── filters/                     # Business logic filtering (25+ tests)
│   ├── builders/                    # Template data building
│   ├── renderers/                   # Template rendering with typed imports
│   ├── generators/                  # Top-level orchestrators
│   └── wasm/                        # WASM runtime utilities
├── runtime/                         # @protoc-gen-go-wasmjs/runtime NPM package
│   ├── src/client/                  # WASMServiceClient base class
│   ├── src/browser/                 # BrowserServiceManager
│   └── src/schema/                  # Type utilities
├── proto/wasmjs/v1/                 # WASM annotation definitions
├── example/                         # Complete demo with per-service clients, typed callbacks, and browser services
└── docs/                            # Architecture and development guides
```

## Development

For detailed development instructions, testing guidelines, and contribution workflows, see [DEVELOPMENT.md](DEVELOPMENT.md).

### Quick Start

```bash
# Run the test suite
./test.sh

# Build the generators  
make

# Test with examples
cd example && make buf && make wasm

# Run framework tests
go test ./pkg/... -v
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
