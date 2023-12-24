LEVELDB_PATH = "level.db"

.PHONY: format
format:
	gofmt -w .

.PHONY: test
test:
	go test ./...  -coverpkg=./... -coverprofile ./coverage.out
	go tool cover -func ./coverage.out

.PHONY: build
build:
	go build -o impulse ./cmd/impulse

.PHONY: run
run: build
	./impulse --engine=LEVELDB --leveldb=${LEVELDB_PATH} --verbose
