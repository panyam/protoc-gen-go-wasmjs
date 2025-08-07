# Connect4 Example Makefile
# Demonstrates stateful proxy generation for real-time multiplayer games

.PHONY: all generate wasm web clean test 

# Run generate seperately for now
all: parent web wasm web

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
	mkdir -p wasm
	GOOS=js GOARCH=wasm go build -o ./static/wasm/multiplayer_connect4.wasm cmd/wasm/*.go

# Start development web server
web: 
	cd web && make
	@echo "Starting Connect4 game server..."
	@echo "Open http://localhost:8080 to play!"
	go run cmd/serve/*.go

# Generate only stateful proxies (for development)
stateful:
	../../bin/protoc-gen-go-wasmjs-stateful \
		-I proto \
		-I ../../proto \
		--go-wasmjs-stateful_out=./web/gen \
		proto/connect4/game.proto

# Clean generated files
clean:
	rm -rf gen/
	rm -rf web/gen/
	rm -rf static/gen/
	rm -f web/*.wasm

# Test the stateful proxy generation
test: ../../bin/protoc-gen-go-wasmjs-stateful
	@echo "Testing stateful proxy generation..."
	protoc \
		--plugin=protoc-gen-go-wasmjs-stateful=../../bin/protoc-gen-go-wasmjs-stateful \
		--go-wasmjs-stateful_out=./web/gen \
		-I proto \
		-I ../../proto \
		proto/connect4/game.proto
	@echo "✅ Stateful proxy generation successful!"

# Build dependencies
../../bin/protoc-gen-go-wasmjs-stateful:
	cd ../.. && make stateful
