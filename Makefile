GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
SOURCE=cmd/portfolio/main.go
TARGET=portfolio
OUTPUT=/usr/local/bin/$(TARGET)

.PHONY: build test clean

all: build test

build: 
	$(GOBUILD) -o $(OUTPUT) -v $(SOURCE)

test: 
	$(GOTEST) -v ./...

compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(TARGET)-darwin $(SOURCE)
	GOOS=linux GOARCH=amd64 go build -o bin/$(TARGET)-linux $(SOURCE)
	GOOS=windows GOARCH=amd64 go build -o bin/$(TARGET)-windows.exe $(SOURCE)

clean:
	rm $(TARGET)