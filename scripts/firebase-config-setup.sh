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

# TODO(https://github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/issues/28):
# Fetch the API key once this command becomes generally available:
# https://cloud.google.com/sdk/gcloud/reference/alpha/services/api-keys/list

usage() {
  local name
  name=$(basename "$0")
  echo "Usage: ${name} PROJECT_ID API_KEY"
  echo
  echo "API_KEY is the Web API Key found at: https://console.firebase.google.com/project/${PROJECT_ID:-YOUR_PROJECT_ID}/settings/general"
  exit 1
}

export readonly PROJECT_ID="$1"
readonly APP_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
readonly FIREBASECONFIG_SAMPLE="${APP_ROOT}/webui/firebaseConfig.sample.ts"
readonly FIREBASECONFIG_JS="${APP_ROOT}/webui/firebaseConfig.ts"

if [[ "$#" -ne 2 ]]; then
  usage
fi

export readonly API_KEY="$2"
echo ${API_KEY}

cat ${FIREBASECONFIG_SAMPLE}
cat ${FIREBASECONFIG_SAMPLE} | envsubst > ${FIREBASECONFIG_JS}

echo
echo "Wrote to $(realpath ${FIREBASECONFIG_JS}):"
echo
cat "${FIREBASECONFIG_JS}"
