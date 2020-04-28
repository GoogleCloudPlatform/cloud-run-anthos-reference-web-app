include env.mk

# Shared cluster substitution args
CLUSTER_ARGS := \
	_CLUSTER_LOCATION=$(CLUSTER_LOCATION) \
	_CLUSTER_NAME=$(CLUSTER_NAME) \
	_NAMESPACE=$(NAMESPACE)

# Shared istio substitution args
ISTIO_ARGS := \
	_ISTIO_INGRESS_NAMESPACE=$(ISTIO_INGRESS_NAMESPACE) \
	_ISTIO_INGRESS_SERVICE=$(ISTIO_INGRESS_SERVICE)

# backend/cloudbuild.yaml
BACKEND_SUBS := $(CLUSTER_ARGS) \
	_BACKEND_IMAGE_NAME=$(BACKEND_IMAGE_NAME) \
	_BACKEND_KSA=$(BACKEND_KSA) \
	_BACKEND_SERVICE_NAME=$(BACKEND_SERVICE_NAME) \
	_GIT_USER_ID=$(GIT_USER_ID) \
	_GIT_REPO_ID=$(GIT_REPO_ID)

BACKEND_TEST_SUBS := _GIT_USER_ID=$(GIT_USER_ID) \
	_GIT_REPO_ID=$(GIT_REPO_ID)

FRONTEND_E2E_SUBS := _DOMAIN=$(DOMAIN)

# cloudbuild.yaml
INFRA_SUBS := $(CLUSTER_ARGS) $(ISTIO_ARGS) \
	_BACKEND_GSA=$(BACKEND_GSA) \
	_BACKEND_KSA=$(BACKEND_KSA) \
	_BACKEND_SERVICE_HOST_NAME=$(BACKEND_SERVICE_HOST_NAME) \
	_DOMAIN=$(DOMAIN) \
	_MANAGED_ZONE_NAME=$(MANAGED_ZONE_NAME) \
	_SSL_CERT_NAME=$(SSL_CERT_NAME)

# cloudbuild-provision-cluster.yaml
PROVISION_SUBS := $(CLUSTER_ARGS) $(ISTIO_ARGS) \
	_CLUSTER_GKE_VERSION=$(CLUSTER_GKE_VERSION)

# webui/cloudbuild.yaml
WEBUI_SUBS := _DOMAIN=$(DOMAIN)

# Comma separate substitution args
comma := ,
empty :=
space := $(empty) $(empty)
BACKEND_SUBS := $(subst $(space),$(comma),$(BACKEND_SUBS))
BACKEND_TEST_SUBS := $(subst $(space),$(comma),$(BACKEND_TEST_SUBS))
INFRA_SUBS := $(subst $(space),$(comma),$(INFRA_SUBS))
PROVISION_SUBS := $(subst $(space),$(comma),$(PROVISION_SUBS))
WEBUI_SUBS := $(subst $(space),$(comma),$(WEBUI_SUBS))

# Open API args
CUSTOM_TEMPLATES=backend/templates
OPENAPI_GEN_JAR=openapi-generator-cli-4.3.0.jar
OPENAPI_GEN_URL="https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/4.3.0/$(OPENAPI_GEN_JAR)"
OPENAPI_GEN_SERVER_ARGS=-g go-server -i openapi.yaml -o backend/src --api-name-suffix= --git-user-id=$(GIT_USER_ID) --git-repo-id=$(GIT_REPO_ID) --package-name=service -t $(CUSTOM_TEMPLATES)
OPENAPI_GEN_CLIENT_ARGS=-g typescript-angular -i openapi.yaml -o webui/api-client

IGNORE := $(shell gcloud container clusters describe $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION) 2>&1 > /dev/null)
CLUSTER_MISSING=$(.SHELLSTATUS)

.PHONY: bootstrap clean delete run-local-webui run-local-backend check-license lint-webui lint test-webui-local test-backend-local build-webui test-webui build-backend build-infrastructure build-all test cluster

## RULES FOR LOCAL DEVELOPMENT
bootstrap:
	./bootstrap.sh $(PROJECT_ID)

