build:
	@go mod tidy
	@echo "Building:"
	@go build -o ./bin/main.bin ./main.go
	@echo "Applying: Permissions"
	@chmod +x ./bin/main.bin

run:
	@./bin/server.bin

clean:
	@rm -rf ./bin/*