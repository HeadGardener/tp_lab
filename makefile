build:
	go build -o ./bin -v ./cmd/tp_lab

run:
	go run ./cmd/tp_lab

test:
	gotest -v ./...

vendor:
	go mod vendor -v