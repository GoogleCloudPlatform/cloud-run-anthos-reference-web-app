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
import { browser, ExpectedConditions } from 'protractor';
import { expect } from 'chai';

import { ItemsPage } from '../pages/items.po';

const page = new ItemsPage();
const loadingSpinner = page.getLoadingSpinner();
let lastCount;
const loadingDelay = 200;

When('I go to Items page', async () => {
  await page.navigateTo();
  await browser.sleep(loadingDelay);
  await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
});

Then('I should see some items', async () =>  {
  lastCount = await page.getTableRows().count();
});

Then('I should see {int} more items', async (diff) =>  {
  const newCount = await page.getTableRows().count();
  expect(newCount).equals(lastCount + diff);
  lastCount = newCount;
});

Then('I should see {int} fewer items', async (diff) =>  {
  const newCount = await page.getTableRows().count();
  expect(newCount).equals(lastCount - diff);
  lastCount = newCount;
});

When('I fill in {string} with {string}', async (fieldName, value) => {
  await page.getTextbox(fieldName).clear();
  await page.getTextbox(fieldName).sendKeys(value);
});

When('I click {string} button', async (buttonName) => {
  await page.getButton(buttonName).click();
  await browser.sleep(loadingDelay);
  await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
});

When('I click Submit button', async () => {
  await page.getSubmitButton().click();
  await browser.wait(ExpectedConditions.invisibilityOf(page.getProgressBar()));
  await browser.sleep(loadingDelay);
  await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
});

Then('I should see Item with {string} and {string}', async (name, description) => {
  expect(await page.getItemTitle().getText()).equals(name);
  expect(await page.getItemDescription().getText()).equals(description);
});

When ('I click on link {string}', async (name) => {
  const itemLink = page.getItemLinkByName(name);
  await browser.driver.wait(ExpectedConditions.presenceOf(itemLink));
  await itemLink.click();
  await browser.sleep(loadingDelay);
  await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
});
