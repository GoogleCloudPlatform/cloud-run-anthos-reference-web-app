
import * as admin from 'firebase-admin';
import * as express from 'express';

admin.initializeApp({});
admin.auth();

export const listUsers = async (request: express.Request, response: express.Response) =>  {
  try {
    const result = await admin.auth().listUsers();
    return response.send(result.users);
  } catch (e) {
    return response.status(500).send(e);
  }
}

export const getUser = async (request: express.Request, response: express.Response) => {
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
}

export const updateUser = async (request: express.Request, response: express.Response) => {
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
}