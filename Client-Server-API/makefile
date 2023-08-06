build:
	@go mod tidy
	@echo "Building: Server"
	@go build -o ./bin/server.bin ./server.go
	@echo "Building: Client"
	@go build -o ./bin/client.bin ./client.go
	@echo "Applying: Permissions"
	@chmod +x ./bin/server.bin
	@chmod +x ./bin/client.bin

run-server:
	@./bin/server.bin

run-client:
	@./bin/client.bin

clean:
	@rm -rf ./bin/*
	@rm -rf ./cotacao.db
	@rm -rf ./cotacao.txt