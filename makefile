include env.mk

# You can use bigger machine type n1-highcpu-8 or n1-highcpu-32.
# See https://cloud.google.com/cloud-build/pricing
# for more detail.
ifdef CB_MACHINE_TYPE
	MACHINE_TYPE=--machine-type=$(CB_MACHINE_TYPE)
endif

ifeq ($(EVENTING_ENABLED),true)
WEBUI_E2E_TEST_TAGS="@core or @alerts"
else
WEBUI_E2E_TEST_TAGS="@core"
endif

# Shared cluster substitution args
CLUSTER_ARGS = \
	_CLUSTER_LOCATION=$(CLUSTER_LOCATION) \
	_CLUSTER_NAME=$(CLUSTER_NAME) \
	_NAMESPACE=$(NAMESPACE)

# Shared istio substitution args
ISTIO_ARGS = \
	_ISTIO_INGRESS_NAMESPACE=$(ISTIO_INGRESS_NAMESPACE) \
	_ISTIO_INGRESS_SERVICE=$(ISTIO_INGRESS_SERVICE)

# backend/cloudbuild.yaml
BACKEND_SUBS = $(CLUSTER_ARGS) \
	_BACKEND_IMAGE_NAME=$(BACKEND_IMAGE_NAME) \
	_BACKEND_KSA=$(BACKEND_KSA) \
	_BACKEND_SERVICE_NAME=$(BACKEND_SERVICE_NAME) \
	_EVENT_BROKER_HOSTNAME=$(EVENT_BROKER_HOSTNAME) \
	_EVENTING_ENABLED=$(EVENTING_ENABLED) \
	_GIT_USER_ID=$(GIT_USER_ID) \
	_GIT_REPO_ID=$(GIT_REPO_ID)

# backend/inventory-state-service/cloudbuild.yaml
INVENTORY_STATE_SERVICE_SUBS = $(CLUSTER_ARGS) \
	_INVENTORY_STATE_IMAGE_NAME=$(INVENTORY_STATE_IMAGE_NAME) \
	_INVENTORY_STATE_SERVICE_NAME=$(INVENTORY_STATE_SERVICE_NAME) \
	_BACKEND_CLUSTER_HOST_NAME=$(BACKEND_CLUSTER_HOST_NAME)

# backend/inventory-level-monitor-service/cloudbuild.yaml
INVENTORY_LEVEL_MONITOR_SERVICE_SUBS = $(CLUSTER_ARGS) \
	_INVENTORY_LEVEL_MONITOR_IMAGE_NAME=$(INVENTORY_LEVEL_MONITOR_IMAGE_NAME) \
	_INVENTORY_LEVEL_MONITOR_SERVICE_NAME=$(INVENTORY_LEVEL_MONITOR_SERVICE_NAME) \
	_BACKEND_CLUSTER_HOST_NAME=$(BACKEND_CLUSTER_HOST_NAME)

# backend/inventory-balance-monitor-service/cloudbuild.yaml
INVENTORY_BALANCE_MONITOR_SERVICE_SUBS = $(CLUSTER_ARGS) \
	_INVENTORY_BALANCE_MONITOR_IMAGE_NAME=$(INVENTORY_BALANCE_MONITOR_IMAGE_NAME) \
	_INVENTORY_BALANCE_MONITOR_SERVICE_NAME=$(INVENTORY_BALANCE_MONITOR_SERVICE_NAME) \
	_BACKEND_CLUSTER_HOST_NAME=$(BACKEND_CLUSTER_HOST_NAME)

BACKEND_TEST_SUBS = _GIT_USER_ID=$(GIT_USER_ID) \
	_GIT_REPO_ID=$(GIT_REPO_ID)

USER_SVC_SUBS = $(CLUSTER_ARGS) \
	_USER_SVC_IMAGE=$(USER_SVC_IMAGE_NAME) \
	_USER_SVC_KSA=$(USER_SVC_KSA) \
	_USER_SVC_NAME=$(USER_SVC_NAME) \

FRONTEND_E2E_SUBS = _DOMAIN=$(DOMAIN)

# cloudbuild.yaml
INFRA_SUBS = $(CLUSTER_ARGS) $(ISTIO_ARGS) \
	_BACKEND_GSA=$(BACKEND_GSA) \
	_BACKEND_KSA=$(BACKEND_KSA) \
	_BACKEND_SERVICE_HOST_NAME=$(BACKEND_SERVICE_HOST_NAME) \
	_USER_SVC_KSA=$(USER_SVC_KSA) \
	_USER_SVC_GSA=$(USER_SVC_GSA) \
	_USER_SVC_HOST_NAME=$(USER_SVC_HOST_NAME) \
	_DOMAIN=$(DOMAIN) \
	_MANAGED_ZONE_NAME=$(MANAGED_ZONE_NAME) \
	_SSL_CERT_NAME=$(SSL_CERT_NAME)

# cloudbuild-provision-cluster.yaml
PROVISION_SUBS = $(CLUSTER_ARGS) $(ISTIO_ARGS) \
	_CLUSTER_GKE_VERSION=$(CLUSTER_GKE_VERSION)

