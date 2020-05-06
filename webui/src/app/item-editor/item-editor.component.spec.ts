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

import { ItemEditorComponent } from './item-editor.component';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { ActivatedRoute } from '@angular/router';
import { ActivatedRouteStub } from '../testing/activated-route-stub';
import { InventoryService, Item } from 'api-client';
import { of, Observable } from 'rxjs';
import { Location } from '@angular/common';
import { DummyComponent } from '../testing/dummy.component';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { SetFormValue } from '../testing/helpers';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

describe('ItemEditorComponent', () => {
  let component: ItemEditorComponent;
  let fixture: ComponentFixture<ItemEditorComponent>;
  let activatedRouteSub: ActivatedRouteStub;
  let location: Location;
  let inventoryService: InventoryService;
  let getItemSpy: jasmine.Spy<
    (id: string, observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Item>
  >;
  let newItemSpy: jasmine.Spy<
    (item?: Item, observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Item>
  >;
  let updateItemSpy: jasmine.Spy<
    (id: string, item?: Item, observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Item>
  >;

  beforeEach(async(() => {
    activatedRouteSub = new ActivatedRouteStub();
    TestBed.configureTestingModule({
      declarations: [ ItemEditorComponent ],
      imports: [
        HttpClientTestingModule,
        MatButtonModule,
        MatCardModule,
        MatFormFieldModule,
        MatInputModule,
        MatProgressBarModule,
        MatProgressSpinnerModule,
        ReactiveFormsModule,
        NoopAnimationsModule,
        RouterTestingModule.withRoutes([
          { path: 'items', component: DummyComponent },
          { path: 'items/:id', component: DummyComponent },
        ]),
      ],
      providers: [
        {provide: ActivatedRoute, useValue: activatedRouteSub}
      ],
    })
    .compileComponents();
    inventoryService = TestBed.inject(InventoryService);
    location = TestBed.inject(Location);
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ItemEditorComponent);
    component = fixture.componentInstance;
    getItemSpy = spyOn(inventoryService, 'getItem');
    newItemSpy = spyOn(inventoryService, 'newItem');
    updateItemSpy = spyOn(inventoryService, 'updateItem');
  });

  it('should create', () => {
    expect(component).toBeTruthy();
    expect(component.isNew).toBeTruthy();
  });

  [
    {
      testDesc: 'should load item with only name',
      item: {id: '123', name: 'test'},
      expectFormValue: {id: '123', name: 'test', description: ''}
    },
    {
      testDesc: 'should load item with name and description',
      item: {id: '123', name: 'test', description: 'test desc'},
      expectFormValue: {id: '123', name: 'test', description: 'test desc'},
    }
  ].forEach(testParam => {
    it(testParam.testDesc, () => {
      const itemId = testParam.item.id;
      getItemSpy.withArgs(itemId).and.returnValue(of(testParam.item));
      activatedRouteSub.setParamMap({ id: itemId });
      fixture.detectChanges();
      expect(component).toBeTruthy();
      expect(getItemSpy).toHaveBeenCalledTimes(1);
      expect(component.itemForm.value).toEqual(testParam.expectFormValue);
      expect(component.isNew).toBeFalsy();
    });
  });

  it('should create item', () => {
    const value = {id: '', name: 'test', description: 'desc'};
    const expectValue = {id: '123', name: 'test', description: 'desc'};
    const spy = newItemSpy;
    spy.withArgs(value).and.returnValue(of(expectValue));
    component.isNew = true;
    fixture.detectChanges();
    expect(component).toBeTruthy();
    SetFormValue(fixture, '[formControlName="name"]', value.name);
    SetFormValue(fixture, '[formControlName="description"]', value.description);
    fixture.debugElement.nativeElement.querySelector('button[type="submit"]').click();
    expect(spy).toHaveBeenCalledTimes(1);
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/items/123');
    });
  });

  it('should update item', () => {
    const value = {id: '123', name: 'test', description: 'desc'};
    const expectValue = {id: '123', name: 'test', description: 'desc'};
    const spy = updateItemSpy;
    spy.withArgs(value.id, value).and.returnValue(of(expectValue));
    component.isNew = false;
    fixture.detectChanges();
    expect(component).toBeTruthy();
    SetFormValue(fixture, '[formControlName="id"]', value.id);
    SetFormValue(fixture, '[formControlName="name"]', value.name);
    SetFormValue(fixture, '[formControlName="description"]', value.description);
    fixture.debugElement.nativeElement.querySelector('button[type="submit"]').click();
    expect(spy).toHaveBeenCalledTimes(1);
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/items/123');
    });
  });

  it('should go back on create', () => {
    fixture.detectChanges();
    fixture.debugElement.nativeElement.querySelector('[data-testid="cancelBtn"]').click();
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/items');
    });
    expect(component).toBeTruthy();
  });

  it('should go back on edit', () => {
    fixture.detectChanges();
    component.isNew = false;
    component.itemId = '123';
    fixture.debugElement.nativeElement.querySelector('[data-testid="cancelBtn"]').click();
    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/items/123');
    });
    expect(component).toBeTruthy();
  });
});
