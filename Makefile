build:
	go build -o ${GOPATH}/bin/dfg main.go

clean:
	rm -rf ./bin

.PHONY: build clean