
NAME = blockchaind

clean:
	[ -x $(NAME) ] && rm $(NAME)

test:
	go test -cover .

blockchaind:
	go build -o $(NAME) ./cmd/

all: clean blockchaind
