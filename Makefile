default: build

build: imports vet
	go build -v -o ./yandex-dns-cli-manager .

doc:
	godoc -http=:6060 -index

fmt:
	go fmt ./

# https://godoc.org/golang.org/x/tools/cmd/goimports
# go get golang.org/x/tools/cmd/goimports
imports:
	goimports -w ./

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

run: build
	./yandex-dns-cli-manager

test:
	go test ./

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./