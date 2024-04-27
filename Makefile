build:
	go build -o ${GOPATH}/bin/dfg main.go
	go build -o ${GOPATH}/bin/dfgapi httpx/main.go

clean:
	rm -rf ./bin

.PHONY: build clean