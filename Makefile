all: test build

test:
	go test -v ./...

format:
	go fmt ./...

deps:
	go mod tidy

build:
	go build -o bin/janus

release:
	go build -ldflags "-s -w" -o bin/janus

content:
	cd _pages; npm run build

package: release content
	rm -rf _package/janus;
	mkdir -p _package/janus;
	cp -R bin _package/janus/bin;

	mkdir -p _package/janus/content;
	cp -R _template _package/janus/content/template;
	cp -R _static _package/janus/content/static;

	cp -R _pages/build _package/janus/content/static/component;

	mkdir -p _package/janus/util
	cp _package/janus.service _package/janus/util;

	tar -czf _package/janus.tar.gz -C _package janus;
