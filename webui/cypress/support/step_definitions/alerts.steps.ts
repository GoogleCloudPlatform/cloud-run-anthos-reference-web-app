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
import { BasePage } from '../../pages/base.po';
import { testItem } from '../../data.config';

const page = new BasePage();

When('I go to alerts page', () => {
  page.navigateToPath('alerts');
  cy.wait('@alertList');
});

Then('I should see the latest alert is for test item contains {string}', async (alertText) => {
  const firstRowSelector = 'table[data-testid="alerts"] tbody tr:first-child';
  cy.get(firstRowSelector + ' td.mat-column-item').should('have.text', testItem.Name);
  cy.get(firstRowSelector + ' td.mat-column-text').contains(alertText);
});

When('I dismiss the latest alert', async (alertText) => {
  const firstRowSelector = 'table[data-testid="alerts"] tbody tr:first-child';
  cy.get(firstRowSelector + ' td.mat-column-actions button').click();
  cy.wait('@alertDelete');
});

