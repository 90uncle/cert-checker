PROJECT="cert-checker"

BINARY_NAME = bin/cert-checker

.PHONY: build-win
build-win:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe main.go

.PHONY: build-linux
build-linux:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME) main.go

.PHONY: build-all
build-all: build-win build-linux

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)*