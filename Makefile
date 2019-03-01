CLI_VERSION = 0.0.2
DAEMON_VERSION = 0.0.1
IMAGE = frinkahedron:$(DAEMON_VERSION)

.PHONY: build setup build-linux

build:
	CGO_ENABLED=0 go build -o bin/frinkacli -ldflags="-X main.version=${CLI_VERSION}" cmd/frinkacli/frinkacli.go
	CGO_ENABLED=0 go build -o bin/frinkahedron -ldflags="-X main.version=${CLI_VERSION}" cmd/frinkahedron/frinkahedron.go

build-linux:
    # GODEBUG is here to force use of the pure go resolver, just in case.
	CGO_ENABLED=0 GODEBUG=netdns=go GOOS=linux go build -a -installsuffix cgo -o bin/frinkacli -ldflags="-X main.version=${CLI_VERSION}" cmd/frinkacli/frinkacli.go

setup:
	dep ensure -vendor-only

