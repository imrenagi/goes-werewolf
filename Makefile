GOPATH=$(shell go env GOPATH)
NAMESPACE=dev
IMAGE_REGISTRY:=imrenagi
APP_NAME=goes-werewolf
GIT_COMMIT=$(shell git rev-parse --short HEAD)
ENV=dev
IMAGE_TAG=latest

.PHONY: all test clean

build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o $(APP_NAME) cmd/server/main.go

mock:
	mockgen -destination=internal/app/werewolf/mocks/mock_player_dao.go -package=mocks github.com/imrenagi/goes-werewolf/internal/app/werewolf/services PlayerDAO
	mockgen -destination=internal/app/werewolf/mocks/mock_services_polldao.go -package=mocks github.com/imrenagi/goes-werewolf/internal/app/werewolf/services PollDAO
	mockgen -destination=internal/app/werewolf/mocks/mock_state_context.go -package=mocks github.com/imrenagi/goes-werewolf/internal/app/werewolf/states Context
	mockgen -destination=internal/app/werewolf/mocks/mock_polling_service.go -package=mocks github.com/imrenagi/goes-werewolf/internal/app/werewolf/states PollingService

test:
	go test ./... -cover -vet -all -short

docker:
	docker build -t $(IMAGE_REGISTRY)/$(APP_NAME):latest -f build/package/Dockerfile .

docker-build:
	docker build -t $(IMAGE_REGISTRY)/$(APP_NAME):$(GIT_COMMIT) -f build/package/Dockerfile .
	docker tag $(IMAGE_REGISTRY)/$(APP_NAME):$(GIT_COMMIT) $(IMAGE_REGISTRY)/$(APP_NAME):latest

docker-push:
	docker push $(IMAGE_REGISTRY)/$(APP_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_REGISTRY)/$(APP_NAME):latest
