include env.mk

# You can use bigger machine type n1-highcpu-8 or n1-highcpu-32.
# See https://cloud.google.com/cloud-build/pricing
# for more detail.
ifdef CB_MACHINE_TYPE
	MACHINE_TYPE=--machine-type=$(CB_MACHINE_TYPE)
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
	_GIT_USER_ID=$(GIT_USER_ID) \
	_GIT_REPO_ID=$(GIT_REPO_ID)

BACKEND_TEST_SUBS = _GIT_USER_ID=$(GIT_USER_ID) \
	_GIT_REPO_ID=$(GIT_REPO_ID)

FRONTEND_E2E_SUBS = _DOMAIN=$(DOMAIN)

# cloudbuild.yaml
INFRA_SUBS = $(CLUSTER_ARGS) $(ISTIO_ARGS) \
	_BACKEND_GSA=$(BACKEND_GSA) \
	_BACKEND_KSA=$(BACKEND_KSA) \
	_BACKEND_SERVICE_HOST_NAME=$(BACKEND_SERVICE_HOST_NAME) \
	_DOMAIN=$(DOMAIN) \
	_MANAGED_ZONE_NAME=$(MANAGED_ZONE_NAME) \
	_SSL_CERT_NAME=$(SSL_CERT_NAME)

# cloudbuild-provision-cluster.yaml
PROVISION_SUBS = $(CLUSTER_ARGS) $(ISTIO_ARGS) \
	_CLUSTER_GKE_VERSION=$(CLUSTER_GKE_VERSION)

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
OPENAPI_GEN_CLIENT_ARGS=-g typescript-angular -i openapi.yaml -o webui/api-client

CLUSTER_MISSING=$(shell gcloud --project=$(PROJECT_ID) container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null; echo $$?)

.PHONY: clean delete run-local-webui run-local-backend lint-webui lint test-webui-local test-backend-local build-webui test-webui build-backend build-infrastructure build-all test cluster

## RULES FOR LOCAL DEVELOPMENT
clean:
	rm -rf webui/node_modules webui/api-client
	git clean -d -f -X backend/

/tmp/$(OPENAPI_GEN_JAR):
	wget $(OPENAPI_GEN_URL) -P /tmp/

webui/api-client: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_CLIENT_ARGS)

webui/node_modules:
	cd webui && npm ci

backend/api-service/src/api/openapi.yaml: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml $(CUSTOM_TEMPLATES)/*.mustache
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_SERVER_ARGS)

# Uses port 4200
run-local-webui: webui/api-client
	cd webui && ng serve --proxy-config proxy.conf.json

# Uses port 8080
run-local-backend: backend/api-service/src/api/openapi.yaml
	cd backend/api-service && go run main.go

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

test-webui-local: webui/api-client webui/node_modules
	cd webui && npm run test -- --watch=false --browsers=ChromeHeadless

test-webui-e2e-local: webui/api-client webui/node_modules
	cd webui && npm run e2e -- --dev-server-target= --base-url=http://localhost:4200 

## RULES FOR CLOUD DEVELOPMENT
GCLOUD_BUILD=gcloud --project=$(PROJECT_ID) builds submit $(MACHINE_TYPE) --verbosity=info .

cluster:
	if ! gcloud --project=$(PROJECT_ID) container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null; then \
	  echo creating cluster $(CLUSTER_NAME); \
	  $(GCLOUD_BUILD) --config cloudbuild-provision-cluster.yaml --substitutions $(call join_subs,$(PROVISION_SUBS)); \
	  gcloud --project=$(PROJECT_ID) container clusters get-credentials $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION); \
	fi

delete:
	$(GCLOUD_BUILD) --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=delete,$(call join_subs,$(INFRA_SUBS))

build-webui: cluster
	$(GCLOUD_BUILD) --config ./webui/cloudbuild.yaml --substitutions $(call join_subs,$(WEBUI_SUBS))

test-backend:
	$(GCLOUD_BUILD) --config ./backend/api-service/cloudbuild-test.yaml --substitutions $(call join_subs,$(BACKEND_TEST_SUBS))

test-webui:
	$(GCLOUD_BUILD) --config ./webui/cloudbuild-test.yaml .

test-webui-e2e:
	$(GCLOUD_BUILD) --config ./webui/e2e/cloudbuild.yaml --substitutions $(call join_subs,$(FRONTEND_E2E_SUBS))

build-backend: cluster
	$(GCLOUD_BUILD) --config ./backend/api-service/cloudbuild.yaml --substitutions $(call join_subs,$(BACKEND_SUBS))

build-infrastructure: cluster
	$(GCLOUD_BUILD) . --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=apply,$(call join_subs,$(INFRA_SUBS))

build-all: build-infrastructure build-webui build-backend

test: test-backend test-webui
