#!/bin/bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

usage() {
  local name
  name=$(basename "$0")
  echo "Usage: ${name} PROJECT_ID"
  exit 1
}

if [[ "$#" -ne 1 ]]; then
  usage
fi

readonly PROJECT_ID="$1"
readonly PROJECT_NUMBER=$(gcloud projects describe "${PROJECT_ID}" --format="value(projectNumber)")

# Enable Cloud Build, Kubernetes Engine, and Cloud Resource Manager APIs
gcloud --project "${PROJECT_ID}" services enable {cloudbuild,container,cloudresourcemanager}.googleapis.com

# Grant Cloud Build service account permissions
# Service Account Admin, roles/iam.serviceAccountAdmin
# Service Account User, roles/iam.serviceAccountUser
# Kubernetes Engine Admin, roles/container.admin
# Project IAM Admin, roles/resourcemanager.projectIamAdmin
# Compute Load Balancer Admin, roles/compute.loadBalancerAdmin
# Compute Network Admin, roles/compute.networkAdmin

for role in iam.serviceAccount{Admin,User} container.admin resourcemanager.projectIamAdmin compute.{loadBalancer,network}Admin; do
  gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member serviceAccount:"${PROJECT_NUMBER}"@cloudbuild.gserviceaccount.com \
  --role roles/"${role}"
done

# Create env.mk if not present
if [[ ! -f "env.mk" ]]; then
  cp env.mk.sample env.mk
fi

# Substitute default PROJECT_ID value
sed "s/^PROJECT_ID=project-id$/PROJECT_ID=${PROJECT_ID}/" -i env.mk
