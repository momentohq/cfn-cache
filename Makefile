.PHONY: build test clean

build: clean
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -tags="logging,lambda.norpc,$(TAGS)" -o bin/bootstrap cmd/main.go

test:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags="logging,lambda.norpc,$(TAGS)" -o bin/bootstrap cmd/main.go

clean:
	rm -rf bin
