
# E2E test for the frontend

## Prerequisite

### Setup firebase private key

This is not needed if run the e2e test on Cloud Build,
in which the Cloud Build service account already have permission configured.

Follow the [Firebase Admin Initialize SDK]
to setup `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

## Run the test

You can run the e2e test against any running instance of the this application.

Both frontend and backend need to be working together under the same URL.

```bash
npm run e2e -- --config baseUrl=$TARGET_URL
```

`$TARGET_URL` could be `http://localhost:4200` if you run it locally,
or `https://$DOMAIN` if you run it on Google Cloud.

[Firebase Admin Initialize SDK]: https://firebase.google.com/docs/admin/setup#initialize-sdk
