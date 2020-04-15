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



# PROJECT_ID
PROJECT_ID=$1
PROJECT_NUMBER=$(gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")

# Enable cloud run api
gcloud services enable cloudbuild.googleapis.com

# Enable kubernetes engine api
gcloud services enable container.googleapis.com

# Enable cloud resource manager api
gcloud services enable cloudresourcemanager.googleapis.com

# Grant cloud build service account permissions

# Service Account User, roles/iam.serviceAccountUser
# Kubernetes Engine Admin, roles/container.admin
# Project IAM Admin, roles/resourcemanager.projectIamAdmin
# Service Account Admin, roles/iam.serviceAccountAdmin

for role in roles/container.developer roles/iam.serviceAccountUser roles/container.admin roles/resourcemanager.projectIamAdmin roles/iam.serviceAccountAdmin roles/compute.loadBalancerAdmin roles/compute.networkAdmin
do
  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
  --role ${role}
done

