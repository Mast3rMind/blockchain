
test:
	go test -cover ./...

build:
	go build -o bcd ./cli/*.go
