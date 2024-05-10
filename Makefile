BINARY_NAME=main.out
BINARY_PATH=bin

build:
	go build -o ${BINARY_PATH}/${BINARY_NAME} ./cmd/api/

run:
	@make build 
	${BINARY_PATH}/${BINARY_NAME}

clean:
	go mod tidy
	go clean
	rm ${BINARY_PATH}/${BINARY_NAME}