# cloudbuild-eventing-triggers.yaml
EVENTING_TRIGGERS_SUBS = $(CLUSTER_ARGS) \
	_INVENTORY_STATE_SERVICE_NAME=$(INVENTORY_STATE_SERVICE_NAME) \
	_INVENTORY_BALANCE_MONITOR_SERVICE_NAME=$(INVENTORY_BALANCE_MONITOR_SERVICE_NAME) \
	_INVENTORY_LEVEL_MONITOR_IMAGE_NAME=$(INVENTORY_LEVEL_MONITOR_IMAGE_NAME)

# webui/cloudbuild.yaml
WEBUI_SUBS = _DOMAIN=$(DOMAIN)

# Comma separate substitution args
comma := ,
empty :=
space := $(empty) $(empty)
join_subs = $(subst $(space),$(comma),$(1))

# Open API args
CUSTOM_TEMPLATES=backend/templates
OPENAPI_GEN_JAR=openapi-generator-cli-4.3.0.jar
OPENAPI_GEN_URL="https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/4.3.0/$(OPENAPI_GEN_JAR)"
OPENAPI_GEN_SERVER_ARGS=-g go-server -i openapi.yaml -o backend/api-service --api-name-suffix= --git-user-id=$(GIT_USER_ID) --git-repo-id=$(GIT_REPO_ID)/api-service --package-name=service -t $(CUSTOM_TEMPLATES) --additional-properties=sourceFolder=src
OPENAPI_GEN_TYPESCRIPT_CLIENT_ARGS=-g typescript-angular -i openapi.yaml -o webui/api-client
OPENAPI_GEN_TYPESCRIPT_USER_CLIENT_ARGS=-g typescript-angular -i backend/user-service/user-api.yaml -o webui/user-svc-client
OPENAPI_GEN_GO_CLIENT_ARGS=-g go -i openapi.yaml -o backend/api-client --package-name=client

CLUSTER_MISSING=$(shell gcloud --project=$(PROJECT_ID) container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null; echo $$?)

.PHONY: clean delete run-local-webui run-local-backend lint-webui lint test-webui-local test-backend-local build-webui test-webui build-backend build-infrastructure build-all test cluster
# Build eventing
.PHONY: build-eventing build-eventing-triggers run-local-inventory-state-service run-local-inventory-level-monitor-service run-local-inventory-balance-monitor-service build-inventory-state-service build-inventory-level-monitor-service build-inventory-balance-monitor-service
# Test eventing
.PHONY: test-eventing test-inventory-state-service-local test-inventory-level-monitor-service-local test-inventory-balance-monitor-service-local test-inventory-state-service test-inventory-level-monitor-service test-inventory-balance-monitor-service

## RULES FOR LOCAL DEVELOPMENT
clean:
	rm -rf webui/node_modules webui/api-client
	git clean -d -f -X backend/

/tmp/$(OPENAPI_GEN_JAR):
	wget $(OPENAPI_GEN_URL) -P /tmp/

webui/api-client: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_TYPESCRIPT_CLIENT_ARGS)

webui/user-svc-client: /tmp/$(OPENAPI_GEN_JAR) backend/user-service/user-api.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_TYPESCRIPT_USER_CLIENT_ARGS)

webui/node_modules:
	cd webui && npm ci

