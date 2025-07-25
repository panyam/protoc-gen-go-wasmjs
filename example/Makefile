
buf:
	rm -Rf gen
	buf generate

test: buf
	cd gen/wasm/gen/wasm && sh build.sh
