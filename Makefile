GIT_TAG=$(shell git describe --abbrev=0 --tags)

build:
	cd ./clank && go build -ldflags "-X main.version=$(GIT_TAG)" && go install

docker-build:
	docker build -t lpxxn/clank:$(GIT_TAG) -f Dockerfile .
docker-build-latest:
	docker build -t lpxxn/clank:latest -f Dockerfile .

docker-buildx-latest:
	docker buildx build --platform linux/amd64,linux/arm64 -t lpxxn/clank:latest -f Dockerfile .

