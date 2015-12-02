test:
	go test ./... -v -run=$(grep)
build:
	go build ./...
update:
	go get github.com/smartystreets/goconvey
delete-branches:
	git branch | grep -v master | xargs -I {} git branch -D {}