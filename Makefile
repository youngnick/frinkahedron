CLI_VERSION = 0.0.2
DAEMON_VERSION = 0.0.1
IMAGE = frinkahedron:$(DAEMON_VERSION)

.PHONY: build setup build-linux

build:
	CGO_ENABLED=0 go build -o bin/frinkacli -ldflags="-X main.version=${CLI_VERSION}" cmd/frinkacli/frinkacli.go
	CGO_ENABLED=0 go build -o bin/frinkahedron -ldflags="-X main.version=${CLI_VERSION}" cmd/frinkahedron/frinkahedron.go
