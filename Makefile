LINTER_IMAGE=golangci/golangci-lint

.PHONY: lint
lint:
	docker run --rm -w /go/in -v $(CURDIR):/go/in $(LINTER_IMAGE) golangci-lint --config=build/golangci-lint.yaml run