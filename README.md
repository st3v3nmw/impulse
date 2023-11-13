# impulse

> Wer Ordnung hält, ist nur zu faul zum Suchen. \
> (If you keep things tidily ordered, you’re just too lazy to go searching.) \
> — German proverb, Designing Data-Intensive Applications

Distributed Key-Value Store in Go

## Roadmap

- [x] HTTP Server API (with `fasthttp`)
  - [x] `Put(key, value)`
  - [x] `Get(key) -> value`
  - [x] `Delete(key)`
- [ ] Storage Engines
  - [x] Hash Map (with `sync.RWMutex`)
  - [x] LevelDB
  - [ ] LSM-Tree
      - [ ] SSTable
      - [ ] Bloom Filter
  - [ ] B Tree
- [ ] Single Leader Replication (leader election with [Raft Consensus Protocol](https://raft.github.io/))
- [ ] Leaderless Replication (peer-to-peer coordination)
- [ ] Chaos Testing (like Netflix's [Chaos Monkey](https://netflix.github.io/chaosmonkey/))
- [ ] Range Queries
- [ ] Sharding

## "Benchmarks"

### benchmarks/concurrent.sh

```console
> ./benchmarks/concurrent.sh
go build -o impulse ./...

HASH_MAP
========

PUT
Running 10s test @ http://127.0.0.1:3000/foo
  100 goroutine(s) running concurrently
224432 requests in 9.93928912s, 12.20MB read
Requests/sec:		22580.29
Transfer/sec:		1.23MB
Avg Req Time:		4.428641ms
Fastest Request:	120.027µs
Slowest Request:	33.188626ms
Number of Errors:	0

GET
Running 10s test @ http://127.0.0.1:3000/foo
  100 goroutine(s) running concurrently
249859 requests in 9.924716107s, 28.59MB read
Requests/sec:		25175.43
Transfer/sec:		2.88MB
Avg Req Time:		3.972126ms
Fastest Request:	99.222µs
Slowest Request:	51.840697ms
Number of Errors:	0

DELETE
Running 10s test @ http://127.0.0.1:3000/foo
  100 goroutine(s) running concurrently
234442 requests in 9.951180146s, 12.74MB read
Requests/sec:		23559.22
Transfer/sec:		1.28MB
Avg Req Time:		4.244623ms
Fastest Request:	112.367µs
Slowest Request:	33.148716ms
Number of Errors:	0

LEVELDB
=======

PUT
Running 10s test @ http://127.0.0.1:3000/foo
  100 goroutine(s) running concurrently
229573 requests in 9.940867386s, 12.48MB read
Requests/sec:		23093.86
Transfer/sec:		1.26MB
Avg Req Time:		4.330155ms
Fastest Request:	112.35µs
Slowest Request:	37.433849ms
Number of Errors:	0

GET
Running 10s test @ http://127.0.0.1:3000/foo
  100 goroutine(s) running concurrently
248876 requests in 9.9313757s, 28.48MB read
Requests/sec:		25059.57
Transfer/sec:		2.87MB
Avg Req Time:		3.990491ms
Fastest Request:	96.86µs
Slowest Request:	43.956081ms
Number of Errors:	0

DELETE
Running 10s test @ http://127.0.0.1:3000/foo
  100 goroutine(s) running concurrently
234623 requests in 9.952669339s, 12.75MB read
Requests/sec:		23573.88
Transfer/sec:		1.28MB
Avg Req Time:		4.241983ms
Fastest Request:	110.965µs
Slowest Request:	31.113818ms
Number of Errors:	0
```

### benchmarks/sequential.sh

```console
> N=2048 ./benchmarks/sequential.sh
go build -o impulse ./...

HASH_MAP
========
PUT: 104.59 requests per second
GET (present): 96.58 requests per second
DELETE: 104.70 requests per second
GET (missing): 96.21 requests per second

LEVELDB
=======
PUT: 104.58 requests per second
GET (present): 95.94 requests per second
DELETE: 104.30 requests per second
GET (missing): 95.17 requests per second
```

### benchmarks/resilience.sh
