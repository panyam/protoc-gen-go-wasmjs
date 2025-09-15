# Default: production mode using published wasmjs protos + local plugin  
buf: 
	rm -Rf gen
	buf generate
	goimports -w `find gen | grep "\.go"`

# Development mode: uses local plugin and symlinked wasmjs annotations
buf-dev: parent symlink-wasmjs
	rm -Rf gen
	# Update dependencies for dev config to avoid conflicts
	rm -f buf.lock
	buf generate --template buf.gen.dev.yaml --config buf.dev.yaml
	goimports -w `find gen | grep "\.go"`

clean:
	rm -Rf gen

parent:
	cd ../../ && make

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

test: buf
	cd gen/wasm/gen/wasm && sh build.sh

# Build WASM module from generated code
build-wasm:
	@echo "Building WASM module..."
	GOOS=js GOARCH=wasm go build -o gen/wasm/browser_example.wasm .
	@echo "Copying wasm_exec.js..."
	@cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/
	@echo "Transpiling TypeScript files to JavaScript..."
	@esbuild gen/wasm/browserServiceManager.ts --outfile=gen/wasm/browserServiceManager.js --format=esm --target=es2020
	@esbuild gen/wasm/browser_exampleClient.client.ts --outfile=gen/wasm/browser_exampleClient.client.js --format=esm --target=es2020

# Serve the example with Python HTTP server
serve: build-wasm
	@echo "Starting web server at http://localhost:8080/web/"
	@echo "Open http://localhost:8080/web/ in your browser"
	@python3 -m http.server 8080

# Alternative: Serve with a simple Go server (if Python is not available)
serve-go: build-wasm
	@echo "Starting Go web server at http://localhost:8080/web/"
	@echo "Open http://localhost:8080/web/ in your browser"
	@go run -e 'package main; import "net/http"; import "log"; func main() { log.Println("Server starting on :8080..."); log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))) }'

# Run the complete example
run: buf-dev build-wasm serve

# Run with production config
run-prod: buf build-wasm serve

.PHONY: build-wasm serve serve-go run run-prod
