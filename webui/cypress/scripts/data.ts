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


import * as admin from 'firebase-admin' ;
const {testItem, testLocation} = require('../data.config');

export const ensureData = async () => {
  await ensureDocByName('items', testItem, 'ItemId');
  await ensureDocByName('locations', testLocation, 'LocationId');
};

export const cleanupData = async () => {
  await cleanupDocByName('items', testItem, 'ItemId');
  await cleanupDocByName('locations', testLocation, 'LocationId');
};

const ensureDocByName = async (collection: string, entry: any, key: string) => {
  const querySnapshot = await admin.firestore().collection(collection).where('Name', '==', entry.Name).get();
  if (querySnapshot.empty) {
    const newEntry = await admin.firestore().collection(collection).add(entry);
    await newEntry.set({...entry, Id: newEntry.id});
  } else {
    console.log(`Test data [${entry.Name}] found, checking stale data.`);
    for (const doc of querySnapshot.docs) {
      await cleanupDoc('inventoryTransactions', key, doc.id);
      await cleanupDoc('alerts', key, doc.id);
    }
  }
};

const cleanupDocByName = async (collection: string, entry: any, key: string) => {
  const querySnapshot = await admin.firestore().collection(collection).where('Name', '==', entry.Name).get();
  if (!querySnapshot.empty) {
    console.log(`Test data [${entry.Name}] found, checking stale data.`);
    for (const doc of querySnapshot.docs) {
      await cleanupDoc('inventoryTransactions', key, doc.id);
      await cleanupDoc('alerts', key, doc.id);
      await doc.ref.delete();
    }
  }
};

const cleanupDoc = async (collection: string, filterKey: string, filterValue: string) => {
  const qs = await admin.firestore().collection(collection).where(filterKey, '==', filterValue).get();
  if (!qs.empty) {
    console.log(`Found ${qs.docs.length} entries in ${collection}, cleaning up.`);
    for (const doc of qs.docs) {
      await doc.ref.delete();
    }
  }
};

