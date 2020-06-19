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

import { TestBed, inject } from '@angular/core/testing';

import { LoginGuard } from './login.guard';
import { HttpHeaders } from '@angular/common/http';
import { AngularFireAuth } from '@angular/fire/auth';
import { InventoryService } from 'api-client';
import { RouterTestingModule } from '@angular/router/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('LoginGuard', () => {

  const AngularFireAuthMock: any = {
    authState: {},
    auth: {
      onAuthStateChanged() {
        return Promise.resolve();
      },
      currentUser: {
        getIdToken(val: boolean) {
          return Promise.resolve();
        }
      }
    },
  };

  class InventoryServiceMock {
    public defaultHeaders = new HttpHeaders();

    constructor() {}
  }

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        HttpClientTestingModule,
        RouterTestingModule,
        MatSnackBarModule,
      ],
      providers: [
        LoginGuard,
        {
          provide: AngularFireAuth,
          useValue: AngularFireAuthMock,
        },
        {
          provide: InventoryService,
          useClass: InventoryServiceMock,
        }
      ]
    });
  });

  it('should be created', inject([LoginGuard], (guard: LoginGuard) => {
    expect(guard).toBeTruthy();
  }));
});
