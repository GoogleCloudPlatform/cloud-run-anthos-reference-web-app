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

import { When, Then } from 'cucumber';
import { BasePage } from '../pages/base.po';
import { expect } from 'chai';
import { element, by } from 'protractor';

const page = new BasePage();

When('I go to {string} page', async (path) => {
  await page.navigateToPath(path);
});

When('I fill in {string} with {string}', async (fieldName, value) => {
  await page.getFormField(fieldName).clear();
  await page.getFormField(fieldName).sendKeys(value);
});

When('I click {string} button', async (buttonName) => {
  await page.clickButton(buttonName);
});

When ('I click on link {string}', async (name) => {
  await page.clickLinkByName(name);
});

When ('I click on icon button', async () => {
  const plusButton = element(by.css('button mat-icon'));
  await page.clickElement(plusButton);
});

When ('I select {string} in selector {string}', async (optionText, selectName) => {
  await page.getFormField(selectName).sendKeys(optionText);
});

When ('I check radio button {string}', async (value) => {
  const radioButton = element(by.css(`mat-radio-button[value=${value}]`));
  await page.clickElement(radioButton);
});

let lastCount = 0;

Then('I should see some entries', async () =>  {
  lastCount = await page.getTableRows().count();
});

Then('I should see {int} more entries', async (diff) =>  {
  const newCount = await page.getTableRows().count();
  expect(newCount).to.equals(lastCount + diff);
  lastCount = newCount;
});

Then('I should see {int} fewer entries', async (diff) =>  {
  const newCount = await page.getTableRows().count();
  expect(newCount).to.equal(lastCount - diff);
  lastCount = newCount;
});

Then('I should see the latest transaction is for item {string} in location {string} for {string}', async (item, loc, diff) => {
  const row = element(by.css('table[data-testid="transactions"] tbody tr:first-child'));
  // tslint:disable-next-line: no-unused-expression
  expect(await row.isPresent()).to.be.true;
  expect(await row.element(by.css('td.mat-column-item')).getText()).to.equal(item);
  expect(await row.element(by.css('td.mat-column-location')).getText()).to.equals(loc);
  expect(await row.element(by.css('td.mat-column-diff')).getText()).to.equals(diff);
});
