.PHONY: generate build build-all format test

format:
	go fmt ./cmd/... ./pkg/...

test:
	go test ./cmd/... ./pkg/...

generate:
	operator-sdk generate k8s
	operator-sdk generate crds

build:
	go build ./cmd/manager/...

build-all: generate build

install:
	find deploy/crds -name "*_crd.yaml" | xargs -I{} oc apply -f {}

run-local: install
	operator-sdk run --local


