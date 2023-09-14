TEST?=$$(go list ./... |grep -v 'vendor|e2e')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

VERSION ?= v0.0.1
OUTPUT_TYPE ?= type=docker
BUILDPLATFORM ?= linux/amd64
IMG_TAG ?= $(subst v,,$(VERSION))
REGISTRY ?= stackpath.com
IMG_NAME ?= virtual-kubelet
IMAGE ?= $(REGISTRY)/$(IMG_NAME)

BUILD_DATE ?= $(shell date '+%Y-%m-%dT%H:%M:%S')
VERSION_FLAGS := "-ldflags=-X main.buildVersion=$(IMG_TAG) -X main.buildTime=$(BUILD_DATE)"

default: build

build:
	@echo "==> Building..."
	go build -o bin/virtual-kubelet cmd/virtual-kubelet/*
	@echo
test:
	@echo "==> Running tests..."
	go test -v $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
	@echo

vet:
	@echo "==> Checking for suspicious Go constructs..."
	@echo "go vet \$$(go list ./... | grep -v vendor/)"
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -ne 0 ]; then \
		echo; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		echo; \
		exit 1; \
	fi
	@echo

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w $(GOFMT_FILES)
	@echo

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements..."
	gofmt_files=$$(gofmt -l $$(find . -name '*.go' | grep -v vendor))
	@if [[ -n "$(gofmt_files)" ]]; then \
        echo; \
		echo "gofmt needs running on the following files:"; \
		echo "$(gofmt_files)"; \
		echo "You can use the command: \`make fmt\` to reformat code."; \
		echo; \
		exit 1; \
	fi;
	@echo

generate:
	@which swagger ; if [ $$? -ne 0 ] ; then \
		echo "Please install go-swagger to generate StackPath API client code"; \
		echo "See: https://goswagger.io/install.html"; \
		echo; \
		exit 1; \
	fi

	@echo "==> Generating code from StackPath API swagger specs..."
	swagger generate client \
		--spec=swagger/stackpath_workload.oas2.json \
		--target=internal/api/workload \
		--model-package=workload_models \
		--client-package=workload_client

	@echo "==> Generating mocks..."

	mockgen -mock_names ClientService=InstancesClientService -destination=internal/mocks/mock_instances_client.go \
		-package=mocks github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/instances ClientService
	mockgen -mock_names ClientService=WorkloadClientService -destination=internal/mocks/mock_workload_client.go \
		-package=mocks github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/workload ClientService
	mockgen -mock_names ClientService=InstanceClientService -destination=internal/mocks/mock_instance_client.go \
		-package=mocks github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/instance ClientService
	mockgen -mock_names ClientService=InstanceLogsClientService -destination=internal/mocks/mock_instance_logs_client.go \
		-package=mocks github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/instance_logs ClientService
	mockgen -destination=internal/mocks/mock_k8s_listers.go -package=mocks k8s.io/client-go/listers/core/v1 ComponentStatusLister,ConfigMapLister,ConfigMapNamespaceLister,EndpointsLister,EndpointsNamespaceLister,EventLister,EventNamespaceLister,LimitRangeLister,LimitRangeNamespaceLister,NamespaceLister,NodeLister,PersistentVolumeLister,PersistentVolumeClaimLister,PersistentVolumeClaimNamespaceLister,PodLister,PodNamespaceLister,PodTemplateLister,PodTemplateNamespaceLister,ReplicationControllerLister,ReplicationControllerNamespaceLister,ResourceQuotaLister,ResourceQuotaNamespaceLister,SecretLister,SecretNamespaceLister,ServiceLister,ServiceNamespaceLister,ServiceAccountLister,ServiceAccountNamespaceLister


build-image: test vet fmt fmtcheck
	docker buildx build \
		--file docker/Dockerfile \
		--build-arg VERSION_FLAGS=$(VERSION_FLAGS) \
		--output=$(OUTPUT_TYPE) \
		--platform="$(BUILDPLATFORM)" \
		--pull \
		--tag $(IMAGE):$(IMG_TAG) .

.PHONY: build test vet fmt fmtcheck build-image
