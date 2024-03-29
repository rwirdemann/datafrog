build:
	go build -o ${GOPATH}/bin/dbd main.go

clean:
	rm -rf ./bin

.PHONY: build clean