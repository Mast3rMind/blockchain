
NAME = blockchaind

clean:
	rm -f $(NAME)
	rm -rf vendor

test:
	go test -cover .

blockchaind:
	go build -o $(NAME) ./cmd/

all: clean blockchaind
