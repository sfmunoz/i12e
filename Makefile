all: bin/i12e

bin/i12e: ./cmd/i12e/main.go ./internal/install/i12e_install.go
	go build -trimpath -buildvcs=false -ldflags="-s -w" -o bin/i12e ./cmd/i12e

clean:
	rm -rf bin

.PHONY: all clean
