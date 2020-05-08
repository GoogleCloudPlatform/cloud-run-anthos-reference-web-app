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

readonly APP_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

echo "This script will interactively walk you through setting up your custom domain"
echo "for use by the Cloud Run on Anthos Reference Web App."
echo

# Read the project id to use
readonly DEFAULT_PROJECT_ID=$(gcloud config get-value project)
read -r -p "Enter project id (default: ${DEFAULT_PROJECT_ID}): " PROJECT_ID
readonly PROJECT_ID="${PROJECT_ID:-$DEFAULT_PROJECT_ID}"
echo
echo "Using project: [${PROJECT_ID}]."
echo

# Ask user to choose whether to create or provide a custom domain
echo "A custom domain is required for this application to work properly."
echo "Choose an action to take:"
echo " [1] Create a custom subdomain"
echo " [2] Use existing domain/subdomain"
read -r -p "Please enter your numeric choice (default: 1): " choice
choice="${choice:-1}"
while [[ -n "${choice}" ]] && (("${choice}" != 1)) && (("${choice}" != 2)); do
  read -r -p "Please enter a valid numeric choice: " choice
  choice="${choice:-1}"
done
echo

# Create or read a custom domain
if (("${choice}" == 1)); then
  # Based on https://cloud-tutorial.dev/
  echo "This step will claim ownership of a subdomain of cloud-tutorial.dev."
  read -r -p "Enter desired zone name (e.g. my-cool-zone): " zone
  wget -P /tmp https://cloud-tutorial.dev/claim.sh
  chmod +x /tmp/claim.sh
  gcloud config set project "${PROJECT_ID}"
  DOMAIN="${zone}.cloud-tutorial.dev"
  if /tmp/claim.sh "${zone}"; then
    gcloud config set project "${DEFAULT_PROJECT_ID}"
    echo
    echo "The domain ${zone}.cloud-tutorial.dev has been successfully claimed."
  else
    gcloud config set project "${DEFAULT_PROJECT_ID}"
    echo
    echo "Unable to claim '${DOMAIN}'. Please try again with a different zone."
    exit 1
  fi
else
  read -r -p "Enter your existing domain (e.g. example.com): " DOMAIN
  echo
  echo "Please ensure:"
  echo " - Your custom domain is associated with a Cloud DNS Managed Zone in ${PROJECT_ID}."
  echo "   Additional info: https://cloud.google.com/dns/zones"
  echo " - Your name server records point to your managed zone."
  echo "   Additional info: https://cloud.google.com/dns/docs/migrating#update_your_registrars_name_server_records"
  echo
  read -r -n 1 -p "Once ready, press any key to continue."
fi
readonly DOMAIN
echo
echo "Using domain: [${DOMAIN}]."
echo

# Ensure a corresponding managed zone exists
zone=$(gcloud --project "${PROJECT_ID}" dns managed-zones list --format="csv[no-heading](name)" --filter="dnsName:${DOMAIN}")
if [[ -z "${zone}" ]]; then
  echo "No Cloud DNS Managed Zone corresponding to ${DOMAIN} found in ${PROJECT_ID}."
  echo "Please ensure your custom domain is associated with a Cloud DNS Managed Zone"
  echo "in ${PROJECT_ID} before trying again. For additional info, visit:"
  echo "https://cloud.google.com/dns/zones"
  exit 1
fi

# Verify the domain
url="https://search.google.com/search-console?resource_id=sc-domain:${DOMAIN}"
echo "In your web browser, open:"
echo "${url}"
echo
read -r -p "Does the page say you don't have access to the property [y/n]? " answer
while [[ "${answer}" != "y" ]] && [[ "${answer}" != "n" ]]; do
  read -r -p "Please enter 'y' or 'n': " answer
done
echo

if [[ "${answer}" == "y" ]]; then
  echo "Click VERIFY YOUR OWNERSHIP."
  echo "Copy and paste the TXT record data below (starts with 'google-site-verification='):"
  read -r TXT_RECORD
  readonly TXT_RECORD
  echo
  gcloud --project "${PROJECT_ID}" dns record-sets transaction start --zone "${zone}"
  gcloud --project "${PROJECT_ID}" dns record-sets transaction add "${TXT_RECORD}" \
    --ttl=300 --type=TXT \
    --name="${DOMAIN}" \
    --zone="${zone}"
  if ! gcloud --project "${PROJECT_ID}" dns record-sets transaction execute --zone "${zone}"
  then
    echo "Failed to create the TXT record set. You may need to manually create the TXT"
    echo "record set with its TXT record data via the Cloud Console:"
    echo "https://console.cloud.google.com/net-services/dns/zones/${zone}?project=${PROJECT_ID}"
    exit 1
  fi

  echo
  echo "The TXT record has been configured for domain ${DOMAIN}."
  echo "Within a few minutes you should be able to finish domain verification at:"
  echo "${url}"
  echo
  echo "In the meantime, continue with the remaining prerequisites in the README"
else
  echo "Your custom domain is ready for use."
fi

# Create env.mk if not present
if [[ ! -f "${APP_ROOT}/env.mk" ]]; then
  cp "${APP_ROOT}"/env.mk.sample "${APP_ROOT}"/env.mk
fi

# Substitute default PROJECT_ID and DOMAIN values
sed -e "s/^PROJECT_ID=project-id$/PROJECT_ID=${PROJECT_ID}/" \
    -e "s/^DOMAIN=my-zone.cloud-tutorial.dev$/DOMAIN=${DOMAIN}/" \
    -i "${APP_ROOT}"/env.mk
