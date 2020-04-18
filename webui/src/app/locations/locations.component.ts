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
import { InventoryService, Location } from 'api-client';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'app-locations',
  templateUrl: './locations.component.html',
  styleUrls: ['./locations.component.scss']
})
export class LocationsComponent implements OnInit {
  displayedColumns: string[] = ['name', 'warehouse'];
  dataSource = new MatTableDataSource<Location>();
  loading = false;

  constructor(
    private inventoryService: InventoryService) {
  }

  ngOnInit() {
    this.loadLocations();
  }

  loadLocations(): void {
    this.loading = true;
    this.inventoryService.listLocations().subscribe( res => {
      this.dataSource.data = res;
      this.loading = false;
    }, () => {
      this.dataSource.data = [];
      this.loading = false;
    });
  }
}
