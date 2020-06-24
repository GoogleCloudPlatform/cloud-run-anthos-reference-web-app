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


readonly BASE_PATH=$(dirname "${BASH_SOURCE[0]}")/..
readonly ADMIN_PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
readonly WORKDER_PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

cat ${BASE_PATH}/credentials.template.js | \
sed "s/\${ADMIN_PASSWORD}/${ADMIN_PASSWORD}/g" | \
sed "s/\${WORKDER_PASSWORD}/${WORKDER_PASSWORD}/g" > ${BASE_PATH}/credentials.js

echo "Test user passwords generated."
