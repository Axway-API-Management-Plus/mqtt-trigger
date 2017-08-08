VERSION := 0.0.2
NAME := mqtt-trigger
DATE := $(shell date +'%Y-%M-%d_%H:%M:%S')
BUILD := $(shell git rev-parse HEAD | cut -c1-8)
LDFLAGS :=-ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.Date=$(DATE)"
IMAGE := jdavanne/$(NAME)
REGISTRY := davinci976
PUBLISH := $(REGISTRY)/$(IMAGE)
.PHONY: docker all

all: build

build:
	(cd src ; go build -o ../$(NAME) $(LDFLAGS))

dev:
	ls -d src/* | entr -r sh -c "make && ./mqtt-trigger --conf ./mqtt-trigger.yml"

docker-test:
	docker-compose -f docker-compose.test.yml down
	docker-compose -f docker-compose.test.yml build
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit || (docker-compose -f docker-compose.test.yml down ; exit 1)
	docker-compose -f docker-compose.test.yml down

docker-test-logs:
	docker-compose -f docker-compose.test.yml logs

clean:
	rm -f $(NAME) $(NAME).tar.gz

test:
	for dir in $$(find . -name "*_test.go" | grep -v ./vendor | xargs dirname | sort -u -r); do echo "$$dir..."; go test -v $$dir || exit 1 ; done | tee output.txt
	cat output.txt | egrep -- "--- FAIL:|--- SKIP:" || true

test-specific:
	go test -v $$(ls *.go | grep -v "_test.go") $(ARGS)

deps:
	go list -f '{{range .TestImports}}{{.}} {{end}} {{range .Imports}}{{.}} {{end}}' ./src/... | tr ' ' '\n' | grep -e "^[^/_\.][^/]*\.[^/]*/" |sort -u >.deps

deps-install:
	go get -v $$(cat .deps)
	#for dep in $$(cat .deps); do echo "installing '$$dep'... "; go get -v $$dep; done

deps-install-force: deps
	go get -u -v $$(cat .deps)
	#for dep in $$(cat .deps); do echo "installing '$$dep'... "; go get -u -v $$dep; done

docker-run:
	docker-compose up

docker:
	docker build -t $(IMAGE) .

docker-publish-all: docker-publish docker-publish-version

docker-publish-version:
	docker tag $(IMAGE) $(PUBLISH):$(VERSION)
	docker push $(PUBLISH):$(VERSION)

docker-publish: docker
	docker tag $(IMAGE) $(PUBLISH):latest
	docker push $(PUBLISH):latest
