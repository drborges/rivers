test:
	go test ./... -v -run=$(grep)
build:
	go build ./...
update:
	go get -u ./...
delete-branches:
	git branch | grep -v master | xargs -I {} git branch -D {}