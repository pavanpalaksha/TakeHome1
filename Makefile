export GOTRACEBACK := all
GO := go
SERVER := server
CLIENT := client

all: build

build:
	$(GO) build -o $(SERVER) ./cmd/server
	$(GO) build -o $(CLIENT) ./cmd/client	

deps:
	$(GO) get -u github.com/golang/glog
	$(GO) get -u google.golang.org/grpc

clean:
	$(GO) clean -i ./...
	rm -f $(SERVER) $(CLIENT)
