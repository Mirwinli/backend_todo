include .env
export
start-service:
	go run ./cmd/todo/main.go --config=./config/local.yaml
migrate-up:
	migrate -path migrations -database ${CONN-STRING} up
migrate-down:
	migrate -path migrations -database ${CONN-STRING} down