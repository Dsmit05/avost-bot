APP = avost-bot
SERVICE_PATH = github.com/Dsmit05
BR = `git rev-parse --symbolic-full-name --abbrev-ref HEAD`
VER = `git describe --tags --abbrev=0`
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
FORMAT = $(VER)-$(TIMESTM)
DOCTAG = $(VER)-$(BR)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP) -ldflags "-X $(SERVICE_PATH)/$(APP)/internal/config.BuildVersion=$(FORMAT)" cmd/bot-v1/main.go

.PHONY: build-image
build-image:
	docker build -t $(APP):$(DOCTAG) .

.PHONY: run-app
run-app:
	docker run -d -v ${HOME}/APPS/$(APP)/data:/data --name=$(APP)-$(VER) $(APP):$(DOCTAG)

.PHONY: del-app
del-app:
	docker rm $(APP)-$(VER)

.PHONY: test_api
test_api:
	go test -tags=test_api  ./...

.PHONY: swag_init
swag_init:
	swag init -g cmd/bot-v1/main.go