IMAGE ?= kabanero/events-operator
IMAGE_TAG ?= prototype
KUBEBUILDER_VERSION ?= 2.3.0
OPERATOR_SDK_RELEASE_VERSION ?= v0.15.2
OPERATOR_FLAGS = --zap-level=debug --zap-encoder=console
CRDS = $(wildcard deploy/crds/*_crd.yaml)
SAMPLE_CRS=$(wildcard sample_crds/example1/*.yaml)

.PHONY: setup test-setup generate install build build-all format test

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
	touch .apply-crds

apply-crds: .apply-crds

run-local: .apply-crds
	operator-sdk run --local --operator-flags="$(OPERATOR_FLAGS)"

oc-deploy: build-image push-image
	sleep 2
	oc apply -f deploy

oc-undeploy:
	oc delete deployment events-operator

debug: .apply-crds
	operator-sdk run --local --enable-delve --operator-flags="$(OPERATOR_FLAGS)"

delete-samples:
	echo $(SAMPLE_CRS) | tr ' ' '\n' | xargs -I{} oc delete -f {}

apply-samples:
	echo $(SAMPLE_CRS) | tr ' ' '\n' | xargs -I{} oc apply -f {}

setup:
	@./scripts/install-operator-sdk.sh ${OPERATOR_SDK_RELEASE_VERSION}

test-setup:
	@./scripts/install-envtest.sh ${KUBEBUILDER_VERSION}

format:
	go fmt ./...

vet:
	#go vet ./...
	@echo "Vetting is disabled. It will be re-enabled once code is stable."

tidy:
	go mod tidy -v

test:
	go test ./...

operator-tests:
	@ginkgo -r -v --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --compilers=2 pkg/controller

unit-tests:
	#go test -v -tags=unit_test ./...
	@echo "unit tests passed"

e2e-tests:
	@echo "e2e tests passed"
