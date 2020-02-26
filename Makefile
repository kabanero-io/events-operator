.PHONY: generate build build-all format test

CRDS = $(wildcard deploy/crds/*_crd.yaml)
OPERATOR_FLAGS = --zap-level=debug --zap-encoder=console

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

.install: $(CRDS)
	echo $(CRDS) | tr ' ' '\n' | xargs -I{} oc apply -f {}
	touch .install

install: .install

run-local: .install
	operator-sdk run --local --operator-flags="$(OPERATOR_FLAGS)"

debug: .install
	operator-sdk run --local --enable-delve --operator-flags="$(OPERATOR_FLAGS)"
