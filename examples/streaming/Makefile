# Streaming Example Makefile
# Demonstrates server streaming capabilities with WASM

.PHONY: all generate wasm web clean test 

# Run generate separately for now
all: parent wasm web

parent:
	cd ../../ && make

# Default: production mode using published wasmjs protos + local plugin
generate:
	rm -Rf gen
	buf generate

# Development mode: uses local plugin and symlinked wasmjs annotations
generate-dev: parent symlink-wasmjs
	rm -Rf gen
	# Remove buf.lock to avoid conflicts with published wasmjs dependency
	rm -f buf.lock
	buf generate --template buf.gen.dev.yaml --config buf.dev.yaml

# Create symlink to wasmjs annotations for development
symlink-wasmjs:
	@if [ ! -L proto/wasmjs ]; then \
		echo "Creating wasmjs symlink for development..."; \
		ln -s ../../../proto/wasmjs proto/wasmjs; \
	fi

# Remove symlink (for switching back to production mode)
remove-symlink:
	@if [ -L proto/wasmjs ]; then \
		echo "Removing wasmjs symlink..."; \
		rm proto/wasmjs; \
	fi

# Build WASM service
wasm: 
	mkdir -p web/static/wasm
	GOOS=js GOARCH=wasm go build -o ./web/static/wasm/streaming_service.wasm cmd/wasm/*.go

# Start development web server
web: 
	cd web && make
	@echo "Starting Streaming demo server..."
	@echo "Open http://localhost:8080 to test streaming!"
	go run cmd/serve/*.go

# Clean generated files
clean:
	rm -rf gen/
	rm -rf web/gen/
	rm -rf static/gen/
	rm -Rf web/static/gen
	rm -f web/*.wasm

# Test the streaming generation
test: generate-dev
	@echo "Testing server streaming generation..."
	@echo "✅ Server streaming generation successful!"
	@echo "Generated files:"
	@find gen/ -name "*.go" -o -name "*.ts" | head -10