clean:
	rm -rf webui/node_modules webui/api-client
	git clean -d -f -X backend/

/tmp/$(OPENAPI_GEN_JAR):
	wget $(OPENAPI_GEN_URL) -P /tmp/

${GOPATH}/bin/addlicense:
	@echo addlicense is not installed, installing it now ...
	go get -u github.com/google/addlicense
	@echo ... installation complete

webui/api-client: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_CLIENT_ARGS)

webui/node_modules:
	cd webui && npm ci

backend/src/api/openapi.yaml: /tmp/$(OPENAPI_GEN_JAR) openapi.yaml $(CUSTOM_TEMPLATES)/*.mustache
	java -jar /tmp/$(OPENAPI_GEN_JAR) generate $(OPENAPI_GEN_SERVER_ARGS)

# Uses port 4200
run-local-webui: webui/api-client
	cd webui && ng serve --proxy-config proxy.conf.json

# Uses port 8080
run-local-backend: backend/src/api/openapi.yaml
	cd backend/src && go run main.go

check-license: ${GOPATH}/bin/addlicense
	${GOPATH}/bin/addlicense -check .

lint-webui: webui/node_modules
	cd webui && npm run lint

lint: check-license lint-webui

test-backend-local: backend/src/api/openapi.yaml
	docker stop firestore-emulator 2>/dev/null || true
	docker run --detach --rm -p 9090:9090 --name=firestore-emulator jdlk7/firestore-emulator
	docker run --network=host jwilder/dockerize:0.6.1 dockerize -timeout=60s -wait=tcp://localhost:9090
	cd backend/src/go && FIRESTORE_EMULATOR_HOST=localhost:9090 go test -tags=emulator -v
	docker stop firestore-emulator

test-webui-local: webui/api-client webui/node_modules
	cd webui && npm run test

test-webui-e2e-local: webui/api-client webui/node_modules
	cd webui && npm run e2e -- --dev-server-target= --base-url=http://localhost:4200 

## RULES FOR CLOUD DEVELOPMENT

cluster:
ifneq ($(CLUSTER_MISSING),0)
	@echo Cluster $(CLUSTER_NAME) does not exist, creating cluster
	gcloud builds submit . --verbosity=info --config cloudbuild-provision-cluster.yaml --substitutions $(PROVISION_SUBS)
	gcloud container clusters get-credentials $(CLUSTER_NAME) --zone $(CLUSTER_LOCATION)
endif
	# Restrict cluster ingress firewall rules to just GCLB's published ranges:
	# https://cloud.google.com/load-balancing/docs/https#source_ip_addresses
	for firewall in $$(gcloud --project=$(PROJECT_ID) compute firewall-rules list \
		  --filter="targetTags.list()~^gke-$(CLUSTER_NAME)-[0-9a-z]*-node$$ AND sourceRanges.list()=0.0.0.0/0" \
		  --format="csv[no-heading](name)"); do \
		gcloud --project=$(PROJECT_ID) compute firewall-rules update "$$firewall" \
		  --source-ranges=35.191.0.0/16,130.211.0.0/22; \
	done

delete:
	gcloud builds submit . --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=delete,$(INFRA_SUBS)

build-webui: cluster
	gcloud builds submit --config ./webui/cloudbuild.yaml --substitutions $(WEBUI_SUBS) .

test-backend:
	gcloud builds submit --config ./backend/cloudbuild-test.yaml --substitutions $(BACKEND_TEST_SUBS)  .

test-webui:
	gcloud builds submit --config ./webui/cloudbuild-test.yaml .

test-webui-e2e:
	gcloud builds submit --config ./webui/e2e/cloudbuild.yaml --substitutions $(FRONTEND_E2E_SUBS) .

build-backend: cluster
	gcloud builds submit --config ./backend/cloudbuild.yaml --substitutions $(BACKEND_SUBS) .

build-infrastructure: cluster
	gcloud builds submit . --config cloudbuild.yaml --substitutions _APPLY_OR_DELETE=apply,$(INFRA_SUBS)

build-all: build-infrastructure build-webui build-backend

test: test-backend test-webui
