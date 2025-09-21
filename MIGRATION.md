# Migration Guide: Bundle-Based Architecture

This guide helps you migrate from the old per-service WASM loading to the new bundle-based architecture.

## What Changed

### Old Architecture (Problematic)
```typescript
// Each service loaded WASM independently
const presenterClient = new Presenter_v1Client();
const browserClient = new Browser_v1Client();

await presenterClient.loadWasm('/module.wasm'); // Load 1
await browserClient.loadWasm('/module.wasm');   // Load 2 (duplicate!)

await presenterClient.presenterService.someMethod();
await browserClient.browserService.someMethod();
```

**Problems:**
- ❌ Duplicate WASM loading for services in same module
- ❌ Resource waste and slower initialization
- ❌ Each service client managed its own WASM state

### New Architecture (Fixed)
```typescript
// Bundle manages WASM loading, services share it
import { Presenter_v1Bundle } from './generated/presenter/v1/presenterServiceClient';

const bundle = new Presenter_v1Bundle();
await bundle.loadWasm('/module.wasm'); // Single load for all services

// All services in the bundle share the same WASM instance
await bundle.presenterService.someMethod();
await bundle.browserService.someMethod();
```

**Benefits:**
- ✅ Single WASM load per module
- ✅ Efficient resource usage
- ✅ Shared state across services
- ✅ Clean separation of concerns

## Migration Steps

### Step 1: Update Imports
```diff
- import { Presenter_v1Client } from './generated/presenter/v1/presenterServiceClient';
+ import { Presenter_v1Bundle } from './generated/presenter/v1/presenterServiceClient';
```

### Step 2: Change Client Creation
```diff
- const client = new Presenter_v1Client();
+ const bundle = new Presenter_v1Bundle();
```

### Step 3: Update WASM Loading
```diff
- await client.loadWasm('/module.wasm');
+ await bundle.loadWasm('/module.wasm');
```

### Step 4: Update Service Calls
```diff
- await client.presenterService.someMethod(request);
+ await bundle.presenterService.someMethod(request);
```

### Step 5: Update Browser Service Registration
```diff
- client.registerBrowserService('BrowserAPI', implementation);
+ bundle.registerBrowserService('BrowserAPI', implementation);
```

## Complete Example Migration

### Before
```typescript
import { Presenter_v1Client } from './generated/presenter/v1/presenterServiceClient';
import { BrowserAPIImpl } from './browser-api-impl';

async function init() {
  const presenterClient = new Presenter_v1Client();
  
  presenterClient.registerBrowserService('BrowserAPI', new BrowserAPIImpl());
  await presenterClient.loadWasm('/browser_example.wasm');
  
  const response = await presenterClient.presenterService.loadUserData({
    userId: 'user123'
  });
}
```

### After
```typescript
import { Presenter_v1Bundle } from './generated/presenter/v1/presenterServiceClient';
import { BrowserAPIImpl } from './browser-api-impl';

async function init() {
  const bundle = new Presenter_v1Bundle();
  
  bundle.registerBrowserService('BrowserAPI', new BrowserAPIImpl());
  await bundle.loadWasm('/browser_example.wasm');
  
  const response = await bundle.presenterService.loadUserData({
    userId: 'user123'
  });
}
```

## Bundle Naming

Bundle names are generated from the module name configured in your `buf.gen.yaml`:

```yaml
# This creates: MyServices_Bundle
- module_name=my_services

# This creates: BrowserCallbacks_Bundle  
- module_name=browser_callbacks
```

The bundle contains all services from protobuf packages that compile to the same WASM module.

## API Compatibility

All service method signatures remain the same - only the client creation and WASM loading changes:

```typescript
// Method signatures are unchanged
await bundle.presenterService.loadUserData(request);      // ✅ Same
await bundle.presenterService.updateUIState(request, cb); // ✅ Same
await bundle.presenterService.savePreferences(request);   // ✅ Same
```

## Regenerating Code

After updating protoc-gen-go-wasmjs, regenerate your clients:

```bash
# Clean and regenerate
buf generate

# Or with make
make clean && make buf
```

The new bundle-based clients will be generated automatically.
