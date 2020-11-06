.PHONY : build run fresh test clean

BIN := grafana-export

build:
	go build -o ${BIN} -ldflags="-s -w"

clean:
	go clean
	- rm -f ${BIN}

dist:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(MAKE) build

fresh: clean build run

lint:
	gofmt -s -w .
	find . -name "*.go" -exec ${GOPATH}/bin/golint {} \;

run:
	./${BIN}

test: lint
	go test
