
# Default: production mode using published wasmjs protos + local plugin  
buf: 
	rm -Rf gen
	buf generate

# Development mode: uses local plugin and symlinked wasmjs annotations
buf-dev: parent symlink-wasmjs
	rm -Rf gen
	# Update dependencies for dev config to avoid conflicts
	rm -f buf.lock
	buf generate --template buf.gen.dev.yaml --config buf.dev.yaml

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
