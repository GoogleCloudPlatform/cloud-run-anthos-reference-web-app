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

import { ItemsComponent } from './items.component';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { InventoryService, Item } from 'api-client';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { of, Observable } from 'rxjs';
import { Location } from '@angular/common';
import { DummyComponent } from '../testing/dummy.component';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MockAllowedDirective } from '../allowed.directive.mock';

describe('ItemsComponent', () => {
  let component: ItemsComponent;
  let fixture: ComponentFixture<ItemsComponent>;
  let inventoryService: InventoryService;
  let location: Location;
  let listItemSpy: jasmine.Spy<
    (observe?: 'body', reportProgress?: boolean, options?: {})
      => Observable<Array<Item>>
  >;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ItemsComponent, DummyComponent, MockAllowedDirective ],
      imports: [
        HttpClientTestingModule,
        MatButtonModule,
        MatCardModule,
        MatIconModule,
        MatTableModule,
        MatProgressSpinnerModule,
        ReactiveFormsModule,
        RouterTestingModule.withRoutes([
          { path: 'items/new', component: DummyComponent },
          { path: 'items/:id', component: DummyComponent },
         ]),
        NoopAnimationsModule,
      ]
    })
    .compileComponents();
    inventoryService = TestBed.inject(InventoryService);
    location = TestBed.inject(Location);
    initComponent();
  }));

  const initComponent = () => {
    fixture = TestBed.createComponent(ItemsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  };

  it('should create', () => {
    listItemSpy = spyOn(inventoryService, 'listItems').and.callThrough();
    initComponent();
    expect(component).toBeTruthy();
    expect(listItemSpy).toHaveBeenCalledTimes(1);
  });

  const testListParams = [
    {
      description: 'should display empty list',
      items: [],
      count: 0,
    }, {
      description: 'should display list with 1 item',
      items: [
        { name: '测试产品', id: 'id-0', description: 'desc-0' },
      ],
      count: 1,
    }, {
      description: 'should display list with 2 items',
      items: [
        { name: 'test1', id: 'id-1', description: 'desc-1' },
        { name: 'test2', id: 'id-2', description: 'desc-2' }
      ],
      count: 2,
    }
  ];
  testListParams.forEach(param => {
    it(param.description, () => {
      listItemSpy = spyOn(inventoryService, 'listItems');
      listItemSpy.and.returnValue(of(param.items));
      initComponent();
      expect(component).toBeTruthy();
      expect(listItemSpy).toHaveBeenCalledTimes(1);
      expect(component.dataSource.data.length).toEqual(param.count);
      const ele: HTMLElement = fixture.nativeElement;
      param.items.forEach((item) => {
        expect(ele.textContent).toContain(item.name);
      });
    });
  });

  it('should navigate to create page', async(() => {
    const button = fixture.debugElement.nativeElement.querySelector('button.create-btn');
    button.click();

    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/items/new');
    });
  }));

  it('should navigate to item page', async(() => {
    listItemSpy = spyOn(inventoryService, 'listItems');
    listItemSpy.and.returnValue(of([{ name: 'test0', id: 'id-0', description: 'desc-0' }]));
    initComponent();

    const button = fixture.debugElement.nativeElement.querySelector('a.item-link:first-child');
    button.click();

    fixture.whenStable().then(() => {
      expect(location.path()).toEqual('/items/id-0');
    });
  }));


});
