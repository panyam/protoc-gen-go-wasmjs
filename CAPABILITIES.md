# protoc-gen-go-wasmjs

## Version
0.0.33

## Provides
- wasm-bindings: Go-to-WASM binding generation from protobuf services
- typescript-client-gen: TypeScript client generation with full type safety
- four-step-pipeline: Artifact processing (collect, filter, decorate, generate)
- composition-bundles: Bundle classes with dependency injection for services
- dual-target: Separate Go WASM and TypeScript generators
- local-first: Same service logic runs in server or browser via WASM
- runtime-package: @protoc-gen-go-wasmjs/runtime for browser integration
- invocation-styles: SYNC/PROMISE/CALLBACK control for WASM-JS threading
- streaming-support: Server streaming from WASM to JavaScript via callbacks
- cross-package-imports: Cross-package TypeScript import resolution from proto descriptors
- schema-deserialization: Schema-aware deserialization with factory composition

- **@protoc-gen-go-wasmjs/runtime** (@protoc-gen-go-wasmjs/runtime): Browser runtime with base classes for generated TS clients
  - location: runtime/
  - module: @protoc-gen-go-wasmjs/runtime

## Module
github.com/panyam/protoc-gen-go-wasmjs

## Location
newstack/protoc-gen-go-wasmjs/main

## Stack Dependencies
- protokit (github.com/panyam/protokit)

## Integration

### Buf Plugins
Add as buf plugins in buf.gen.yaml:
```yaml
plugins:
  - local: protoc-gen-go-wasmjs-go
    out: ./gen/wasm
  - local: protoc-gen-go-wasmjs-ts
    out: ./frontend/src/generated
```

### Go Module
```go
// go.mod
require github.com/panyam/protoc-gen-go-wasmjs v0.0.33

// Local development
replace github.com/panyam/protoc-gen-go-wasmjs => ~/newstack/protoc-gen-go-wasmjs/main
```

### Key Imports
```go
import "github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm"
```

### Runtime Package (npm)
```bash
npm install @protoc-gen-go-wasmjs/runtime
```

## Status
Active

## Conventions
- Split generator architecture (BaseGenerator + TSGenerator + GoGenerator)
- Composition-based bundles
- Artifact pipeline (collect → classify → map → create)
- InvocationStyle annotations for async patterns

## Migrations
