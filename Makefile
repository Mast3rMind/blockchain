
clean:
	rm -rf ./build

test:
	go test -cover ./...

build:
	[ -d ./build ] || mkdir ./build

build/blockchaind: build
	go build -v -o build/blockchaind ./cmd/blockchaind/*.go

all: build/blockchaind

