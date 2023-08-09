.PHONY: build test clean

build:
	env GOOS=linux GOARCH=amd64 TAGS=logging go build -ldflags="-s -w" -tags="$(TAGS)" -o bin/handler cmd/main.go

test:
	cfn generate
	env GOOS=linux GOARCH=amd64 TAGS=logging go build -ldflags="-s -w" -o bin/handler cmd/main.go

clean:
	rm -rf bin
