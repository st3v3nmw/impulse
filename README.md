# pithered

> adj. frustrated that you can’t force yourself to remember something, even
though it’s right on the tip of your tongue—wishing you could simply rifle
through your own files directly, rather than having to toss random scraps to
your team of mental archivists, who evidently need hours to sift through the
pile before they come up with an answer, just as you’re falling asleep.

Distributed Key Value Store in Go

## Features

### v0 (Basic Interface on a Single Machine)

- [x] HTTP Server
- [x] `Put(key, value)`
- [x] `Get(key)`
- [x] `Delete(key)`
- [ ] Redis as a read & write buffer

### v1 (Master-Slave Architecture)

- [ ] Replicas
- [ ] Chaos Testing

### v2 (Sharding)

- [ ] Sharding

### v3 (Complex Querying)

- [ ] Range Queries
- [ ] Search

### v4 (No Leader Nodes)

- [ ] [Raft Consensus](https://raft.github.io/)
