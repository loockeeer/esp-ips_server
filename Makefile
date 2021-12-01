BINARY_NAME=build/main
build:
	go build -o ${BINARY_NAME} main.go

clean:
	go clean
	rm ${BINARY_NAME}

format:
	go fmt espips_server/...