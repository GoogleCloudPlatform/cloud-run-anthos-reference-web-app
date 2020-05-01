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

import { browser, by, element, ExpectedConditions, ElementFinder, ElementArrayFinder } from 'protractor';

export class BasePage {
  loadingDelay = 200;
  loadingSpinner: ElementFinder;
  progressBar: ElementFinder;

  constructor() {
    this.loadingSpinner = this.getLoadingSpinner();
    this.progressBar = this.getProgressBar();
    // Disable this because browser.waitForAngular() never resolve, and all tests end up timing out.
    // Tested with ChromeDriver@81.0.4044.69, Protractor@5.4.4, Angular@9.1.0.
    // If these versions are updated, we might try again to enable this.
    browser.waitForAngularEnabled(false);
  }

  async navigateToPath(path: string) {
    await browser.get(`/${path}`);
    await browser.sleep(this.loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(this.loadingSpinner));
  }

  getPageTitle(): ElementFinder  {
    return element(by.css('mat-card-title'));
  }

  getTableRows(): ElementArrayFinder {
    return element.all(by.tagName('tr'));
  }

  getFormField(name: string): ElementFinder  {
    return element(by.css(`[formcontrolname=${name}]`));
  }

  async waitForElement(css: string) {
    await browser.wait(ExpectedConditions.presenceOf(element(by.css(css))));
  }

  getButton(name: string): ElementFinder  {
    return element(by.buttonText(name));
  }

  async clickButton(name: string) {
    this.getButton(name).click();
    await browser.sleep(this.loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(this.progressBar));
    await browser.sleep(this.loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(this.loadingSpinner));
  }

  getLinkByName(name: string): ElementFinder {
    return element(by.linkText(name));
  }

  async clickElement(link: ElementFinder) {
    await browser.wait(ExpectedConditions.presenceOf(link));
    await link.click();
    await browser.sleep(this.loadingDelay);
    await browser.wait(ExpectedConditions.invisibilityOf(this.loadingSpinner));
  }

  async clickLinkByName(name: string) {
    await this.clickElement(this.getLinkByName(name));
  }

  getLoadingSpinner(): ElementFinder  {
    return element(by.css('mat-progress-spinner'));
  }

  getProgressBar(): ElementFinder  {
    return element(by.css('mat-progress-bar'));
  }

}
