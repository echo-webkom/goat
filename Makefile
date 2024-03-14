BULD_DIR=build

.PHONY: install build run test tidy fmt clean

install:
	go mod download

build:
	go build -o $(BULD_DIR)/server cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./test/**/* -v -cover

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

clean:
	rm -rf $(BULD_DIR)

