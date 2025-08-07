# Library Example Setup

This example demonstrates how to use protoc-gen-go-wasmjs with published wasmjs proto dependencies.

## Setup Modes

### Production Mode (Default)
Uses published wasmjs proto dependencies from buf.build + locally installed plugin:

```bash
# Install the plugin locally
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest

# Generate using published wasmjs proto dependencies
make buf
```

**Configuration:**
- `buf.yaml`: References `buf.build/panyam/protoc-gen-go-wasmjs` for wasmjs proto dependencies
- `buf.gen.yaml`: Uses `local: ../../bin/protoc-gen-go-wasmjs` for plugin execution

### Development Mode
For developing changes to the wasmjs annotations themselves:

```bash
# Generate with local wasmjs protos (requires symlink)
make buf-dev
```

**Configuration:**
- `buf.gen.dev.yaml`: Uses local plugin
- Creates symlink to local wasmjs protos for development
- Automatically removes symlink when switching back to production mode

## Benefits of This Approach

1. **No BSR Plugin Limitations**: Uses local plugin installation, avoiding buf.build plugin registry limitations
2. **Published Proto Dependencies**: Users don't need to manage wasmjs proto symlinks
3. **Standard Installation**: Users can `go install` the plugin normally
4. **Development Flexibility**: Plugin developers can still work with local modifications

## For Plugin Users

Simply add the wasmjs proto dependency to your `buf.yaml`:

```yaml
version: v2
modules:
  - path: proto
deps:
  - buf.build/googleapis/googleapis
  - buf.build/protocolbuffers/wellknowntypes
  - buf.build/panyam/protoc-gen-go-wasmjs  # Add this line
```

And install the plugin:

```bash
go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs@latest
```

Then use `local: protoc-gen-go-wasmjs` in your `buf.gen.yaml`.