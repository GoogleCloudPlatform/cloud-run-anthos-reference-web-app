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


FROM node:12-alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY package*.json /
RUN npm ci --only=production
COPY dist /dist

# Run the web service on container startup.
CMD ["node", "/dist/index.js"]