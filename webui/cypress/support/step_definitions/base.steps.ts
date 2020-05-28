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

import { When, Then, Before } from 'cypress-cucumber-preprocessor/steps';
import { BasePage } from '../../pages/base.po';

const page = new BasePage();

Before(async () => {
  cy.server();
  cy.route('GET', '/api/**').as('api');
  cy.route('GET', '/api/items').as('itemList');
  cy.route('GET', '/api/items/*').as('itemGet');
  cy.route('POST', '/api/items').as('itemCreate');
  cy.route('PUT', '/api/items/**').as('itemUpdate');
  cy.route('GET', '/api/locations').as('locationList');
  cy.route('GET', '/api/locations/*').as('locationGet');
  cy.route('DELETE', '/api/items/**').as('itemDelete');
  cy.route('POST', '/api/locations').as('locationCreate');
  cy.route('PUT', '/api/locations/**').as('locationUpdate');
  cy.route('DELETE', '/api/locations/**').as('locationDelete');
  cy.route('POST', '/api/inventoryTransactions').as('invTransCreate');
  cy.route('GET', '/api/*/*/inventoryTransactions').as('invTransList');
  cy.route('POST', 'https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword*').as('verifyPassword');
  cy.route('POST', 'https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo*').as('getAccountInfo');
});

When('I go to {string} page', async (path) => {
  page.navigateToPath(path);
  cy.wait('@api');
});

When('I fill in {string} with {string}', async (fieldName, value) => {
  page.getFormField(fieldName).clear();
  page.getFormField(fieldName).type(value);
});

When('I click {string} button', async (buttonName) => {
  page.getButton(buttonName).click();
});

When('I click {string} button and wait', async (buttonName) => {
  page.getButton(buttonName).click();
  cy.wait('@api');
});

When('wait for {string}', (alias) => {
  cy.wait(alias);
});

When ('I click on link {string}', async (name) => {
  page.getLinkByName(name).click();
});


When ('I click on the plus icon button and wait', async () => {
  cy.get('button mat-icon').click();
  cy.wait('@itemList');
  cy.wait('@locationList');
});

When ('I submit the transaction inventory', () => {
  page.getButton('Submit').click();
  cy.wait('@invTransCreate');
  cy.wait(1000);
  cy.wait('@invTransList');
});

When ('I select {string} in selector {string}', async (optionText, selectName) => {
  page.getFormField(selectName).click();
  cy.get('mat-option').contains(optionText).click();
});

When ('I check radio button {string}', async (value) => {
  cy.get(`mat-radio-button[value=${value}]`).click();
});

let lastCount = 0;

Then('I should see some entries', async () =>  {
  page.getTableRows().then(elm => lastCount = elm.length);
});

Then('I should see {int} more entries', async (diff) =>  {
  page.getTableRows().should('have.length', lastCount + diff);
  page.getTableRows().then(elm => lastCount = elm.length);
});

Then('I should see {int} fewer entries', async (diff) =>  {
  page.getTableRows().should('have.length', lastCount - diff);
  page.getTableRows().then(elm => lastCount = elm.length);
});

Then('I should see the latest transaction is for item {string} in location {string} for {string}', async (item, loc, diff) => {
  const firstRowSelector = 'table[data-testid="transactions"] tbody tr:first-child';
  cy.get(firstRowSelector + ' td.mat-column-item').should('have.text', item);
  cy.get(firstRowSelector + ' td.mat-column-location').should('have.text', loc);
  cy.get(firstRowSelector + ' td.mat-column-diff').should('have.text', diff);
});

Then('I should see page title {string}', async (title) => {
  page.getPageTitle().should('have.text', title);
});
