# Cloud Run for Anthos Reference Web App

This repository, including all associated workflows and automations, represents
an opinionated set of best practices aimed at demonstrating a reference architecture
for building a web application on Google Cloud using Cloud Run for Anthos.

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

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2FGoogleCloudPlatform%2Fcloud-run-anthos-reference-web-app&cloudshell_git_branch=master)

#### Local Setup

Follow the steps to [set up gcloud][] in your local environment,
then `git clone` this repo.

### Custom Domain

For this reference application to work properly, you will need a custom domain
that has been set up properly as described below.

You can follow the steps at [cloud-tutorial.dev][] to get a
custom subdomain and managed zone that are ready to use.

#### Managed Zone

Your custom domain must be associated with a [Cloud DNS Managed Zone][] in the
same project that you are using for this application.

To re-use an existing custom domain that you already own,
be sure to [update name server records][] to point to your managed zone.

#### Confirm Ownership of Domain

In order for the reference application to work, you must complete
[domain ownership verification][].

You can create a TXT record using the following steps:

1. In [Cloud DNS][], navigate to the managed zone associated with your custom domain.
1. Click **Add record set**.
1. Set the **Resource Record Type** to **TXT**.
1. In the **TXT data** field, paste the TXT record provided from following the
   [domain ownership verification][] steps.
1. Click **Create**.

### Setup Identity Platform for Auth

1. Follow [Setting up OAuth 2.0 guide][] to setup [OAuth consent screen][].
1. Enable Identity Platform and add **Google** as an Identity provider:
   * Go to the [Identity Platform page in the GCP console][].
   * Select your project from the **Select a project** drop-down.
   * Click **Enable Identity Platform**.
   * On the **Providers** page, click **Add a provider**.
   * Select **Google** from the list.
   * Fill in the **Web Client ID** and **Web Client Secret** fields with those
     from the OAuth client ID created in the previous step.
1. Add your custom domain as Authorized Domain on
[Identity Platform -> Settings][] page, Security tab.
1. Follow the example in [webui/firebaseConfig.js.sample](webui/firebaseConfig.js.sample)
   to create `webui/firebaseConfig.js`
   * **apiKey** and **authDomain** can be found following the
    [Identity Platform quickstart guide][]

### Set up Firestore security rules

Add the following security rules from [`firestore/firestore.rules`](firestore/firestore.rules)
to your Firebase project in the rules tab at
<https://console.firebase.google.com/project/$PROJECT_ID/database/firestore/rules>

## Deploying the Application for the First Time

This project uses [Cloud Build][] and [Config Connector][] to automate code and
infrastructure deployments.
The instructions below describe how to deploy the application.

### 1. Configure GCP Project

You will need to bootstrap the services and permissions required by this example.
The easiest way to do so is by running [bootstrap.sh](bootstrap.sh):

```bash
./bootstrap.sh $PROJECT_ID
```

### 2. Create an environment file

Copy [env.mk.sample](env.mk.sample) to `env.mk`:

```bash
cp env.mk.sample env.mk
```

### 3. Fill out TODO sections in `env.mk`

Address the TODO comment at the top of `env.mk` and ensure values are correct.

### 4. Create a GKE Cluster

Run `make cluster`

### 5. Add a verified owner for the domain

Add the following service account as an [additional verified owner][]:

`cnrm-system@${PROJECT_ID}.iam.gserviceaccount.com`

where `${PROJECT_ID}` is replaced by your Google Cloud project ID.

### 6. Build and deploy

Run `make build-all`

## Try Out the Application

Once your application is deployed, you can try it out by navigating to `https://$DOMAIN`,
where `$DOMAIN` is the custom domain
you configured in `env.mk`.

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
However, you must manually delete your Cloud Run service and GKE Cluster.

[Cloud Build]: https://cloud.google.com/cloud-build/docs
[Config Connector]: https://cloud.google.com/config-connector/docs
[Cloud DNS Managed Zone]: https://cloud.google.com/dns/zones
[Cloud DNS]: https://console.cloud.google.com/net-services/dns/zones
[update name server records]: https://cloud.google.com/dns/docs/migrating#update_your_registrars_name_server_records
[domain ownership verification]: https://cloud.google.com/storage/docs/domain-name-verification#verification
[additional verified owner]: https://cloud.google.com/storage/docs/domain-name-verification?_ga=2.256052552.-234301672.1582050261#additional_verified_owners
[Identity Platform quickstart guide]: https://cloud.google.com/identity-platform/docs/quickstart-email-password#sign_the_user_in
[Identity Platform page in the GCP console]: https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity
[OAuth consent screen]: https://console.cloud.google.com/apis/credentials/consent
[Identity Platform -> Settings]: https://console.cloud.google.com/customer-identity/settings
[Setting up OAuth 2.0 guide]: https://support.google.com/cloud/answer/6158849?hl=en
[set up gcloud]: https://cloud.google.com/sdk/docs
[`makefile`]: makefile
[Owner permission]: https://console.cloud.google.com/iam-admin/roles/details/roles%3Cowner
[cloud-tutorial.dev]: https://cloud-tutorial.dev/