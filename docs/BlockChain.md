# blockchain

Below is a brief of the [original paper](https://bitcoin.org/bitcoin.pdf) in a more general
context such that it could be applied to non-financial use cases.

This can be useful in distributed systems over the internet and the security there in.

## Proof of Work

The proof-of-work involves scanning for a value that when hashed, such as with SHA-256, the
hash begins with a number of zero bits. The average work required is exponential in the number
of zero bits required and can be verified by executing a single hash.

For our timestamp network, we implement the proof-of-work by incrementing a nonce in the
block until a value is found that gives the block's hash the required zero bits. Once the CPU
effort has been expended to make it satisfy the proof-of-work, the block cannot be changed
without redoing the work. As later blocks are chained after it, the work to change the block
would include redoing all the blocks after it.

The proof-of-work also solves the problem of determining representation in majority decision
making. If the majority were based on one-IP-address-one-vote, it could be subverted by anyone
able to allocate many IPs. Proof-of-work is essentially one-CPU-one-vote.

The majority decision is represented by the longest chain, which has the greatest proof-of-work
effort invested in it. If a majority of CPU power is controlled by honest nodes, the honest
chain will grow the fastest and outpace any competing chains. To modify a past block, an attacker
would have to redo the proof-of-work of the block and all blocks after it and then catch up with
and surpass the work of the honest nodes.

To compensate for increasing hardware speed and varying interest in running nodes over time,
the proof-of-work difficulty is determined by a moving average targeting an average number of
blocks per hour. If they're generated too fast, the difficulty increases.

## Incentive

By convention, the first transaction in a block is a special transaction that starts a new coin
owned by the creator of the block. This adds an incentive for nodes to support the
network and provides a way to initially distribute coins into circulation, since there
is no central authority to issue them.  The steady addition of a constant of amount of new coins
is analogous to gold miners expending resources to add gold to circulation. In our case, it is CPU
time and electricity that is expended.  The incentive can also be funded with transaction fees. If
the output value of a transaction is less than its input value, the difference is a transaction fee
that is added to the incentive value of the block containing the transaction. Once a predetermined
number of coins have entered circulation, the incentive can transition entirely to transaction
fees and be completely inflation free.

The incentive may help encourage nodes to stay honest. If a greedy attacker is able to
assemble more CPU power than all the honest nodes, he would have to choose between using it
to defraud people by stealing back his payments, or using it to generate new coins. He ought to
find it more profitable to play by the rules, such rules that favor him with more new coins than
everyone else combined, than to undermine the system and the validity of his own wealth.

## Reclaiming Disk Space

Once the latest transaction in a coin is buried under enough blocks, the spent transactions before
it can be discarded to save disk space. To facilitate this without breaking the block's hash,
transactions are hashed in a [Merkle Tree](https://en.wikipedia.org/wiki/Merkle_tree), with only
the root included in the block's hash.  Old blocks can then be compacted by stubbing off branches
of the tree. The interior hashes do not need to be stored.

A block header with no transactions would be about 80 bytes. If we suppose blocks are
generated every 10 minutes, 80 bytes * 6 * 24 * 365 = 4.2MB per year. With computer systems
typically selling with 2GB of RAM as of 2008, and Moore's Law predicting current growth of
1.2GB per year, storage should not be a problem even if the block headers must be kept in
memory.
