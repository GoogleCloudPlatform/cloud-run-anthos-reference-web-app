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

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { UsersComponent } from './users.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Observable, of } from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { UserService, User } from 'user-svc-client';
import { By } from '@angular/platform-browser';

describe('UsersComponent', () => {
  let component: UsersComponent;
  let fixture: ComponentFixture<UsersComponent>;
  let userService: UserService;
  let listUserSpy: jasmine.Spy<
    (observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Array<User>>
  >;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ UsersComponent ],
      imports: [
        HttpClientTestingModule,
        MatCardModule,
        MatTableModule,
      ],
    })
    .compileComponents();
    userService = TestBed.inject(UserService);
    listUserSpy = spyOn(userService, 'listUsers');
  }));

  const initComponent = () => {
    fixture = TestBed.createComponent(UsersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  };

  it('should create', () => {
    listUserSpy.and.returnValue(of([]));
    initComponent();
    expect(component).toBeTruthy();
    expect(listUserSpy).toHaveBeenCalledTimes(1);
  });

  it('should show users', () => {
    listUserSpy.and.returnValues(of([
      {
        uid: '12345abcde',
        email: 'worker@test.org',
        customClaims: { roles: 'worker' }
      },
      {
        uid: '23456bcdef',
        email: 'admin@test.org',
        customClaims: { roles: 'admin' }
      },
    ]));
    initComponent();
    expect(listUserSpy).toHaveBeenCalledTimes(1);
    const rows = fixture.debugElement.queryAll(By.css('tbody tr'));
    expect(rows.length).toEqual(2);
    expect(rows[0].query(By.css('.mat-column-email')).nativeElement.textContent).toContain('worker@test.org');
    expect(rows[1].query(By.css('.mat-column-email')).nativeElement.textContent).toContain('admin@test.org');
  });
});
