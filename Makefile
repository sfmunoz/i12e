FILES = ./cmd/i12e/main.go ./internal/install/install.go ./internal/install/i12e_install.go ./internal/install/k3s_install.go ./internal/cmdutil/cmdutil.go ./internal/fsutil/fsutil.go ./internal/pull/pull.go

all: build/i12e

build/i12e: $(FILES)
	go build -trimpath -buildvcs=false -ldflags="-s -w" -o build/i12e ./cmd/i12e

clean:
	rm -rf build

.PHONY: all clean
