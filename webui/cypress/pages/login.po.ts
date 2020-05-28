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

// import { browser, by, element } from 'protractor';
import {email, password} from '../credentials';
import { BasePage } from './base.po';

export class LoginPage extends BasePage {

  async login() {
    this.navigateToPath('login');
    this.getFormField('email').type(email);
    this.getFormField('password').type(password);
    this.getButton('LOGIN').click();
    cy.wait('@verifyPassword');
    cy.wait('@getAccountInfo');
  }

}
