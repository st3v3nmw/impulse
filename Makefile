LEVELDB_PATH = "level.db"

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: test
test:
	cd src/ && go test -cover

.PHONY: build
build:
	cd src/ && go build -o pithered
	mv src/pithered .

.PHONY: run
run: build
	./pithered --leveldb=${LEVELDB_PATH} --verbose
