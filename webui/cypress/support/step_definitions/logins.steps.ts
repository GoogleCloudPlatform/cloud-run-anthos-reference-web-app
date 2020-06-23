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

import { Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

import { LoginPage } from '../../pages/login.po';

const page = new LoginPage();

Given('I log(ged) in as admin', async () => {
  await page.loginAsAdmin();
});

Given('I log(ged) in as worker', async () => {
  await page.loginAsWorker();
});

Then('my avatar image should be set', async () => {
  page.getAvatar().should('exist');
});

When('I go to users page', () => {
  page.navigateToPath('users');
  cy.wait('@userList');
});

Then('I should see user with name {string} and role {string}', (name, role) => {
  cy.contains('tbody td.mat-column-name', name).siblings().contains('td.mat-column-role', role);
});

When('I select role {string} for user {string}', (role, name) => {
  cy.contains('tbody td.mat-column-name', name).siblings().find('mat-select').click();
  cy.contains('mat-option', role).click();
  cy.wait('@userUpdate');
  cy.wait('@userList');
});
