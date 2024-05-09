server:
	@go build -o bin/app main.go
	@./bin/app

test:
	@go test -v ./...

db:
	@docker run -d \
    --name db-container \
    -p 5432:5432 \
    -e POSTGRES_DB=empdb \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=securepass \
    -v postgres_data:/var/lib/postgresql/data \
    postgres:latest

.PHONY: server test db