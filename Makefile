GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
SOURCE=cmd/portfolio/main.go
TARGET=/usr/local/bin/portfolio

.PHONY: build test clean

all: build test

build: 
	$(GOBUILD) -o $(TARGET) -v $(SOURCE)

test: 
	$(GOTEST) -v ./...

clean:
	rm $(TARGET)