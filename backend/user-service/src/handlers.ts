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

import * as admin from 'firebase-admin';
import * as express from 'express';

admin.initializeApp({});
admin.auth();

export const listUsers = async (
  request: express.Request,
  response: express.Response
) => {
  try {
    const result = await admin.auth().listUsers();
    return response.send(result.users);
  } catch (e) {
    return response.status(500).send(e);
  }
};

export const getUser = async (
  request: express.Request,
  response: express.Response
) => {
  const uid = request.params.uid;
  try {
    const user = await admin.auth().getUser(uid);
    return response.send(user);
  } catch (e) {
    if (e.code === 'auth/user-not-found') {
      return response.status(404).send(e);
    }
    return response.status(500).send(e);
  }
};

export const updateUser = async (
  request: express.Request,
  response: express.Response
) => {
  const uid = request.params.uid;
  const role = request.query.role;
  if (uid && typeof uid === 'string' && role) {
    try {
      await admin.auth().setCustomUserClaims(uid, {role});
      return response.sendStatus(201);
    } catch (e) {
      if (e.code === 'auth/user-not-found') {
        return response.status(404).send(e);
      }
      return response.status(500).send(e);
    }
  }
  return response.sendStatus(400);
};
