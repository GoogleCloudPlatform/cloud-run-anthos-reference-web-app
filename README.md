![npm-audit-periodic](https://github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/workflows/npm-audit-periodic/badge.svg)

**English** | [Español](docs/README_sp.md)

# Cloud Run for Anthos Reference Web App

This repository, including all associated workflows and automations, represents
an opinionated set of best practices aimed at demonstrating a reference architecture
for building a web application on Google Cloud using Cloud Run for Anthos.

A detailed description of the architecture of the web app can be found in [architecture.md][].

## Prerequisites

### Development Environment

*NOTE: the steps in this guide assume that you are working in a POSIX-based
development environment.*

The only requirement to run this example out of the box is a working
installation of `gcloud`. Optionally, having `make` installed will allow you
to make use of the convenience targets provided in the [`makefile`][].

*NOTE: Your `gcloud` user account must have [Owner permission][] in order
to complete setup of the application.*

#### Cloud Shell

This example can be run directly from Cloud Shell!

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2FGoogleCloudPlatform%2Fcloud-run-anthos-reference-web-app&cloudshell_git_branch=main)

#### Local Setup

Follow the steps to [set up gcloud][] in your local environment,
then `git clone` this repo.

### Custom Domain

For this reference application to work properly, you will need a custom domain
that has been set up properly and verified.

The easiest way to do this is by running the interactive script [domain-setup.sh][]:

```bash
./scripts/domain-setup.sh
```

This script:

* Allows you to create a custom subdomain or use an existing one.
* Creates custom subdomains and managed zones ready for use using the steps at
  [cloud-tutorial.dev][].
* Ensures any custom domains are associated with a [Cloud DNS Managed Zone][]
  in the same project that you are using for this application.
* For provided custom domains, links to documentation to
  [update name server records][] to point to your managed zone.
* Walks you through [domain ownership verification][].

### Identity Platform for Auth and Firestore Setup

1. [Enable Identity Platform][] for your project.
   * This will create an OAuth 2.0 Client ID that can be used by the web application.
   * This additionally creates a Firebase project where Cloud Firestore can be used.

1. Whitelist your custom domain in Identity Platform.
   * In the GCP console, navigate to [Identity Platform > Settings][].
   * Click on the **Security** tab.
   * Add your custom domain under **Authorized Domains**.
   * Click **Save**.

1. Authorize your OAuth 2.0 Client ID to be usable by your custom domain.
   * In the GCP console, navigate to [APIs & Services > Credentials][].
   * Click on the OAuth 2.0 Client ID that was auto created.
     * "(auto created by Google Service)" appears in the name.
     * **$PROJECT_ID.firebaseapp.com** _should_ appear under
       **Authorized JavaScript origins**.
   * Take note of the **Client ID** and **Client secret**.
     You'll use them in the next step.
   * Under **Authorized JavaScript origins**,
     add your custom domain prefixed with `https://`.
   * Click **Save**.

1. Add **Google** as an Identity Provider in Identity Platform:
   * In the GCP console, navigate to [Identity Platform > Providers][].
   * Click **Add a provider**.
   * Select **Google** from the list.
   * Fill in the **Web Client ID** and **Web Client Secret** fields with those
     from the OAuth 2.0 Client ID created in the previous step.
   * Click **Save**.

