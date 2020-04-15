# Configure all these variables for your project/application

PROJECT_ID:=spencersmall-knative-dev

# Parameters for code generation
GIT_USER_ID=GoogleCloudPlatform
GIT_REPO_ID=cloud-run-gke-reference-web-app

# The name of the GKE cluster to use.
# TODO: Replace with YOUR cluster's name
CLUSTER_NAME:=crfa-canonical-web-app

# Cluster information.
CLUSTER_LOCATION=$(shell gcloud container clusters list --filter="name=$(CLUSTER_NAME)" --format="csv[no-heading](location)" )
CLUSTER_GKE_VERSION=$(shell gcloud container clusters list --filter="name=$(CLUSTER_NAME)" --format="csv[no-heading](currentMasterVersion)")

# If you are provisioning a cluster, manually set these values instead
ifeq ($(CLUSTER_LOCATION),)
	# TODO: Specify the appropriate region for cluster creation
	CLUSTER_LOCATION=us-west1-a
endif
ifeq ($(CLUSTER_GKE_VERSION),)
	# TODO: Specify the appropriate GKE version for cluster creation
	CLUSTER_GKE_VERSION=1.15
endif

# Subdomain of cloud-tutorial.dev used for the demo
# TODO: Set your own subdomain name here
DOMAIN=spencersmall.run
MANAGED_ZONE_NAME=$(shell gcloud dns managed-zones list --format="csv[no-heading](name)" --filter="dnsName:$(DOMAIN)")

# Namespace to be used by app and KCC resources
NAMESPACE=web-app

# Istio Ingress information
ISTIO_INGRESS_SERVICE=istio-ingress
ISTIO_INGRESS_NAMESPACE=gke-system

# Backend service name
BACKEND_IMAGE_NAME=backserv
BACKEND_SERVICE_NAME=$(BACKEND_IMAGE_NAME)
BACKEND_SERVICE_HOST_NAME=$(BACKEND_SERVICE_NAME).$(NAMESPACE).example.com

SSL_CERT_NAME=spencersmall-run-ssl-certificate

# Workload Identity service account names
BACKEND_KSA=$(NAMESPACE)-$(BACKEND_SERVICE_NAME)
BACKEND_GSA=ksa-$(BACKEND_KSA)