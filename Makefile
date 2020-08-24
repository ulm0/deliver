VERSION=$(CI_COMMIT_REF_NAME)-$(CI_COMMIT_SHORT_SHA)
CONTAINER_REGISTRY=$(CI_REGISTRY)
CONTAINER_IMAGE=$(CI_REGISTRY_IMAGE)
ifdef CI_COMMIT_TAG
CONTAINER_REGISTRY=docker.io
CONTAINER_IMAGE=ulm0/deliver
VERSION=$(shell ./hack/version)
endif
LD_FLAGS=-s -w -extldflags -static -X github.com/ulm0/deliver/pkg/cli.Version=$(VERSION)
CONTAINER_IMAGE_FULL=$(CONTAINER_REGISTRY)/$(CONTAINER_IMAGE)
CONTAINER_LIST=$(CONTAINER_IMAGE_FULL):$(VERSION)
ifdef CI_COMMIT_TAG
# cheap fix for multiple tags
CONTAINER_LIST+=-t $(CONTAINER_IMAGE_FULL):latest
endif

build:
	@CGO_ENABLED=0 go build -a -ldflags="$(LD_FLAGS)" -installsuffix cgo -o build/deliver github.com/ulm0/deliver/cmd/deliver

build-docker:
	@echo $(CONTAINER_LIST)
	@docker build --build-arg VERSION=$(VERSION) -t $(CONTAINER_LIST) .

push-docker:
	@docker push $(CONTAINER_IMAGE_FULL)
