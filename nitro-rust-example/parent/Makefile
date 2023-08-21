ARCH := x86_64
# ARCH := aarch64
RUST_DIR := $(shell readlink -m $(shell dirname $(firstword $(MAKEFILE_LIST))))

build:
	curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh -s -- -y
	rustup target install ${ARCH}-unknown-linux-gnu
	cargo build --manifest-path=${RUST_DIR}/Cargo.toml --target=${ARCH}-unknown-linux-gnu --release

client:
	docker build --build-arg TARGETPLATFORM=$(ARCH)-unknown-linux-gnu -t enclave-client -f Dockerfile .
	nitro-cli build-enclave --docker-uri enclave-client:latest --output-file enclave-client.eif

.PHONY: clean
clean:
	cargo clean
