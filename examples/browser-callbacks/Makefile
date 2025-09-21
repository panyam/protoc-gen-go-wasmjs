# Default: production mode using new split generators

buf: ensureenv clean parent 
	buf generate
	goimports -w `find gen | grep "\.go"`

all: buf build-wasm serve

ensureenv:
	@test -f buf.yaml && test -f buf.gen.yaml && echo "buf.yaml does not exist.  Run 'make setupdev' or 'make setupprod' to setup your environment..."

setupdev: cleanall symlink-wasmjs
	ln -s buf.gen.yaml.dev buf.gen.yaml
	ln -s buf.yaml.dev buf.yaml

setupprod: cleanall remove-symlink
	ln -s buf.gen.yaml.prod buf.gen.yaml
	ln -s buf.yaml.prod buf.yaml

cleanall: clean
	rm -f buf.yaml
	rm -f buf.gen.yaml

clean:
	rm -Rf gen
	rm -f buf.lock

parent:
	cd ../../ && make default old

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

# Build WASM module and setup TypeScript project
build-wasm:
	@echo "Building WASM module..."
	GOOS=js GOARCH=wasm go build -o web/public/browser_example.wasm .
	@echo "Copying wasm_exec.js..."
	@cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/public/
	@echo "Installing TypeScript dependencies..."
	@cd web && pnpm install

# Serve the example with Vite dev server
serve: build-wasm
	@echo "Starting Vite dev server at http://localhost:8080/"
	@echo "TypeScript will be compiled on-the-fly"
	@cd web && pnpm run dev

# Alternative: Serve with a simple Go server (if Python is not available)
serve-go: build-wasm
	@echo "Starting Go web server at http://localhost:8080/web/"
	@echo "Open http://localhost:8080/web/ in your browser"
	@go run -e 'package main; import "net/http"; import "log"; func main() { log.Println("Server starting on :8080..."); log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))) }'

# Run the complete example - choose dev or prod setup with `make setupdev` or `make setupprod`
run: buf build-wasm serve

# TypeScript development commands
tsbuild:
	@cd web && pnpm run typecheck && pnpm run build

ts-preview:
	@cd web && pnpm run preview

# Quick development workflow
dev: setupdev buf build-wasm serve

.PHONY: build-wasm serve serve-go run setupdev setupprod cleanall ts-install ts-typecheck ts-build ts-preview dev
