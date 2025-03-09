.PHONY: build
build:
	go build -v ./...

.PHONY: start
start:
	docker-compose up --build --detach

.PHONY: watch
watch:
	docker-compose up --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: restart
restart: stop start

.PHONY: clean
clean:
	docker system prune -a

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	@go test -race -cover -count=3 ./...

.PHONY: coverage
coverage:
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	go-test-coverage --config=./.testcoverage.yml --badge-file-name=coverage.svg
