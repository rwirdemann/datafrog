build:
	go build -o ${GOPATH}/bin/dfgapi cmd/dfgapi/main.go
	go build -o ${GOPATH}/bin/dfgweb cmd/dfgweb/main.go

clean:
	rm -rf ./bin

.PHONY: build clean