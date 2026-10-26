# Image used by deploy and build-installer. The default matches what `make image`
# builds locally and what config/manager/kustomization.yaml points at.
IMG ?= ko.local/respawn:dev
# YEAR defines the year value used for substituting the YEAR placeholder in the boilerplate header.
YEAR ?= $(shell date +%Y)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'.

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate ClusterRole and CustomResourceDefinition objects.
	"$(CONTROLLER_GEN)" rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	"$(CONTROLLER_GEN)" object paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: verify
verify: manifests generate ## Fail if generated code, manifests, or go.mod are out of date.
	go mod tidy -diff
	@test -z "$$(git status --porcelain -- api config)" || { \
		git status --porcelain -- api config; \
		echo "Generated files are out of date. Run 'make generate manifests' and commit the result."; \
		exit 1; \
	}

.PHONY: test
test: manifests generate fmt vet setup-envtest ## Run unit and envtest integration tests with the race detector.
	KUBEBUILDER_ASSETS="$(shell "$(ENVTEST)" use $(ENVTEST_K8S_VERSION) --bin-dir "$(LOCALBIN)" -p path)" \
		go test -race $$(go list ./... | grep -v /e2e) -coverprofile cover.out
	@grep -E '^mode:|/internal/' cover.out > cover.internal.out
	@echo "internal/ coverage: $$(go tool cover -func cover.internal.out | tail -n 1 | awk '{print $$NF}')"

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	"$(GOLANGCI_LINT)" run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	"$(GOLANGCI_LINT)" run --fix

.PHONY: lint-config
lint-config: golangci-lint ## Verify golangci-lint linter configuration
	"$(GOLANGCI_LINT)" config verify

##@ Local cluster (k3d)

# k3s image for local and e2e clusters. Keep in step with the newest supported k3s minor.
K3S_IMAGE ?= rancher/k3s:v1.37.1-k3s1
K3D_CLUSTER ?= respawn-dev
K3D_E2E_CLUSTER ?= respawn-e2e

.PHONY: dev-cluster
dev-cluster: k3d ## Create the local k3d cluster if it does not exist.
	@"$(K3D)" cluster list $(K3D_CLUSTER) >/dev/null 2>&1 && echo "k3d cluster '$(K3D_CLUSTER)' already exists" || \
		"$(K3D)" cluster create $(K3D_CLUSTER) --image $(K3S_IMAGE) --k3s-arg "--disable=traefik@server:0" --wait

.PHONY: dev-cluster-delete
dev-cluster-delete: k3d ## Delete the local k3d cluster.
	"$(K3D)" cluster delete $(K3D_CLUSTER)

.PHONY: dev-deploy
dev-deploy: image k3d ## Build the operator image, import it into the local k3d cluster, and deploy.
	"$(K3D)" image import $(IMG) --cluster $(K3D_CLUSTER)
	$(MAKE) deploy IMG=$(IMG)
	"$(KUBECTL)" -n respawn-system rollout restart deployment/respawn-controller-manager
	"$(KUBECTL)" -n respawn-system rollout status deployment/respawn-controller-manager --timeout=2m

.PHONY: setup-test-e2e
setup-test-e2e: k3d ## Create the k3d cluster for e2e tests if it does not exist.
	@"$(K3D)" cluster list $(K3D_E2E_CLUSTER) >/dev/null 2>&1 && echo "k3d cluster '$(K3D_E2E_CLUSTER)' already exists" || \
		"$(K3D)" cluster create $(K3D_E2E_CLUSTER) --image $(K3S_IMAGE) --k3s-arg "--disable=traefik@server:0" --wait

# kubectl kuberc is disabled by default for test isolation; enable with KUBECTL_KUBERC=true.
.PHONY: test-e2e
test-e2e: setup-test-e2e manifests generate fmt vet ## Run the e2e tests against a dedicated k3d cluster.
	"$(K3D)" kubeconfig merge $(K3D_E2E_CLUSTER) --kubeconfig-switch-context >/dev/null
	K3D="$(K3D)" K3D_CLUSTER=$(K3D_E2E_CLUSTER) go test -tags=e2e ./test/e2e/ -v -ginkgo.v -timeout 20m
	$(MAKE) cleanup-test-e2e

.PHONY: cleanup-test-e2e
cleanup-test-e2e: k3d ## Tear down the k3d cluster used for e2e tests.
	@"$(K3D)" cluster delete $(K3D_E2E_CLUSTER)

##@ Build

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager ./cmd

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host against the current kubeconfig.
	go run ./cmd

# Platforms for release images. Raspberry Pi (arm64) clusters are a core audience.
PLATFORMS ?= linux/amd64,linux/arm64

