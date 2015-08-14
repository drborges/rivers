test:
	goapp test ./... -v -run=$(grep)
build:
	goapp build ./...
update:
	goapp get -u ./...
delete-branches:
	git branch | grep -v master | xargs -I {} git branch -D {}