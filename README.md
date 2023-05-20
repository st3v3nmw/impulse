# impulse

> Wer Ordnung hält, ist nur zu faul zum Suchen. \
> (If you keep things tidily ordered, you’re just too lazy to go searching.) \
> — German proverb

Distributed Key Value Store in Go

## Features

### v-1 (Storage)

- [x] LevelDB
- [ ] In Memory Map
- [ ] SSTable
- [ ] LSM-Tree
- [ ] B-Tree

### v0 (Basic Interface on a Single Machine)

- [x] HTTP Server
- [x] `Put(key, value)`
- [x] `Get(key)`
- [x] `Delete(key)`

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