.PHONY: image
image: ko ## Build the manager image for this machine's architecture into the local Docker daemon as $(IMG).
	KO_DOCKER_REPO=$(firstword $(subst :, ,$(IMG))) "$(KO)" build ./cmd --bare --tags $(lastword $(subst :, ,$(IMG))) \
		--platform linux/$(shell go env GOARCH)

.PHONY: image-all
image-all: ko ## Build the manager image for all $(PLATFORMS) without pushing it.
	KO_DOCKER_REPO=build-only.invalid/respawn "$(KO)" build ./cmd --bare --platform $(PLATFORMS) --push=false

.PHONY: build-installer
build-installer: manifests generate kustomize ## Generate a consolidated YAML with CRDs and deployment.
	mkdir -p dist
	cd config/manager && "$(KUSTOMIZE)" edit set image controller=${IMG}
	"$(KUSTOMIZE)" build config/default > dist/install.yaml

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	@out="$$( "$(KUSTOMIZE)" build config/crd 2>/dev/null || true )"; \
	if [ -n "$$out" ]; then echo "$$out" | "$(KUBECTL)" apply -f -; else echo "No CRDs to install; skipping."; fi

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	@out="$$( "$(KUSTOMIZE)" build config/crd 2>/dev/null || true )"; \
	if [ -n "$$out" ]; then echo "$$out" | "$(KUBECTL)" delete --ignore-not-found=$(ignore-not-found) -f -; else echo "No CRDs to delete; skipping."; fi

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && "$(KUSTOMIZE)" edit set image controller=${IMG}
	"$(KUSTOMIZE)" build config/default | "$(KUBECTL)" apply -f -

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	"$(KUSTOMIZE)" build config/default | "$(KUBECTL)" delete --ignore-not-found=$(ignore-not-found) -f -

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
KO ?= $(LOCALBIN)/ko
K3D ?= $(LOCALBIN)/k3d

## Tool Versions
KUSTOMIZE_VERSION ?= v5.8.1
CONTROLLER_TOOLS_VERSION ?= v0.23.0
# Keep in step with the version in .custom-gcl.yml.
GOLANGCI_LINT_VERSION ?= v2.15.0
KO_VERSION ?= v0.19.1
K3D_VERSION ?= v5.9.0

#ENVTEST_VERSION is the controller-runtime version to use for setup-envtest, derived from go.mod
ENVTEST_VERSION ?= $(shell v='$(call gomodver,sigs.k8s.io/controller-runtime)'; \
  [ -n "$$v" ] || { echo "Set ENVTEST_VERSION manually (controller-runtime replace has no tag)" >&2; exit 1; }; \
  printf '%s\n' "$$v")

#ENVTEST_K8S_VERSION is the version of Kubernetes to use for setting up ENVTEST binaries (i.e. 1.31)
ENVTEST_K8S_VERSION ?= $(shell v='$(call gomodver,k8s.io/api)'; \
  [ -n "$$v" ] || { echo "Set ENVTEST_K8S_VERSION manually (k8s.io/api replace has no tag)" >&2; exit 1; }; \
  printf '%s\n' "$$v" | sed -E 's/^v?[0-9]+\.([0-9]+).*/1.\1/')

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: setup-envtest
setup-envtest: envtest ## Download the binaries required for ENVTEST in the local bin directory.
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@"$(ENVTEST)" use $(ENVTEST_K8S_VERSION) --bin-dir "$(LOCALBIN)" -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))
	@test -f .custom-gcl.yml && { \
		echo "Building custom golangci-lint with plugins..." && \
		$(GOLANGCI_LINT) custom --destination $(LOCALBIN) --name golangci-lint-custom && \
		mv -f $(LOCALBIN)/golangci-lint-custom $(GOLANGCI_LINT); \
	} || true

.PHONY: ko
ko: $(KO) ## Download ko locally if necessary.
$(KO): $(LOCALBIN)
	$(call go-install-tool,$(KO),github.com/google/ko,$(KO_VERSION))

.PHONY: k3d
k3d: $(K3D) ## Download k3d locally if necessary.
$(K3D): $(LOCALBIN)
	$(call go-install-tool,$(K3D),github.com/k3d-io/k3d/v5,$(K3D_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] && [ "$$(readlink -- "$(1)" 2>/dev/null)" = "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f "$(1)" ;\
GOBIN="$(LOCALBIN)" go install $${package} ;\
mv "$(LOCALBIN)/$$(basename "$(1)")" "$(1)-$(3)" ;\
} ;\
ln -sf "$$(realpath "$(1)-$(3)")" "$(1)"
endef

define gomodver
$(shell go list -m -f '{{if .Replace}}{{.Replace.Version}}{{else}}{{.Version}}{{end}}' $(1) 2>/dev/null)
endef
