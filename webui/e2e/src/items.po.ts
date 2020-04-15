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

import { browser, by, element, ExpectedConditions } from 'protractor';

export class ItemsPage {
  navigateTo() {
    return browser.get('/items');
  }

  getCreateButton() {
    return element(by.css('button.create-btn'));
  }

  getPageTitle() {
    return element(by.css('mat-card-title'));
  }

  getTableRows() {
    return element.all(by.tagName('tr'));
  }

  getNameTextbox() {
    return element(by.css('input[formcontrolname=name]'));
  }

  getDescriptionTextbox() {
    return element(by.css('input[formcontrolname=description]'));
  }

  getSubmitButton() {
    return element(by.css('button[type=submit]'));
  }

  getItemTitle() {
    return element(by.css('.item-info mat-card-title'));
  }

  getItemDescription() {
    return element(by.css('.item-info mat-card-content'));
  }

  getEditButton() {
    return element(by.buttonText('Edit'));
  }

  getBackButton() {
    return element(by.buttonText('Back'));
  }

  getDeleteButton() {
    return element(by.buttonText('Delete'));
  }

  getItemLinkByName(name: string) {
    return element(by.linkText(name));
  }
}
