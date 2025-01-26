# Set bin directory variable
BIN_DIR := ./bin

SOCKET := /tmp/ft_supervisor.sock

all: build

build:
	go build -o ${BIN_DIR}/supervisord cmd/supervisord/*.go
	go build -o ${BIN_DIR}/supervisorctl cmd/supervisorctl/*.go

clean:
	rm -rf ${BIN_DIR}/*
	rm ${SOCKET}
