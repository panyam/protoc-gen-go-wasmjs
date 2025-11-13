# Browser Callbacks Example

This example demonstrates how WASM services can call back into browser-provided APIs using the `browser_provided` annotation.

## Architecture

```
┌─────────────────────────────────────────────┐
│                Browser (JS)                 │
│                                              │
│  ┌─────────────────────────────────────┐    │
│  │    BrowserAPI Implementation        │    │
│  │  - fetch()                          │◄───┼──── Calls from WASM
│  │  - localStorage access              │    │     (blocking experience)
│  │  - cookie access                    │    │
│  │  - alert()                          │    │
│  └─────────────────────────────────────┘    │
│                     ▲                        │
│                     │ FIFO Queue             │
│                     ▼                        │
│  ┌─────────────────────────────────────┐    │
│  │    Browser Service Manager          │    │
│  │  - Processes calls sequentially     │    │
│  │  - Handles timeouts                 │    │
│  │  - Routes to implementations        │    │
│  └─────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
                      ▲
                      │ Browser Channel
                      ▼
┌─────────────────────────────────────────────┐
│               WASM Module                    │
│                                              │
│  ┌─────────────────────────────────────┐    │
│  │     PresenterService (MVP)          │    │
│  │  - Drives UI logic                  │────┼──── Streams UI updates
│  │  - Fetches data via BrowserAPI      │    │     to browser
│  │  - Manages state                    │    │
│  └─────────────────────────────────────┘    │
│                     │                        │
│                     ▼                        │
│  ┌─────────────────────────────────────┐    │
│  │     BrowserAPIClient                │    │
│  │  - Blocking API for WASM code       │    │
│  │  - Queues calls through channel     │    │
│  │  - Handles timeouts                 │    │
│  └─────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

## Key Features

1. **Browser-Provided Services**: Services marked with `browser_provided` annotation generate clients on WASM side
2. **Blocking Experience**: WASM code uses simple blocking calls that are queued and processed FIFO
3. **MVP Pattern**: Presenter runs in WASM, driving UI logic through browser callbacks
4. **Server Streaming**: UI updates are streamed from WASM to browser
5. **Type Safety**: Full proto type safety across WASM/JS boundary

## Running the Example

### Development Mode (Local Plugin)

```bash
# Build plugin and generate code
make dev

# Or step by step:
make build-plugin    # Build the plugin
make generate-dev    # Generate with local plugin
make build-wasm      # Build WASM module
make run             # Start web server
```

### Production Mode (Installed Plugin)

```bash
# Install plugin first
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest

# Generate and run
make prod
```

Then open http://localhost:8080 in your browser.

## Files Structure

```
example/
├── proto/
│   ├── browser/v1/
│   │   └── browser.proto         # Browser API service (browser_provided)
│   └── presenter/v1/
│       └── presenter.proto        # Presenter service (runs in WASM)
├── services/
│   └── presenter.go               # Presenter implementation
├── cmd/wasm/
│   └── main.go                    # WASM entry point
├── web/
│   └── index.html                 # Test page
├── gen/                           # Generated code (git ignored)
│   ├── go/                        # Go protobuf types
│   ├── wasm/                      # WASM bindings and TS client
│   └── ts/                        # TypeScript types
├── buf.yaml.{dev,prod}            # Buf configurations
├── buf.gen.yaml.{dev,prod}        # Generation configurations
└── Makefile                       # Build automation
```

## How It Works

1. **Proto Definition**: Services marked with `browser_provided` option
2. **Code Generation**: Plugin generates:
   - WASM side: Client that calls through browser channel
   - JS side: Server interface for implementation
3. **Runtime**:
   - WASM calls are queued FIFO
   - Browser processes them sequentially
   - Responses delivered back to waiting goroutines
4. **No Deadlocks**: Single channel ensures sequential processing

## Usage

The example uses the new bundle-based architecture:

```typescript
import { ExampleBundle} from './generated/presenter/v1/presenterServiceClient';

// Create bundle - manages WASM loading for all services
const bundle = new ExampleBundle();

// Register browser API implementation
bundle.registerBrowserService('BrowserAPI', new BrowserAPIImpl());

// Load WASM once for all services in the bundle
await bundle.loadWasm('/browser_example.wasm');

// Use individual service clients (all share the same WASM)
await bundle.presenterService.loadUserData({ userId: 'user123' });
```

**Key Benefits:**
- Single WASM load per module (not per service)
- Multiple services share the same WASM instance
- Clean separation: bundle manages WASM, service clients handle business logic

## Testing

Click the buttons in the web UI to:

1. **Load User Data**: Fetches from API (mocked) and caches in localStorage
2. **Update UI State**: Streams UI updates from WASM presenter
3. **Save Preferences**: Stores preferences in localStorage with alert confirmation

Watch the console output to see the flow of calls between WASM and browser.
