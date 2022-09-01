BINARY_NAME=build/main.out

build:
	go get
	go build -o ${BINARY_NAME} main.go

clean:
	go clean
	rm ${BINARY_NAME}