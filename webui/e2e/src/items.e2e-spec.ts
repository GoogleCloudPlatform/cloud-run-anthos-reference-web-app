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

import { browser, logging, element, by, ExpectedConditions, ElementFinder } from 'protractor';
import { LoginPage } from './login.po';
import { ItemsPage } from './items.po';

describe('Item pages', () => {
  let loginPage: LoginPage;
  let page: ItemsPage;
  let loadingSpinner: ElementFinder;
  const testItemName = 'test item 1';
  const loadingDelay = 200;
  let startCount;

  beforeEach(() => {
    loginPage = new LoginPage();
    loginPage.login();
    page = new ItemsPage();
    loadingSpinner = page.getLoadingSpinner();
  });

  it('should create item', async () => {
    page.navigateTo();
    // on item list
    const cardTitle = page.getPageTitle();
    await browser.wait(ExpectedConditions.presenceOf(cardTitle));
    expect(await cardTitle.getText()).toEqual('Items');
    loadingSpinner = page.getLoadingSpinner();
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    startCount = await page.getTableRows().count();
    page.getCreateButton().click();
    // on item create page
    page.getNameTextbox().sendKeys(testItemName);
    page.getDescriptionTextbox().sendKeys('description of test item');
    page.getSubmitButton().click();
    // on item view page
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    expect(page.getItemTitle().getText()).toEqual(testItemName);
    expect(page.getItemDescription().getText()).toEqual('description of test item');
  });

  it('should edit item', async () => {
    page.navigateTo();
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    const itemLink = page.getItemLinkByName(testItemName);
    await browser.driver.wait(ExpectedConditions.presenceOf(itemLink));
    itemLink.click();
    page.getEditButton().click();
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    page.getDescriptionTextbox().clear();
    page.getDescriptionTextbox().sendKeys('edited description of test item');
    page.getSubmitButton().click();
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    expect(page.getItemDescription().getText()).toEqual('edited description of test item');
  });

  it('should delete item', async () => {
    page.navigateTo();
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    expect(page.getTableRows().count()).toEqual(startCount + 1);
    page.getItemLinkByName(testItemName).click();
    // on item view page again
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    page.getDeleteButton().click();
    browser.sleep(loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(loadingSpinner));
    expect(page.getTableRows().count()).toEqual(startCount);
  });

  afterEach(async () => {
    // Assert that there are no errors emitted from the browser
    const logs = await browser.manage().logs().get(logging.Type.BROWSER);
    expect(logs).not.toContain(jasmine.objectContaining({
      level: logging.Level.SEVERE,
    } as logging.Entry));
  });
});
