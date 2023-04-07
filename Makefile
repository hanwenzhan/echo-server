test:
	@go mod tidy
	@go fmt ./...
	@go test ./...

update:
	@go get -u ./...

##### build docker image
PROJECT_NAME  ?= line-echo-server
VERSION       ?= v0.0.0
IMG_NAME      := registry.icetech.com.tw:5000/$(PROJECT_NAME):$(VERSION)
PLATFORM      := linux/amd64,linux/arm64

.xbuilder:
	@docker ps --format "{{.ID}} {{.Image}} {{.Names}} {{.Command}}" --filter "NAME=buildx_buildkit_xbuilder0"
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	@-docker buildx create --use --name xbuilder
	@docker buildx inspect --bootstrap

build-image: .xbuilder
	@echo "IMG_NAME=$(IMG_NAME)"
	@docker buildx build --platform $(PLATFORM) \
		--build-arg "VERSION=$(VERSION)" \
		--build-arg "BASED_ON=alpine" -t $(IMG_NAME) . --push

