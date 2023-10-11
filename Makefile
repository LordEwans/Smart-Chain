build:
	@go build -o bin
	
run: build
	@./bin

test:
	@go test -v ./...