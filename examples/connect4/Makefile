# Connect4 Example Makefile
# Demonstrates stateful proxy generation for real-time multiplayer games

.PHONY: all generate wasm web clean test

all: parent generate wasm web

parent:
	cd ../../ && make

# Generate all protobuf code using buf
generate:
	echo buf generate

# Build WASM service
wasm: generate
	cd wasm && GOOS=wasip1 GOARCH=wasm go build -o ../web/connect4_service.wasm service.go
	cd cmd/wasm && GOOS=wasip1 GOARCH=wasm go build -o ../../web/multiplayer_connect4.wasm *.go

# Start development web server
web:
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
