.PHONY: build run-stdio run-http clean

BINARY=yapi-mcp

build:
	go build -o $(BINARY) ./cmd/yapi-mcp

run-stdio: build
	./$(BINARY) -transport stdio -config config.yaml

run-http: build
	./$(BINARY) -transport http -config config.yaml -port 8080

clean:
	rm -f $(BINARY)
	rm -f $(BINARY).exe
