VERSION=$(shell ./hack/version)
LD_FLAGS=-s -w -extldflags -static -X github.com/ulm0/deliver/pkg/cli.Version=$(VERSION)

build:
	@CGO_ENABLED=0 go build -a -ldflags="$(LD_FLAGS)" -installsuffix cgo -o build/deliver github.com/ulm0/deliver/cmd/deliver

build-docker:
	@docker build --build-arg VERSION=$(VERSION) -t ulm0/deliver:$(VERSION) .

push-docker:
	@docker push ulm0/deliver:$(VERSION)
