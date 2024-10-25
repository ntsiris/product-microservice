build:
	@go build -o bin/product cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/product

clean:
	@ rm bin/product