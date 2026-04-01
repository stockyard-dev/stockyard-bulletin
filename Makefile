build:
	CGO_ENABLED=0 go build -o bulletin ./cmd/bulletin/

run: build
	./bulletin

test:
	go test ./...

clean:
	rm -f bulletin

.PHONY: build run test clean
