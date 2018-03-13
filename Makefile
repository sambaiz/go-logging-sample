.PHONY: setup
setup:
	go get -u github.com/golang/lint/golint
	go get -u github.com/golang/dep/cmd/dep

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	go list ./... | xargs golint -set_exit_status
	
.PHONY: dep
dep:
	dep ensure
