
This is a sample cli tool.  Each line via stdin is submitted as a transaction.


### Running
Below shows how to run a single node and multiple nodes.

#### First Node
```
go run cmd/*.go
```

#### Additional Node
Spin up additional nodes by changing the port to avoid conflict.

```
go run cmd/*.go -b 127.0.0.1:45455 -j 127.0.0.1:45454
```

### Usage

```
go run cmd/*.go --help
```
