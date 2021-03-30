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

import { Component, OnInit } from '@angular/core';

import { AngularFireAuth } from '@angular/fire/auth';
import { AuthProcessService } from 'ngx-auth-firebaseui';

import { firebase } from `firebase/app`;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'CRfA Canonical Web App';

  private photoUrl: string | null = null;

  constructor(
    public afAuth: AngularFireAuth,
    public authProcess: AuthProcessService,
  ) {
  }

  ngOnInit(): void {
    this.afAuth.onAuthStateChanged((u: firebase.User | null) => {
      if (u) {
        this.photoUrl = u.photoURL;
      }
    });
  }

  public get avatarImageUrl(): string | null {
    return this.photoUrl;
  }

  signOut() {
    this.authProcess.signOut()
      .catch(e => console.error('An error happened while signing out!', e));
  }
}
