
# E2E test for the frontend

## Prerequisite

### Prepare a test user

1. Go to [Identity Platform Providers][] to enable "Email/Password" provider.
1. Go to [Identity Platform Users][] to create a user for e2e test.

### Configure test credentials

With the credentials of the test user, use `credentials.sample.ts` as an
example to create `credentials.ts`

## Run the test

You can run the e2e test against any running instance of the this application.

Both frontend and backend need to be working together under the same URL.

```bash
npm run e2e -- --config baseUrl=$TARGET_URL
```

`$TARGET_URL` could be `http://localhost:4200` if you run it locally,
or `https://$DOMAIN` if you run it on Google Cloud.
