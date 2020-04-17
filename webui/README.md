# Cloud Run for Anthos Reference Web App - Frontend

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 8.3.20.

# Development

## Prerequisites

### Install Angular CLI

```bash
npm install -g @angular/cli
```

(If you see an `EACCES` error, see [1]

### `npm install`

Run `npm install` in the same directory as `package.json` to install all dependencies.

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.


Run with proxy `ng serve --proxy-config proxy.conf.json` if you are also running the backend locally on port 80.

### Setting up Backend

Replace API_BASE_PATH in `src/environments/environment.ts` with your desired backend endpoint if you want to target something besides `localhost:80/api` when local testing.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `--prod` flag for a production build.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).

[1]: https://docs.npmjs.com/resolving-eacces-permissions-errors-when-installing-packages-globally