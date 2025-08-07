# Connect4 Example Makefile
# Demonstrates stateful proxy generation for real-time multiplayer games

.PHONY: all generate wasm web clean test compile-ts

all: parent generate web compile-ts wasm web

parent:
	cd ../../ && make

# Generate all protobuf code using buf
generate:
	buf generate

# Build WASM service
wasm: compile-ts
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
