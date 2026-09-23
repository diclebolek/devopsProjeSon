.PHONY: test test-rust test-go fmt

test: test-rust test-go

test-rust:
	cargo test --manifest-path firmware/Cargo.toml

test-go:
	cd pipeline && go test -race ./...

fmt:
	cargo fmt --all --manifest-path firmware/Cargo.toml
	cd pipeline && gofmt -w .
