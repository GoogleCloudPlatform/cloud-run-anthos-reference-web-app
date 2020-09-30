
# E2E test for the frontend

## Prerequisite

### Setup firebase private key

This is not needed if running the e2e test on Cloud Build,
when the Cloud Build service account already has permissions configured.

Follow the [Firebase Admin Initialize SDK]
to setup `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

## Run the test

### Running locally

You can run the e2e test against any running instance of the this application.

Both frontend and backend need to be working together under the same URL.

```bash
npm run e2e -- --config baseUrl=$TARGET_URL
```

`$TARGET_URL` could be `http://localhost:4200` if you run it locally,
or `https://$DOMAIN` if you run it on Google Cloud.

### Running in Cloud Build

You can also run the e2e tests in Cloud Build via [cloudbuild.yaml][].

In order to do this, edit `env.mk` and make sure that `TEST_ARTIFACTS_LOCATION`
is set to a valid GCS bucket location.

[Firebase Admin Initialize SDK]: https://firebase.google.com/docs/admin/setup#initialize-sdk
[cloudbuild.yaml]: ./cloudbuild.yaml