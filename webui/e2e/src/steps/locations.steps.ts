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

import { Then, Given } from 'cucumber';
import { expect } from 'chai';

import { LocationsPage } from '../pages/locations.po';

const page = new LocationsPage();

Then('I should see Location named {string}', async (name) => {
  expect(await page.getLocationTitle().getText()).equals(name);
});

Then('I should see Location in warehouse {string}', async (warehouse) => {
  expect(await page.getLocationWarehouse().getText()).equals(warehouse);
});

Given('There is a location named {string}', async (name) => {
  await page.navigateTo();
  const link = page.getLinkByName(name);
  if (!await link.isPresent()) {
    page.clickButton('Create');
    await page.getFormField('name').sendKeys(name);
    await page.getFormField('warehouse').sendKeys(`WH ${name}`);
    await page.clickButton('Submit');
  }
});
