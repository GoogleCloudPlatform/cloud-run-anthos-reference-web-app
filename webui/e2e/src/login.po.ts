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
import {email, password} from './credentials';

export class LoginPage {
  navigateTo() {
    return browser.get('/login');
  }

  getEmailTextbox() {
    return element(by.css('input[formcontrolname=email]'));
  }

  getPasswordTextbox() {
    return element(by.css('input[formcontrolname=password]'));
  }

  getLoginButton() {
    return element(by.css('button#loginButton'));
  }

  getSnackBarMessage() {
    return element(by.css('simple-snack-bar > span'));
  }

  getSnackBarButton() {
    return element(by.css('simple-snack-bar button'));
  }

  login() {
    this.navigateTo();
    this.getEmailTextbox().sendKeys(email);
    this.getPasswordTextbox().sendKeys(password);
    this.getLoginButton().click();
    browser.sleep(2000);
  }
}
