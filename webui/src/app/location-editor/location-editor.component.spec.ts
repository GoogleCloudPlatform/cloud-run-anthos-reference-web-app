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

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LocationEditorComponent } from './location-editor.component';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { DummyComponent } from '../testing/dummy.component';
import { ActivatedRouteStub } from '../testing/activated-route-stub';
import { InventoryService, Location } from 'api-client';
import { Observable, of } from 'rxjs';
import { SetFormValue } from '../testing/helpers';
import { Location as Loc } from '@angular/common';
import { ActivatedRoute } from '@angular/router';

describe('LocationEditorComponent', () => {
  let component: LocationEditorComponent;
  let fixture: ComponentFixture<LocationEditorComponent>;
  let activatedRouteSub: ActivatedRouteStub;
  let location: Loc;
  let inventoryService: InventoryService;
  let getLocationSpy: jasmine.Spy<
    (id: string, observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Location>
  >;
  let newLocationSpy: jasmine.Spy<
    (location?: Location, observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Location>
  >;
  let updateLocationSpy: jasmine.Spy<
    (id: string, location?: Location, observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Location>
  >;

  beforeEach(async(() => {
    activatedRouteSub = new ActivatedRouteStub();
    TestBed.configureTestingModule({
      declarations: [ LocationEditorComponent ],
      imports: [
        HttpClientTestingModule,
        MatButtonModule,
        MatCardModule,
        MatFormFieldModule,
        MatInputModule,
        MatProgressBarModule,
        MatProgressSpinnerModule,
        ReactiveFormsModule,
        RouterTestingModule,
        NoopAnimationsModule,
        RouterTestingModule.withRoutes([
          { path: 'locations', component: DummyComponent },
          { path: 'locations/:id', component: DummyComponent },
        ]),
      ],
      providers: [
        {provide: ActivatedRoute, useValue: activatedRouteSub}
      ],
    })
    .compileComponents();
    inventoryService = TestBed.inject(InventoryService);
    location = TestBed.inject(Loc);
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LocationEditorComponent);
    component = fixture.componentInstance;
    getLocationSpy = spyOn(inventoryService, 'getLocation');
    newLocationSpy = spyOn(inventoryService, 'newLocation');
    updateLocationSpy = spyOn(inventoryService, 'updateLocation');
  });

  it('should initialize', () => {
    expect(component).toBeTruthy();
    expect(component.isNew).toBeTruthy();
  });

  it('should load location with name and warehouse', () => {
    const testLocation = {id: '123', name: 'test', warehouse: 'test wh'};
    const expectFormValue = {id: '123', name: 'test', warehouse: 'test wh'};
    const locationId = testLocation.id;
    getLocationSpy.withArgs(locationId).and.returnValue(of(testLocation));
    activatedRouteSub.setParamMap({ id: locationId });
    fixture.detectChanges();
    expect(component).toBeTruthy();
    expect(getLocationSpy).toHaveBeenCalledTimes(1);
    expect(component.locationForm.value).toEqual(expectFormValue);
    expect(component.isNew).toBeFalsy();
  });

  it('should create location', () => {
    const value = {id: '', name: 'test', warehouse: 'wh'};
    const expectValue = {id: '123', name: 'test', warehouse: 'wh'};
    const spy = newLocationSpy;
    spy.withArgs(value).and.returnValue(of(expectValue));
    component.isNew = true;
    fixture.detectChanges();
    expect(component).toBeTruthy();
    SetFormValue(fixture, '[formControlName="name"]', value.name);
    SetFormValue(fixture, '[formControlName="warehouse"]', value.warehouse);
    fixture.debugElement.nativeElement.querySelector('button[type="submit"]').click();
    expect(spy).toHaveBeenCalledTimes(1);
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/locations/123');
    });
  });

  it('should update location', () => {
    const value = {id: '123', name: 'test', warehouse: 'wh'};
    const expectValue = {id: '123', name: 'test', warehouse: 'wh'};
    const spy = updateLocationSpy;
    spy.withArgs(value.id, value).and.returnValue(of(expectValue));
    component.isNew = false;
    fixture.detectChanges();
    expect(component).toBeTruthy();
    SetFormValue(fixture, '[formControlName="id"]', value.id);
    SetFormValue(fixture, '[formControlName="name"]', value.name);
    SetFormValue(fixture, '[formControlName="warehouse"]', value.warehouse);
    fixture.debugElement.nativeElement.querySelector('button[type="submit"]').click();
    expect(spy).toHaveBeenCalledTimes(1);
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/locations/123');
    });
  });

  it('should go back on create', () => {
    fixture.detectChanges();
    fixture.debugElement.nativeElement.querySelector('[data-testid="cancelBtn"]').click();
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/locations');
    });
    expect(component).toBeTruthy();
  });

  it('should go back on edit', () => {
    fixture.detectChanges();
    component.isNew = false;
    component.locationId = '123';
    fixture.debugElement.nativeElement.querySelector('[data-testid="cancelBtn"]').click();
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/locations/123');
    });
    expect(component).toBeTruthy();
  });
});
