
buf: parent
	rm -Rf gen
	buf generate

parent:
	cd ../../ && make

test: buf
	cd gen/wasm/gen/wasm && sh build.sh
