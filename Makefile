all: build/i12e

build/i12e: ./cmd/i12e/main.go ./internal/install/i12e_install.go
	go build -trimpath -buildvcs=false -ldflags="-s -w" -o build/i12e ./cmd/i12e

clean:
	rm -rf build

.PHONY: all clean
