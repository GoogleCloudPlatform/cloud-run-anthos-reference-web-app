/**
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Then, When } from 'cypress-cucumber-preprocessor/steps';

import { LocationsPage } from '../../pages/locations.po';

const page = new LocationsPage();

Then('I should see Location named {string}', async (name) => {
  page.getLocationTitle().should('have.text', name);
});

Then('I should see Location in warehouse {string}', async (warehouse) => {
  page.getLocationWarehouse().should('contain.text', warehouse);
});

When('I go to locations page', () => {
  page.navigateToPath('locations');
  cy.wait('@locationList');
});

When('wait for location to load', () => {
  cy.wait('@locationGet');
  cy.wait('@invTransList');
});
