include help.mk

NAME           	= $(shell basename $(CURDIR))
BUILD         	= $(shell git rev-parse --short HEAD)
DATE          	= $(shell date -uIseconds)
VERSION  	  	= $(shell git describe --always --tags)

IMAGE_BUILD    	= $(NAME):$(BUILD)
IMAGE_VERSION   = $(NAME):$(VERSION)

NETWORK_NAME	= network_$(NAME)$(PIPELINE_ID)
MONGO_NAME		= mongodb_$(NAME)$(PIPELINE_ID)

.PHONY: clean
clean: ##@dependencies Remove locally generated files.
	rm -rf vendor

.PHONY: install
install: clean  ##@dependencies Download dependencies.
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

.PHONY: lint
lint: ##@check Run lint on docker.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain --tag $(IMAGE_BUILD) --target=lint .

.PHONY: audit
audit: ##@check Run audit on docker.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain --tag $(IMAGE_BUILD) --target=audit .

.PHONY: test-unit
test-unit: ##@check Run unit tests on docker.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain --tag $(IMAGE_BUILD) --target=test-unit .

.PHONY: test-integration
test-integration: ##@check Run integration tests on docker.
	-docker network create $(NETWORK_NAME)
	-docker run --rm -d --name $(MONGO_NAME) --network $(NETWORK_NAME) mongo:6 --nojournal

	@(docker build --progress=plain --build-arg MONGO_URL=mongodb://${MONGO_NAME}:27017 --network $(NETWORK_NAME) --tag $(IMAGE_BUILD) --target=test-integration .; \
	status=$$?; \
	docker rm -vf $(MONGO_NAME); \
	docker network rm $(NETWORK_NAME); \
	exit $${status})	

.PHONY: build
build: ##@build Build image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain --build-arg VERSION=$(VERSION) --build-arg BUILD=$(BUILD) --build-arg DATE=$(DATE) --tag $(IMAGE_BUILD) --target=build .

.PHONY: image
image: ##@build Create release docker image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain --build-arg VERSION=$(VERSION) --build-arg BUILD=$(BUILD) --build-arg DATE=$(DATE) --tag $(IMAGE_BUILD) --target=image .

.PHONY: tag
tag: image ##@build Add docker tag in release image.
	docker tag $(IMAGE_BUILD) $(IMAGE_VERSION)

.PHONY: run
run: ##@run Run application on docker compose.
	DOCKER_BUILDKIT=1 \
	COMPOSE_DOCKER_CLI_BUILD=1 \
	docker compose up --build --remove-orphans


.PHONY: stop
stop: ##@run Stop application running on docker compose.
	docker compose down