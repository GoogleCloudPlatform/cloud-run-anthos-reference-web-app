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

export class BasePage {

  async navigateToPath(path: string) {
    cy.visit(`/${path}`, {failOnStatusCode: false});
  }

  getPageTitle()  {
    return cy.get('mat-card-title');
  }

  getTableRows() {
    return cy.get('tbody tr');
  }

  getFormField(name: string)  {
    return cy.get(`[formcontrolname=${name}]`);
  }

  getButton(name: string)  {
    return cy.contains('button', name);
  }

  getLinkByName(name: string) {
    return cy.contains('a', name);
  }

  getLoadingSpinner()  {
    return cy.get('mat-progress-spinner');
  }

  getProgressBar()  {
    return cy.get('mat-progress-bar');
  }

}
