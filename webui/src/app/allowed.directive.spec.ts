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

import { AllowedDirective } from './allowed.directive';
import { LoginGuard } from './login.guard';
import { Component } from '@angular/core';
import { async, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { AngularFireAuth } from '@angular/fire/auth';
import { RouterTestingModule } from '@angular/router/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('AllowedDirective', () => {
  @Component({
    template: `<div *appAllowed="['admin']"><span data-testid="content">test</span></div>`
  })
  class TestComponent {
  }

  let loginGuard: LoginGuard;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AllowedDirective, TestComponent ],
      imports: [
        HttpClientTestingModule,
        RouterTestingModule,
        MatSnackBarModule,
      ],
      providers: [
        {
          provide: AngularFireAuth,
          useValue: {},
        }
      ]
    })
    .compileComponents();
    loginGuard = TestBed.inject(LoginGuard);
  }));

  it('should not display without allowed role', () => {
    loginGuard.userRole = '';
    const fixture = TestBed.createComponent(TestComponent);
    fixture.detectChanges();
    const directive = fixture.debugElement.queryAllNodes(By.directive(AllowedDirective));
    expect(directive).toBeTruthy();
    console.log(fixture.debugElement);
    const testComponent = fixture.debugElement.query(By.css('span[data-testid="content"]'));
    expect(testComponent).toBeFalsy();
  });

  it('should display with allowed role', () => {
    loginGuard.userRole = 'admin';
    const fixture = TestBed.createComponent(TestComponent);
    fixture.detectChanges();
    const directive = fixture.debugElement.queryAllNodes(By.directive(AllowedDirective));
    expect(directive).toBeTruthy();
    console.log(fixture.debugElement);
    const testComponent = fixture.debugElement.query(By.css('span[data-testid="content"]'));
    expect(testComponent).toBeTruthy();
  });


});
