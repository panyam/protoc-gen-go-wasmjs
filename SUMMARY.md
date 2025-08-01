# protoc-gen-go-wasmjs Summary

## Overview
protoc-gen-go-wasmjs is a Protocol Buffers compiler plugin that generates WASM bindings and TypeScript clients for gRPC services. It enables local-first applications where the same service logic can run in both server environments (with full database access) and browser environments (with local storage).

## Core Features

### 1. **Dual-Target Architecture**
- Generate WASM wrappers and TypeScript clients separately or together
- Flexible deployment with TypeScript clients directly in frontend directories
- Clean separation of concerns with independent generation targets

### 2. **Smart Type System Integration**
- Auto-detects TypeScript generator (protoc-gen-es, protoc-gen-ts)
- Smart import detection analyzes proto files for accurate TypeScript imports
- Automatic extension detection (.ts vs .js) based on generator configuration

### 3. **Flexible Proto to JSON Conversion**
- Configurable conversion options to handle Go protojson vs TypeScript differences
- Special handling for oneof fields with flattening support
- Field name transformation (camelCase ↔ snake_case)
- BigInt serialization/deserialization
- Default value management

### 4. **Multi-Service & Multi-Target Support**
- Bundle multiple services in a single WASM module
- Generate optimized bundles per page/use case
- Service filtering for targeted deployments
- Dependency injection through export pattern

### 5. **Developer Experience**
- Seamless buf.build integration
- Generated build scripts and examples
- Comprehensive error handling
- Extensive customization options

## Project Structure

```
├── cmd/protoc-gen-go-wasmjs/     # Plugin entry point
├── pkg/generator/                # Core generation logic
│   ├── templates/                # Embedded templates
│   │   ├── wasm.go.tmpl         # Go WASM wrapper
│   │   ├── client.ts.tmpl       # TypeScript client
│   │   └── build.sh.tmpl        # Build script
│   ├── config.go                # Configuration parsing
│   ├── generator.go             # Main generation logic
│   └── types.go                 # Template data structures
├── proto/wasmjs/v1/             # WASM annotations
├── example/                     # Complete example
└── PROTO_CONVERSION.md          # Conversion documentation
```

## Key Design Decisions

### 1. **Export Pattern over Main Functions**
Instead of generating `main()` functions, we generate export structs that allow full dependency injection:
```go
exports := &ServicesExports{
    UserService: &myUserService{db: database, auth: authService}
}
exports.RegisterAPI()
```

### 2. **Template-Based Generation**
Using Go's `embed` package for templates provides:
- Clean separation of generation logic and output format
- Easy customization through template overrides
- Maintainable code generation

### 3. **Proto-Aware Import Generation**
Analyzes proto file sources to generate accurate imports instead of hardcoded assumptions:
```typescript
// Imports from actual proto sources
import { CreateGameRequest } from './games_pb';
import { CreateUserRequest } from './users_pb';
```

### 4. **Flexible API Structures**
Supports multiple JavaScript API patterns:
- Namespaced: `myapp.service.method()`
- Flat: `myappServiceMethod()`
- Service-based: `services.library.findBooks()`

## Technical Achievements

- **Zero Configuration**: Works out of the box with smart defaults
- **Backward Compatible**: Existing configurations continue to work
- **Production Ready**: Used in real projects with complex proto structures
- **Performance Optimized**: Minimal overhead in proto conversions
- **Type Safe**: Full TypeScript support with proper type inference

## Use Cases

### Local-First Applications
Same business logic runs in both environments:
- **Server**: Full database access, complete dataset
- **Browser**: Local storage, user's subset of data

### Progressive Web Apps
- Offline functionality with WASM services
- Sync when online, work offline
- Reduced server load

### Edge Computing
- Deploy services to CDN edge workers
- Run logic closer to users
- Maintain consistent API across deployments

## Current Status
The project is feature-complete for the core use cases with:
- Comprehensive proto to JSON conversion system
- Multi-target generation support
- Smart import detection
- Full buf.build integration
- Production-ready code generation

Recent additions include enhanced proto conversion options for better handling of differences between Go and TypeScript protobuf implementations, particularly for oneof fields and field naming conventions.