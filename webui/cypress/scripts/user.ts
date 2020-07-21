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
const {adminEmail, adminPassword, workerEmail, workerPassword} = require('../credentials');

const ensureUser = async (email: string, password: string, role: string, displayName: string) => {
  try {
    const userRecord = await admin.auth().getUserByEmail(email);
    console.log(`Test user [${displayName}](${email}) found with id ${userRecord.uid}.`);
    if (userRecord.customClaims && userRecord.customClaims.role !== role) {
      return updateUserRoleClaim(userRecord.uid, role);
    }
  } catch (error) {
    console.log(`[${displayName}](${email}) not found, creating one.`);
    try {
      const newUser = await admin.auth().createUser({
        email,
        password,
        displayName,
        disabled: false
      });
      // See the UserRecord reference doc for the contents of userRecord.
      console.log('Successfully created new user:', newUser.uid);
      return updateUserRoleClaim(newUser.uid, role);
    } catch (e) {
      console.log('Error creating new user:', e);
    }
  }
};

const updateUserRoleClaim = async (uid: string, role: string) => {
  try {
    await admin.auth().setCustomUserClaims(uid, { role });
    console.log(`Successfully update role to [${role}] for user ${uid}.`);
  } catch (error) {
    console.log('Error updating role claim:', error);
  }
};

const deleteUserByEmail = async (email: string) => {
  const userRecord = await admin.auth().getUserByEmail(email);
  console.log(`Cleaning up test user ${email}`);
  await admin.auth().deleteUser(userRecord.uid);
};

export const ensureUsers = async () => {
  await ensureUser(adminEmail, adminPassword, 'admin', 'Test Admin');
  await ensureUser(workerEmail, workerPassword, 'worker', 'Test Worker');
};

export const cleanupUsers = async () => {
  await deleteUserByEmail(adminEmail);
  await deleteUserByEmail(workerEmail);
};

export const setUserRole = async (email: string, role: string) => {
  const userRecord = await admin.auth().getUserByEmail(email);
  if (userRecord.customClaims && userRecord.customClaims.role !== role) {
    return updateUserRoleClaim(userRecord.uid, role);
  }
};
