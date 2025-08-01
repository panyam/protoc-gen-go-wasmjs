# Next Steps

## Recently Completed (August 2025)
- ✅ **Proto to JSON Conversion System**: Implemented flexible conversion options to handle differences between Go protojson and TypeScript protobuf libraries
  - Added `ConversionOptions` interface with `handleOneofs`, `fieldTransformer`, `emitDefaults`, and `bigIntHandler`
  - Enhanced error handling with better context
  - Added runtime configuration via `setConversionOptions()`
  - Updated WASM-side with better protojson marshal/unmarshal options
  - Created comprehensive documentation in `PROTO_CONVERSION.md`

## Immediate Tasks
- [ ] **Browser Demo**: Create a complete browser demo showcasing the proto conversion features
  - Demonstrate oneof field handling
  - Show field name transformation in action
  - Include BigInt serialization examples
  
- [ ] **Performance Testing**: Benchmark the conversion overhead
  - Measure impact of custom conversions
  - Optimize hot paths in the conversion system
  - Add performance monitoring hooks

- [ ] **Type-Aware Conversion**: Enhance conversion system with proto type information
  - Detect oneof fields from proto descriptors instead of heuristics
  - Better BigInt field detection based on proto field types
  - Support for Well-Known Types (Timestamp, Duration, etc.)

## Short Term (Next Month)
- [ ] **Streaming Support**: Research and implement streaming RPC support for WASM
  - Server-streaming RPCs
  - Client-streaming RPCs
  - Bidirectional streaming
  
- [ ] **Advanced Templates**: Template inheritance and partial overrides
  - Allow extending base templates
  - Support for custom template functions
  - Template composition for complex scenarios

- [ ] **Error Recovery**: Implement graceful error recovery
  - Retry mechanisms for transient failures
  - Circuit breaker pattern
  - Fallback to HTTP when WASM fails

## Medium Term (Next Quarter)
- [ ] **Monitoring Integration**: Performance monitoring and analytics
  - OpenTelemetry integration
  - Custom metrics for WASM performance
  - Distributed tracing support
  
- [ ] **IDE Support**: Language server and IDE plugin support
  - VS Code extension for proto to WASM development
  - IntelliJ IDEA plugin
  - Syntax highlighting for WASM annotations

- [ ] **Advanced Conversion Features**:
  - Custom type converters via configuration
  - Conversion middleware system
  - Proto extension support

## Long Term (6+ Months)
- [ ] **Multi-Language Support**: 
  - Rust service implementations with same TypeScript clients
  - C++ WASM generation
  - Shared TypeScript client generation

- [ ] **Edge Computing Optimization**:
  - CDN-friendly WASM deployment
  - Edge worker compatibility
  - Lazy loading strategies

- [ ] **Enterprise Features**:
  - Advanced authentication patterns
  - Authorization middleware
  - Audit logging
  - Compliance features (GDPR, HIPAA)

## Community & Ecosystem
- [ ] **Documentation Improvements**:
  - Video tutorials
  - Example repository with common patterns
  - Migration guides from other solutions
  
- [ ] **Framework Integration**:
  - React hooks for WASM services
  - Vue composables
  - Angular services
  - Svelte stores

- [ ] **Testing Tools**:
  - Mock generation for WASM services
  - Integration testing framework
  - Performance testing suite

## Technical Debt
- [ ] **Code Cleanup**:
  - Refactor template generation for better maintainability
  - Add more unit tests for conversion logic
  - Improve error messages throughout

- [ ] **Performance Optimizations**:
  - Lazy WASM loading
  - Connection pooling for service calls
  - Caching strategies

## Research Topics
- [ ] **WebAssembly Component Model**: Investigate integration with WASM Component Model
- [ ] **SharedArrayBuffer**: Explore using SharedArrayBuffer for better performance
- [ ] **WASM SIMD**: Leverage WASM SIMD for data processing
- [ ] **Module Federation**: Integration with Webpack Module Federation