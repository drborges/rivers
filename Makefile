default:
	go build

install:
	go install

clean:
	go clean

test:
	go test -race ./...
