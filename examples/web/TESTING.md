# Testing Guide for Browser Callbacks Example

## Overview

This example includes comprehensive stability tests to ensure that generated TypeScript code remains stable across generator improvements.

## Test Files

### `src/generation-stability-tests.ts`

Comprehensive import and functionality tests that verify:

- **Service Clients**: All generated service clients can be instantiated
- **Interfaces**: TypeScript type definitions are correctly exported
- **Models**: Concrete class implementations work with default values
- **Schemas**: Schema registries are accessible and properly structured
- **Factory/Deserializer**: Annotation-based factory generation works correctly
- **Cross-package imports**: Types from different packages import correctly
- **Type compatibility**: Interfaces are assignable from concrete models

### `src/smoke-test.ts`

Quick smoke test that runs all stability tests. This file can be imported to verify all generated artifacts compile and work correctly.

## Running Tests

### Type Checking (Recommended)

The most reliable way to test generated code stability:

```bash
pnpm test
```

This runs TypeScript's type checker across all files. If the generated code has breaking changes, compilation will fail with clear error messages.

### Stability Tests Only

To test only the stability test files:

```bash
pnpm test:stability
```

### Build Test

Building the entire project also validates all imports:

```bash
pnpm build
```

## What Gets Tested

The stability tests import and verify ALL generated artifacts:

### Service Clients
- `PresenterServiceClient` - Main WASM service
- `BrowserAPIClient` - Browser-provided service
- `TestServiceClient` (test_one_package.v1)
- `TestServiceClient` (test_multi_packages.v1)

### Message Types
- **Presenter package**: LoadUserDataRequest/Response, StateUpdateRequest, UIUpdate, PreferencesRequest/Response, etc.
- **Browser API package**: FetchRequest/Response, StorageKeyRequest, CookieRequest, AlertRequest, PromptRequest, LogRequest
- **Test packages**: SampleRequest/Response, SecondRequest/Response
- **Utils package**: HelperUtilType, ParentUtilMessage, nested types
- **WASMJS annotations**: StatefulOptions, MessagePatch, PatchBatch, etc.

### Type Artifacts
- **Interfaces** (`interfaces.ts`) - Pure TypeScript type definitions
- **Models** (`models.ts`) - Concrete class implementations with defaults
- **Schemas** (`schemas.ts`) - Runtime field metadata for introspection
- **Factory** (`factory.ts`) - Annotation-based factory + deserializer
- **Service Clients** (`*ServiceClient.ts`) - WASM and browser service proxies

## Why This Matters

These tests ensure that improvements to the code generator don't introduce breaking changes:

1. **Compilation guarantees**: If tests compile, all imports are correct
2. **Type safety**: Catches missing exports, renamed types, wrong paths
3. **Regression detection**: Immediately shows if a change breaks existing code
4. **Documentation**: Serves as living documentation of all generated artifacts

## Continuous Integration

These tests should be run:

- Before committing generator changes
- In CI/CD pipelines after regeneration
- When updating dependencies
- Before releases

## Expected Output

Successful test run produces:
```
✅ All service clients instantiated successfully
✅ All models instantiated successfully with defaults
✅ All schema registries accessible
✅ Factory and deserializer working correctly
✅ Interface and model types are compatible
✅ Cross-package imports working correctly

✅ All generation stability tests passed!
```

## Troubleshooting

### Import Errors

If you see `Module has no exported member 'X'`:
1. Check the actual exports in the generated file
2. Update the import in `generation-stability-tests.ts`
3. This indicates the generator changed export names - document the breaking change

### Type Incompatibility

If you see type assignment errors:
1. Check if interfaces and models are still compatible
2. This may indicate a serious breaking change in the type system
3. Review the generator changes carefully

### Missing Files

If generated files are missing:
1. Run `buf generate` to regenerate
2. Check if the file generation logic changed
3. Update the test to reflect the new file structure
