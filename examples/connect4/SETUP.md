# Connect4 Example Setup

This example demonstrates real-time multiplayer Connect4 using protoc-gen-go-wasmjs with pluggable transport mechanisms.

## Prerequisites

- **Go 1.23+** - For WASM compilation and server
- **Node.js 18+** - For frontend build system
- **pnpm** - Package manager (`npm install -g pnpm`)
- **Buf CLI** - Protocol buffer toolchain (`brew install bufbuild/buf/buf`)

## Setup Modes

### Production Mode (Default)
Uses published wasmjs proto dependencies from buf.build:

```bash
# Build parent plugin first (if needed)
cd ../..
make

# Build Connect4 example
cd examples/connect4
make all
```

### Development Mode
For developing changes to the wasmjs annotations:

```bash
# Generate with local wasmjs protos (creates symlink automatically)
make generate-dev
```

## Running the Game

```bash
# Build everything and start the server (recommended)
make all
# ↳ Runs: make parent → make generate → make wasm → make web

# Or step by step:
make generate  # Generate protobuf + WASM bindings
make wasm      # Build WASM binary (multiplayer_connect4.wasm)
make web       # Build TypeScript frontend + start server
```

**Open http://localhost:8080 to play Connect4!**

## Configuration Files

**Production Mode:**
- `buf.yaml`: References `buf.build/panyam/protoc-gen-go-wasmjs` for wasmjs proto dependencies
- `buf.gen.yaml`: Uses `local: ../../bin/protoc-gen-go-wasmjs` for plugin execution
- `web/package.json`: Frontend dependencies (TypeScript, Webpack, pnpm)

**Development Mode:**
- `buf.gen.dev.yaml`: Uses local plugin and local wasmjs protos via symlink
- `buf.dev.yaml`: Excludes buf.build wasmjs dependency to avoid conflicts

## Build Outputs

After running `make all`:
```
web/static/wasm/multiplayer_connect4.wasm      # Compiled Go service  
web/static/gen/js/index.js                     # TypeScript→JS (games list)
web/static/gen/js/gameViewer.js                # TypeScript→JS (game UI)
web/gen/wasmts/multiplayer_connect4Client.client.ts  # Generated WASM client
cmd/serve/templates/gen/*.html                  # Generated HTML templates
```

## Frontend Development

```bash
cd web
pnpm install                 # Install dependencies
npm run dev                  # Webpack dev server (hot reload)
npm run build                # Production build
```

## Features Demonstrated

- **Real-time multiplayer**: Cross-tab via IndexedDB + polling
- **WASM integration**: Go service compiled to WebAssembly
- **TypeScript generation**: Fully-typed WASM client bindings  
- **Pluggable transports**: IndexedDB, BroadcastChannel, WebSocket-ready
- **State persistence**: Games survive browser restarts
- **Independent module**: Self-contained with own go.mod