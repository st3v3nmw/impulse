LEVELDB_PATH = "level.db"

.PHONY: format
format:
	gofmt -s -w .

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

.PHONY: deploy
deploy:
	docker build . -t localhost:32000/impulse:latest
	docker push localhost:32000/impulse
	microk8s kubectl apply -f deploy/
	microk8s kubectl rollout restart deployment impulse -n impulse
