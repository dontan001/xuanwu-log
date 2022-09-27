# Image URL to use all building/pushing image targets
IMG ?= dontan001/xuanwu-log
IMAGE_TAG_API ?= api
IMAGE_TAG_SCHD ?= schedule

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build API

.PHONY: api-build
api-build: fmt vet ## Build api binary.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/xuanwu-api ./cmd/api/main.go

.PHONY: api-run
api-run: fmt vet ## Run a api server from your host.
	go run ./cmd/api/main.go

.PHONY: api-image
api-image: api-build ## Build docker image with the api.
	docker build -t ${IMG}:${IMAGE_TAG_API} -f docker/api/Dockerfile .

.PHONY: api-push
api-push: ## Push docker image with the api.
	docker push ${IMG}:${IMAGE_TAG_API}

##@ Build Schedule

.PHONY: schedule-build
schedule-build: fmt vet ## Build schedule binary.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/xuanwu-backup ./cmd/schedule/main.go

.PHONY: schedule-run
schedule-run: fmt vet ## Run a schedule from your host.
	go run ./cmd/schedule/main.go

.PHONY: schedule-image
schedule-image: schedule-build ## Build docker image with the schedule.
	docker build -t ${IMG}:${IMAGE_TAG_SCHD} -f docker/schedule/Dockerfile .

.PHONY: schedule-push
schedule-push: ## Push docker image with the schedule.
	docker push ${IMG}:${IMAGE_TAG_SCHD}

##@ Deployment

.PHONY: helm
helm: helmify ## Generate Helm Charts
	kubectl kustomize install/default | $(HELMIFY) -vv helm/xuanwu-log
	kubectl kustomize install/default > helm/xuanwu-log/ALL.yaml

HELMIFY = $(shell pwd)/bin/helmify
.PHONY: helmify
helmify: ## Download helmify locally if necessary.
	$(call go-get-tool-kylin,$(HELMIFY),github.com/arttor/helmify/cmd/helmify@v0.3.4)

# go-get-tool-kylin will 'go get' any package $2 and install it to $1.
# this is used to get kyligence customized tool
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool-kylin
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "replace github.com/arttor/helmify => github.com/kyligence/helmify main" >> go.mod ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef