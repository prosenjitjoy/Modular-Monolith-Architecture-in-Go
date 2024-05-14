tools:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

generate:
	go generate ./...

postgres:
	podman run --name postgres --hostname postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=malldb -v ./database/scripts:/docker-entrypoint-initdb.d -p 5432:5432 -d postgres:latest