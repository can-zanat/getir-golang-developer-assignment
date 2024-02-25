build:
	go build .

run:
	go run .

lint:
	golangci-lint run -c .local/.golangci.yml

unit-test:
	go test ./... -short

generate-mocks:
	mockgen -source=./internal/service.go -destination=./internal/mock_service.go -package=internal
	mockgen -source=./internal/handler.go -destination=./internal/mock_handler.go -package=internal