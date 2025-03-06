.PHONY:
build:
	go build -v ./...

.PHONY:
start:
	docker-compose up --build --detach

.PHONY:
stop:
	docker-compose down

.PHONY:
restart: stop start

.PHONY:
clean:
	docker system prune -a

.PHONY:
lint:
	golangci-lint run

.PHONY:
test:
	@go test -race -cover -count=3 ./...

post:
	curl -d 'http://example.com/' -X POST http://localhost:8080
