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
import { FormControl } from '@angular/forms';
import { InventoryService, Item } from 'api-client';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'app-items',
  templateUrl: './items.component.html',
  styleUrls: ['./items.component.scss'],
})
export class ItemsComponent implements OnInit {
  displayedColumns: string[] = ['name', 'description'];
  dataSource = new MatTableDataSource<Item>();
  loading = false;

  formControl = new FormControl('');

  constructor(
    private inventoryService: InventoryService,
  ) { }

  ngOnInit() {
    this.getItems();
  }

  getItems(): void {
    this.loading = true;
    this.inventoryService.listItems().subscribe(items => {
      this.dataSource.data = items;
      this.loading = false;
    }, () => {
      this.dataSource.data = [];
      this.loading = false;
    });
  }

}
