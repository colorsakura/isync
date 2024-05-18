LDFLAGS := -s -w

all: env fmt test isync

env:
	@go version

fmt:
	go fmt ./...

test:
	go test ./...

isync:
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/isync cmd/isync/main.go

clean:
	rm -rf bin/isync
