GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
SERVER_NAME=appServer
CLIENT_NAME=appClient

all: test build

build:
	cd server && $(GOBUILD) -o $(SERVER_NAME) -v ./...
	mv server/$(SERVER_NAME) $(SERVER_NAME)
	cd client && $(GOBUILD) -o $(CLIENT_NAME) -v ./...
	mv client/$(CLIENT_NAME) $(CLIENT_NAME)

test:
	cd server && $(GOTEST) -v ./...
	cd client && $(GOTEST) -v ./...

clean:
	rm -f $(SERVER_NAME)
	rm -f $(CLIENT_NAME)

run:
	cd server && $(GOBUILD) -o $(SERVER_NAME) -v ./...
	mv server/$(SERVER_NAME) $(SERVER_NAME)
	./$(SERVER_NAME) $$> server.log &


deps:
	$(GOGET) ./..
