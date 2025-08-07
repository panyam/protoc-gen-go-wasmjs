# Connect4 Example Setup

This example demonstrates real-time multiplayer Connect4 using protoc-gen-go-wasmjs with stateful proxy generation.

## Setup Modes

### Production Mode (Default)
Uses published wasmjs proto dependencies from buf.build + locally installed plugin:

```bash
# Install the plugin locally
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest

# Generate using published wasmjs proto dependencies
make generate
```

### Development Mode
For developing changes to the wasmjs annotations themselves:

```bash
# Generate with local wasmjs protos (creates symlink automatically)
make generate-dev
```

## Running the Game

```bash
# Build everything and start the server
make all

# Or step by step:
make generate  # Generate protobuf code
make web       # Build frontend assets
make wasm      # Build WASM service
make web       # Start the server
```

Open http://localhost:8080 to play Connect4!

## Configuration

**Production Mode:**
- `buf.yaml`: References `buf.build/panyam/protoc-gen-go-wasmjs` for wasmjs proto dependencies
- `buf.gen.yaml`: Uses `local: ../../bin/protoc-gen-go-wasmjs` for plugin execution

**Development Mode:**
- `buf.gen.dev.yaml`: Uses local plugin and local wasmjs protos via symlink
- `buf.dev.yaml`: Excludes buf.build wasmjs dependency to avoid conflicts

## Features Demonstrated

- Real-time multiplayer gameplay
- WASM service generation with dual-target architecture
- TypeScript client generation for frontend
- Stateful proxy patterns (when enabled)
- Clean separation between WASM and TypeScript artifacts