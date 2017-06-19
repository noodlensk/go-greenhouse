HAS_GLIDE := $(shell command -v glide;)
IMAGE_NAME ?= noodlensk/hydroponics
.PHONY: hookInstall
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o hydroponics ./main.go

.PHONY: buildLinux
buildLinux:
	env GOOS=linux GOARCH=amd64 go build -o hydroponics_linux ./main.go

.PHONY: buildPi
buildPi:
	env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o hydroponics_pi ./main.go
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor
