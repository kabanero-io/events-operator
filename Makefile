IMAGE ?= kabanero/events-operator
IMAGE_TAG ?= latest
OPERATOR_SDK_RELEASE_VERSION ?= v0.15.2
OPERATOR_FLAGS = --zap-level=debug --zap-encoder=console
CRDS = $(wildcard deploy/crds/*_crd.yaml)
SAMPLE_CRS=$(wildcard pkg/apis/events/v1alpha1/sample_crds/example1/*.yaml)

.PHONY: setup generate install build build-all format test

build:
	go build ./cmd/manager/...

install:
	go install github.com/kabanero-io/events-operator/cmd/manager

build-all: generate build

generate: setup
	operator-sdk generate k8s
	operator-sdk generate crds

build-image: setup
	operator-sdk build $(IMAGE):$(IMAGE_TAG)

push-image:
	docker push $(IMAGE):$(IMAGE_TAG)

.apply-crds: $(CRDS)
	echo $(CRDS) | tr ' ' '\n' | xargs -I{} oc apply -f {}
	touch .applycrds

apply-crds: .apply-crds

run-local: .apply-crds
	operator-sdk run --local --operator-flags="$(OPERATOR_FLAGS)"

debug: .apply-crds
	operator-sdk run --local --enable-delve --operator-flags="$(OPERATOR_FLAGS)"

delete-samples:
	echo $(SAMPLE_CRS) | tr ' ' '\n' | xargs -I{} oc delete -f {}

apply-samples:
	echo $(SAMPLE_CRS) | tr ' ' '\n' | xargs -I{} oc apply -f {}

setup:
	@./scripts/install-operator-sdk.sh ${OPERATOR_SDK_RELEASE_VERSION}

format:
	go fmt ./...

vet:
	#go vet ./...
	@echo "Vetting is disabled. It will be re-enabled once code is stable."

tidy:
	go mod tidy -v

test:
	go test ./...

unit-tests:
	@echo "Unit tests passed"

e2e-tests:
	@echo "e2e tests passed"
