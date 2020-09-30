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
import { ItemsPage } from '../../pages/items.po';

const page = new ItemsPage();

Then('I should see Item named {string}', async (name) => {
  page.getItemTitle().should('have.text', name);
});

Then('I should see Item description as {string}', async (description) => {
  page.getItemDescription().should('contain.text', description);
});

When('I go to items page', () => {
  page.navigateToPath('items');
  cy.wait('@itemList');
});

When('wait for item to load', () => {
  cy.wait('@itemGet');
  cy.wait('@invTransList');
  cy.wait('@itemList');
  cy.wait('@locationList');
});
