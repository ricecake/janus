all: format lint compile test

release: prod-compile

compile: deps build

prod-compile: deps prod-build

quick: js-build go-build

test:
	go test -v ./...

format: go-fmt js-fmt

js-fmt:
	npm run format

go-fmt:
	go fmt ./...

lint: go-lint js-lint

go-lint:
	go vet

js-lint:
	echo "nope"

deps: js-deps go-deps

js-deps:
	npm install

go-deps:
	go mod tidy

js: js-deps js-build

js-build:
	mkdir -p util/content
	npm run build

js-build-prod:
	mkdir -p util/content
	npm run prod-build

go: go-deps go-build

go-build:
	go build -o bin/janus

go-build-prod:
	go build -ldflags "-s -w" -o bin/janus

build: js-build go-build

prod-build: js-build-prod go-build-prod
