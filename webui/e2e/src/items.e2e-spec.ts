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

import { browser, logging, element, by, ExpectedConditions } from 'protractor';
import { LoginPage } from './login.po';
import { ItemsPage } from './items.po';

describe('Item pages', () => {
  let loginPage: LoginPage;
  let page: ItemsPage;

  beforeEach(() => {
    loginPage = new LoginPage();
    loginPage.login();
    page = new ItemsPage();
  });

  it('should create item', async () => {
    const testItemName = 'test item 1';
    page.navigateTo();
    // on item list
    const cardTitle = page.getPageTitle();
    await browser.driver.wait(ExpectedConditions.presenceOf(cardTitle));
    expect(await cardTitle.getText()).toEqual('Items');
    const startCount = await page.getTableRows().count();
    page.getCreateButton().click();
    // on item create page
    page.getNameTextbox().sendKeys(testItemName);
    page.getDescriptionTextbox().sendKeys('description of test item');
    page.getSubmitButton().click();
    // on item view page
    expect(page.getItemTitle().getText()).toEqual(testItemName);
    expect(page.getItemDescription().getText()).toEqual('description of test item');
    page.getEditButton().click();
    // on item edit page
    page.getDescriptionTextbox().clear();
    page.getDescriptionTextbox().sendKeys('edited description of test item');
    page.getSubmitButton().click();
    expect(page.getItemDescription().getText()).toEqual('edited description of test item');
    // on item view page
    page.getBackButton().click();
    // back to item list
    expect(page.getTableRows().count()).toEqual(startCount + 1);
    page.getItemLinkByName(testItemName).click();
    // on item view page again
    page.getDeleteButton().click();
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
