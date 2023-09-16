# Version (set during CI/CD)
VERSION ?= v0.0.1

test:
	go test ./... -coverprofile _cover.out
	go tool cover -func _cover.out | tail -n 1

vet:
	go vet ./...
	gosec -quiet -exclude=G304 ./...

build:
	GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/jo-$(VERSION)-linux-amd64 github.com/mmlt/jo

deploy-local:
	sudo cp ./bin/jo-$(VERSION)-linux-amd64 /usr/local/bin/jo
