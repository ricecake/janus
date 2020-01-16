all: test build content

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

	mkdir -p _package/janus/assets;
	cp -R _template _package/janus/assets/template;
	cp -R _static _package/janus/assets/static;
	cp -R _pages/build _package/janus/assets/content;

	mkdir -p _package/janus/util
	cp _package/janus.service _package/janus/util;
	cp _package/backupCommand.sh _package/janus/util;
	cp -R _schema _package/janus/util/schema;

	tar -czf _package/janus.tar.gz -C _package janus;
