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

# webui/cloudbuild.yaml
WEBUI_SUBS = _DOMAIN=$(DOMAIN)

ISTIO_AUTH_TEST_SUBS = $(ISTIO_ARGS) \
	_CLUSTER_LOCATION=$(CLUSTER_LOCATION) \
	_CLUSTER_NAME=$(CLUSTER_NAME)

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
OPENAPI_GEN_API_CLIENT_ARGS=-g typescript-angular -i openapi.yaml -o webui/api-client
OPENAPI_GEN_USER_CLIENT_ARGS=-g typescript-angular -i backend/user-service/user-api.yaml -o webui/user-svc-client

CLUSTER_MISSING=$(shell gcloud --project=$(PROJECT_ID) container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null; echo $$?)

.PHONY: clean delete delete-cluster run-local-webui run-local-backend lint-webui lint test-webui-local test-backend-local test-istio-auth-local build-webui test-webui test-istio-auth build-backend build-infrastructure build-all test cluster jq

## RULES FOR LOCAL DEVELOPMENT
clean:
	rm -rf webui/node_modules webui/api-client
	git clean -d -f -X backend/

/tmp/$(OPENAPI_GEN_JAR):
	wget $(OPENAPI_GEN_URL) -P /tmp/

webui/api-client: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_API_CLIENT_ARGS)

webui/user-svc-client: /tmp/$(OPENAPI_GEN_JAR) backend/user-service/user-api.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_USER_CLIENT_ARGS)

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

jq:
	@which jq > /dev/null || (echo "'jq' needs to be installed for this target to run. It can be downloaded from https://stedolan.github.io/jq/." && exit 1)

test-backend-local: backend/api-service/src/api/openapi.yaml
	docker stop firestore-emulator 2>/dev/null || true
	docker run --detach --rm -p 9090:9090 --name=firestore-emulator google/cloud-sdk:292.0.0 sh -c \
	 "apt-get install google-cloud-sdk-firestore-emulator && gcloud beta emulators firestore start --host-port=0.0.0.0:9090"
	docker run --network=host jwilder/dockerize:0.6.1 dockerize -timeout=60s -wait=tcp://localhost:9090
	cd backend/api-service/src && FIRESTORE_EMULATOR_HOST=localhost:9090 go test -tags=emulator -v
	docker stop firestore-emulator

FIREBASE_SA=$(shell gcloud --project=$(PROJECT_ID) iam service-accounts list --filter="displayName=firebase-adminsdk" --format="value(email)")
test-istio-auth-local: jq
	gcloud --project=$(PROJECT_ID) iam service-accounts keys create --iam-account=$(FIREBASE_SA) \
		/tmp/istio-auth-test-key.json
	cd istio-auth && API_KEY=$$(grep apiKey ../webui/firebaseConfig.ts | cut -d "'" -f2) \
		HOST_IP=$$(kubectl -n $(ISTIO_INGRESS_NAMESPACE) get service $(ISTIO_INGRESS_SERVICE) -o jsonpath='{.status.loadBalancer.ingress[0].ip}') \
		GOOGLE_APPLICATION_CREDENTIALS=/tmp/istio-auth-test-key.json \
		go test -v || touch /tmp/istio-auth-test.failed
	gcloud --project=$(PROJECT_ID) -q iam service-accounts keys delete --iam-account=$(FIREBASE_SA) \
		$$(jq -r .private_key_id /tmp/istio-auth-test-key.json)
	rm /tmp/istio-auth-test-key.json
	! rm /tmp/istio-auth-test.failed 2>/dev/null

test-webui-local: webui/api-client webui/node_modules
	cd webui && npm run test -- --watch=false --browsers=ChromeHeadless

test-webui-e2e-local: webui/api-client webui/node_modules
	cd webui && npm run e2e

test-webui-e2e-prod: webui/api-client webui/node_modules
	cd webui && npm run e2e -- --headless --config baseUrl=https://${DOMAIN}

## RULES FOR CLOUD DEVELOPMENT
GCLOUD_BUILD=gcloud --project=$(PROJECT_ID) builds submit $(MACHINE_TYPE) --verbosity=info .

cluster:
	if ! gcloud --project=$(PROJECT_ID) container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null; then \
	  echo creating cluster $(CLUSTER_NAME); \
	  $(GCLOUD_BUILD) --config cloudbuild-provision-cluster.yaml --substitutions $(call join_subs,$(PROVISION_SUBS)) && \
	  gcloud --project=$(PROJECT_ID) container clusters get-credentials $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION); \
	fi

delete-cluster:
	gcloud --project=$(PROJECT_ID) container clusters delete $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) --quiet

delete:
	$(GCLOUD_BUILD) --config cloudbuild.yaml --timeout=2400 --substitutions _APPLY_OR_DELETE=delete,$(call join_subs,$(INFRA_SUBS))

build-webui: cluster
	$(GCLOUD_BUILD) --config ./webui/cloudbuild.yaml --substitutions $(call join_subs,$(WEBUI_SUBS))

test-backend:
	$(GCLOUD_BUILD) --config ./backend/api-service/cloudbuild-test.yaml --substitutions $(call join_subs,$(BACKEND_TEST_SUBS))

test-istio-auth:
	$(GCLOUD_BUILD) --config ./istio-auth/cloudbuild-test.yaml --substitutions $(call join_subs,$(ISTIO_AUTH_TEST_SUBS))

test-webui:
	$(GCLOUD_BUILD) --config ./webui/cloudbuild-test.yaml

test-webui-e2e:
	$(GCLOUD_BUILD) --config ./webui/cypress/cloudbuild.yaml --substitutions $(call join_subs,$(FRONTEND_E2E_SUBS))

build-backend: cluster
	$(GCLOUD_BUILD) --config ./backend/api-service/cloudbuild.yaml --substitutions $(call join_subs,$(BACKEND_SUBS))

build-userservice: cluster
	$(GCLOUD_BUILD) --config ./backend/user-service/cloudbuild.yaml --substitutions $(call join_subs,$(USER_SVC_SUBS))

build-infrastructure: cluster
	$(GCLOUD_BUILD) --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=apply,$(call join_subs,$(INFRA_SUBS))

build-infra: build-infrastructure

build-all: build-infrastructure build-backend build-userservice build-webui

test: test-backend test-webui
