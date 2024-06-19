build: 
	@go build -C cmd/scheduler -o ../../bin/scheduler

run: build
	@./bin/scheduler

test: 
	@go test -v ./...