1. Configure the [OAuth consent screen][].
   * **User Type** can be set to either **Internal** or **External**.
   * You'll need to set the **Support email** and the
     **Application homepage link** (your custom domain prefixed with `https://`).
   * Additional information
     [here](https://support.google.com/cloud/answer/6158849?hl=en#userconsent).

1. Setup `webui/firebaseConfig.ts`.
   * Identify your Web API Key by navigating to Project Settings in the Firebase
     console:
     <https://console.firebase.google.com/project/$PROJECT_ID/settings/general>
   * Run [firebase-config-setup.sh][] to create `webui/firebaseConfig.ts`:

   ```bash
   ./scripts/firebase-config-setup.sh $PROJECT_ID $API_KEY
   ```

1. Create Firestore database:
   * Navigate to the Develop > Database in the Firebase console at:
     <https://console.firebase.google.com/project/$PROJECT_ID/database>.
   * Click **Create Database**
   * Choose **production mode**, then click **Next**
   * Use the default location, or customize it as desired, then click **Done**

1. Set up the Firestore security rules:
   * Navigate to the Develop > Database > Rules in the Firebase console at:
     <https://console.firebase.google.com/project/$PROJECT_ID/database/firestore/rules>.
   * Ensure that **Cloud Firestore** is selected in the dropdown above.
     ![firestore rules page screenshot][]
   * Set the security rules to the ones found in [`firestore/firestore.rules`][].

## Deploying the Application for the First Time

This project uses [Cloud Build][] and [Config Connector][] to automate code and
infrastructure deployments.
The instructions below describe how to deploy the application.

### 1. Configure GCP Project

You will need to bootstrap the services and permissions required by this example.
The easiest way to do so is by running [bootstrap.sh][]:

```bash
./scripts/bootstrap.sh $PROJECT_ID
```

This step additionally creates a file named `env.mk` based on [env.mk.sample](env.mk.sample).

### 2. Fill out TODO sections in `env.mk`

Address the TODO comment at the top of `env.mk` and ensure values are correct.

### 3. Create a GKE Cluster

Run `make cluster`

### 4. Add a verified owner for the domain

Add the following service account as an [additional verified owner][]:

`cnrm-system@${PROJECT_ID}.iam.gserviceaccount.com`

where `${PROJECT_ID}` is replaced by your Google Cloud project ID.

### 5. Build and deploy

Run `make build-all`.

## Try Out the Application

Once your application is deployed, you can try it out by navigating to `https://$DOMAIN`,
where `$DOMAIN` is the custom domain you configured in `env.mk`.

### Setup first admin user

After you login at least once to the app, you can use this script to make your
account an `admin`. Afterwards you'll be able to use the Users page to manage
other accounts. To use this script you will need to
[Initialize the Firebase Admin SDK][] and setup
`GOOGLE_APPLICATION_CREDENTIALS` environment variable.

```shell
cd webui
npm install
npm run init-admin <email>
```

## Update the Application

Running `make build-all` will rebuild and deploy the app, including any changes
made to the infrastructure. Note that removing resources from `infrastructure-tpl.yaml`
will not cause them to be deleted. You must either run `make delete` before removing
the resource (then redeploy with `make build-all` after removing it), or manually
delete the resource with `kubectl delete`.

```shell
# builds and deploys backend, frontend, and KCC infrastructure
make build-all

# builds and deploys only the backend Go service
make build-backend

# builds and deploys only the frontend angular webapp
make build-webui
```

## Cleanup

Running `make delete` will delete the Config Connector resources from your cluster,
which will cause Config Connector to delete the associated GCP resources.
However, you must manually delete your Cloud Run for Anthos service and GKE Cluster.

[APIs & Services > Credentials]: https://console.cloud.google.com/apis/credentials
[Cloud Build]: https://cloud.google.com/cloud-build/docs
[Config Connector]: https://cloud.google.com/config-connector/docs
[Cloud DNS Managed Zone]: https://cloud.google.com/dns/zones
[update name server records]: https://cloud.google.com/dns/docs/migrating#update_your_registrars_name_server_records
[domain ownership verification]: https://cloud.google.com/storage/docs/domain-name-verification#verification
[additional verified owner]: https://cloud.google.com/storage/docs/domain-name-verification?_ga=2.256052552.-234301672.1582050261#additional_verified_owners
[Enable Identity Platform]: https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity
[Identity Platform > Providers]: https://console.cloud.google.com/customer-identity/providers
[Identity Platform quickstart guide]: https://cloud.google.com/identity-platform/docs/quickstart-email-password#sign_the_user_in
[Identity Platform page in the GCP console]: https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity
[OAuth consent screen]: https://console.cloud.google.com/apis/credentials/consent
[Identity Platform > Settings]: https://console.cloud.google.com/customer-identity/settings
[Setting up OAuth 2.0 guide]: https://support.google.com/cloud/answer/6158849?hl=en
[set up gcloud]: https://cloud.google.com/sdk/docs
[Owner permission]: https://console.cloud.google.com/iam-admin/roles/details/roles%3Cowner
[cloud-tutorial.dev]: https://cloud-tutorial.dev/
[`makefile`]: makefile
[architecture.md]: ./docs/architecture.md
[bootstrap.sh]: scripts/bootstrap.sh
[firebase-config-setup.sh]: scripts/firebase-config-setup.sh
[domain-setup.sh]: scripts/domain-setup.sh
[firestore rules page screenshot]: docs/img/firestore_rules_page.png
[`firestore/firestore.rules`]: firestore/firestore.rules
[Initialize the Firebase Admin SDK]: https://firebase.google.com/docs/admin/setup#initialize-sdk
