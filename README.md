# impulse

> Wer Ordnung hält, ist nur zu faul zum Suchen. \
> (If you keep things tidily ordered, you’re just too lazy to go searching.) \
> — German proverb

Distributed Key Value Store in Go

## Features

### v-1 (Storage)

- [x] LevelDB
- [x] In Memory Map
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

## "Benchmark"

Results on a severely underpowered machine:

```console
> N=2048 ./benchmark.sh
go build -o impulse ./...

LEVELDB
=======
PUT: 90.32 requests per second
GET: 84.39 requests per second
DELETE: 91.39 requests per second
GET: 83.65 requests per second

IN_MEMORY_MAP
=============
PUT: 90.58 requests per second
GET: 79.94 requests per second
DELETE: 91.10 requests per second
GET: 80.21 requests per second
```
