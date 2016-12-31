# blockchaind

This is a sample cli tool to submit transactions to the chain.  Each line via
stdin is submitted as a transaction.  There is also an http admin interface that
can be enabled by providing a bind address. See the [usage](#usage) for full options.


### Building
```
make all
```

This will generate a `blockchaind` binary in the current working directory.

### Running
Below shows how to run a single node and multiple nodes.

#### First Node

Start the first node with default values:

```
blockchaind
```

#### Additional Node
Spin up additional nodes by changing the bind port to avoid conflicts and providing
the address of the host above as the join address i.e. `-j`.

```
blockchaind -b 127.0.0.1:45455 -j 127.0.0.1:45454
```

### Transactions

To submit a transaction type something in the window and hit enter.   You should
start seeing blocks and transactions being broadcasted along with calls to the
finite state machine.

### Usage

```
blockchaind --help
```

### Development

```
go run cmd/*.go
```
