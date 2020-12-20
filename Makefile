# https://tutorialedge.net/golang/makefiles-for-go-developers/
# https://devhints.io/makefile
# https://gist.github.com/isaacs/62a2d1825d04437c6f08

PROJECT = $(notdir $(CURDIR))

hello:
	@echo "Hello, Nutzi!"

build:
	go build -o bin/$(PROJECT) .

run: build
	go run .

test:
	go test -v ./...

# use live-server: 
# https://github.com/txchen/light-server
html:
	light-server -c live-server.conf.json
