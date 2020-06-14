SERVICE_NAME     =starter
RELEASE_VERSION  =v0.1.3
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: mod test run build exec image show imagerun lint clean, tag
all: test

mod: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

test: mod ## Tests the entire project 
	go test -v -count=1 -race ./...
	# go test -v -count=1 -run TestMakeCPUEvent ./...

run: mod ## Runs the uncompiled code
	go run handler.go main.go 

build: mod ## Builds local release binary
		env CGO_ENABLED=0 go build -ldflags "-X main.Version=$(RELEASE_VERSION)" \
    	-mod vendor -o bin/$(SERVICE_NAME) .

exec: build ## Builds binary and runs it in Dapr
	dapr run --app-id $(SERVICE_NAME) \
         --app-port 8080 \
         --protocol http \
				 --port 3500 \
         --components-path ./config \
         bin/$(SERVICE_NAME) 

event: ## Publishes sample message to Dapr pubsub API 
	curl -v -d '{ "message": "hello" }' \
     -H "Content-type: application/json" \
     "http://localhost:3500/v1.0/publish/events"

image: mod ## Builds and publish docker image 
	docker build --build-arg VERSION=$(RELEASE_VERSION) \
		-t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

clean: ## Cleans all runtime generated directory
	go clean
	rm -fr ./bin/*

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
