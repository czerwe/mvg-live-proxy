# GOROOT := /usr/lib/go
# GOPATH := shell

DOCKERBIN := $(shell which docker)
GOBIN := $(shell which go)

buildindocker: buildimage clean
	${DOCKERBIN} run -v  $(shell pwd):/go/src/mvg-live-proxy --workdir /go/src/mvg-live-proxy abuild make build

dependencys:
	${GOBIN} get github.com/Sirupsen/logrus
	${GOBIN} get github.com/gorilla/mux
	${GOBIN} get github.com/jessevdk/go-flags
armbuild: clean
	env GOOS=linux CGO_ENABLED=0 GOARCH=arm ${GOBIN} build proxymvg.go

build: clean
	env GOOS=linux CGO_ENABLED=0 ${GOBIN} build proxymvg.go


buildimage:
	${DOCKERBIN} build -t abuild -f Dockerfile_build .

.PHONY: clean
clean:
	rm -rf proxymvg