backend/api-service/src/api/openapi.yaml: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml $(CUSTOM_TEMPLATES)/*.mustache
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_SERVER_ARGS)

backend/api-client/openapi.yaml: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_GO_CLIENT_ARGS)

# Uses port 4200
run-local-webui: webui/api-client
	cd webui && ng serve --proxy-config proxy.conf.json

# Uses port 8080
run-local-backend: backend/api-service/src/api/openapi.yaml
	cd backend/api-service && go run main.go

run-local-inventory-state-service: backend/api-client/openapi.yaml
	cd backend/inventory-state-service && go run main.go -backend_cluster_host_name=localhost:8080

run-local-inventory-level-monitor-service: backend/api-client/openapi.yaml
	cd backend/inventory-level-monitor-service && go run main.go -backend_cluster_host_name=localhost:8080

run-local-inventory-balance-monitor-service: backend/api-client/openapi.yaml
	cd backend/inventory-balance-monitor-service && go run main.go -backend_cluster_host_name=localhost:8080

lint-webui: webui/node_modules
	cd webui && npm run lint

lint: lint-webui

test-backend-local: backend/api-service/src/api/openapi.yaml
	docker stop firestore-emulator 2>/dev/null || true
	docker run --detach --rm -p 9090:9090 --name=firestore-emulator google/cloud-sdk:292.0.0 sh -c \
	 "apt-get install google-cloud-sdk-firestore-emulator && gcloud beta emulators firestore start --host-port=0.0.0.0:9090"
	docker run --network=host jwilder/dockerize:0.6.1 dockerize -timeout=60s -wait=tcp://localhost:9090
	cd backend/api-service/src && FIRESTORE_EMULATOR_HOST=localhost:9090 go test -tags=emulator -v
	docker stop firestore-emulator

test-inventory-state-service-local: backend/api-client/openapi.yaml
	cd backend/inventory-state-service && go test -v ./...

test-inventory-level-monitor-service-local: backend/api-client/openapi.yaml
	cd backend/inventory-level-monitor-service/src && go test -v

test-inventory-balance-monitor-service-local: backend/api-client/openapi.yaml
	cd backend/inventory-balance-monitor-service/src && go test -v

test-webui-local: webui/api-client webui/node_modules
	cd webui && npm run test -- --watch=false --browsers=ChromeHeadless

test-webui-e2e-local: webui/api-client webui/node_modules
	cd webui && npm run e2e -e TAGS=$(WEBUI_E2E_TEST_TAGS)

test-webui-e2e-prod: webui/api-client webui/node_modules
	cd webui && npm run e2e -- --headless --config baseUrl=https://${DOMAIN} -e TAGS=$(WEBUI_E2E_TEST_TAGS)

## RULES FOR CLOUD DEVELOPMENT
GCLOUD_BUILD=gcloud --project=$(PROJECT_ID) builds submit $(MACHINE_TYPE) --verbosity=info .

cluster:
	if ! gcloud --project=$(PROJECT_ID) container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null; then \
	  echo creating cluster $(CLUSTER_NAME); \
	  $(GCLOUD_BUILD) --config cloudbuild-provision-cluster.yaml --substitutions $(call join_subs,$(PROVISION_SUBS)) && \
	  gcloud --project=$(PROJECT_ID) container clusters get-credentials $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION); \
	fi

delete:
	$(GCLOUD_BUILD) --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=delete,$(call join_subs,$(INFRA_SUBS))

build-webui: cluster
	$(GCLOUD_BUILD) --config ./webui/cloudbuild.yaml --substitutions $(call join_subs,$(WEBUI_SUBS))

test-backend:
	$(GCLOUD_BUILD) --config ./backend/api-service/cloudbuild-test.yaml --substitutions $(call join_subs,$(BACKEND_TEST_SUBS))

test-inventory-state-service:
	$(GCLOUD_BUILD) --config ./backend/inventory-state-service/cloudbuild-test.yaml

test-inventory-level-monitor-service:
	$(GCLOUD_BUILD) --config ./backend/inventory-level-monitor-service/cloudbuild-test.yaml

test-inventory-balance-monitor-service:
	$(GCLOUD_BUILD) --config ./backend/inventory-balance-monitor-service/cloudbuild-test.yaml

test-webui:
	$(GCLOUD_BUILD) --config ./webui/cloudbuild-test.yaml

test-webui-e2e:
	$(GCLOUD_BUILD) --config ./webui/cypress/cloudbuild.yaml --substitutions $(call join_subs,$(FRONTEND_E2E_SUBS)),_WEBUI_E2E_TEST_TAGS=$(WEBUI_E2E_TEST_TAGS)

build-backend: cluster
	$(GCLOUD_BUILD) --config ./backend/api-service/cloudbuild.yaml --substitutions $(call join_subs,$(BACKEND_SUBS))

build-userservice: cluster
	$(GCLOUD_BUILD) --config ./backend/user-service/cloudbuild.yaml --substitutions $(call join_subs,$(USER_SVC_SUBS))

build-inventory-state-service: cluster
	$(GCLOUD_BUILD) --config ./backend/inventory-state-service/cloudbuild.yaml --substitutions $(call join_subs,$(INVENTORY_STATE_SERVICE_SUBS))

build-inventory-level-monitor-service: cluster
	$(GCLOUD_BUILD) --config ./backend/inventory-level-monitor-service/cloudbuild.yaml --substitutions $(call join_subs,$(INVENTORY_LEVEL_MONITOR_SERVICE_SUBS))
  
build-inventory-balance-monitor-service: cluster
	$(GCLOUD_BUILD) --config ./backend/inventory-balance-monitor-service/cloudbuild.yaml --substitutions $(call join_subs,$(INVENTORY_BALANCE_MONITOR_SERVICE_SUBS))
  
build-eventing-triggers:
	$(GCLOUD_BUILD) --config cloudbuild-eventing-triggers.yaml --substitutions $(call join_subs,$(EVENTING_TRIGGERS_SUBS))

ifeq ($(EVENTING_ENABLED),true)
build-eventing: build-inventory-state-service build-inventory-level-monitor-service build-inventory-balance-monitor-service
	# Eventing triggers have to be set up after the eventing services which they trigger.
	$(MAKE) build-eventing-triggers
else
build-eventing:
	@echo Eventing is disabled. To enable eventing set the EVENTING_ENABLED flag to true in env.mk.
endif

ifeq ($(EVENTING_ENABLED),true)
test-eventing: test-inventory-state-service test-inventory-level-monitor-service test-inventory-balance-monitor-service
else
test-eventing:
	@echo Eventing is disabled. To test eventing set the EVENTING_ENABLED flag to true in env.mk.
endif

build-infrastructure: cluster
	$(GCLOUD_BUILD) --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=apply,$(call join_subs,$(INFRA_SUBS))

build-infra: build-infrastructure

build-all: build-infrastructure build-backend build-userservice build-webui build-eventing

test: test-backend test-webui test-eventing
