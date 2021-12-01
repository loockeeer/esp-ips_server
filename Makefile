BINARY_NAME=build/main.out
build:
	GOOS=linux go build -o ${BINARY_NAME} src/main.go

clean:
	go clean
	rm ${BINARY_NAME}

format:
	go fmt espips_server/...