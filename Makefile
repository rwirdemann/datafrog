build:
	go build -o ${GOPATH}/bin/rt-record record/main.go
	go build -o ${GOPATH}/bin/rt-listen listen/main.go

clean:
	rm -rf ./bin

.PHONY: build clean