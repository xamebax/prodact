ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NAME=prodact

build:
	cd ${ROOT_DIR}/cmd/${NAME} && go build -o ${NAME} .

clean:
	rm -f ${ROOT_DIR}/${NAME}/${NAME}

lint:
	goreportcard-cli -v
