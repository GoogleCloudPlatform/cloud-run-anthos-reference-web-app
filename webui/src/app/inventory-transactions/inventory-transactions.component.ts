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

import { Component, OnInit, Input } from '@angular/core';
import { InventoryTransaction, InventoryService, Item, Location } from 'api-client';
import { MatTableDataSource } from '@angular/material/table';
import * as moment from 'moment';
import { MatDialog } from '@angular/material/dialog';
import { InventoryTransactionDialogComponent } from '../inventory-transaction-dialog/inventory-transaction-dialog.component';

@Component({
  selector: 'app-inventory-transactions',
  templateUrl: './inventory-transactions.component.html',
  styleUrls: ['./inventory-transactions.component.scss']
})
export class InventoryTransactionsComponent implements OnInit {
  @Input() itemId: string | null = null;
  @Input() locationId: string | null = null;

  displayedColumns: string[] = ['item', 'location', 'diff', 'time', 'note', 'createdBy'];
  dataSource = new MatTableDataSource<InventoryTransaction>();
  items: Item[] = [];
  locations: Location[] = [];
  loading = false;

  constructor(
    private inventoryService: InventoryService,
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    if (this.itemId) {
      this.loadTransactionByItemId(this.itemId);
    } else if (this.locationId) {
      this.loadTransactionByLocationId(this.locationId);
    }

    this.inventoryService.listItems().subscribe(items => {
      this.items = items;
    });
    this.inventoryService.listLocations().subscribe(locations => {
      this.locations = locations;
    });
  }

  sortByTime(a: InventoryTransaction, b: InventoryTransaction): number {
    return moment(b.timestamp).unix() - moment(a.timestamp).unix();
  }

  loadTransactionByItemId(itemId: string): void {
    this.loading = true;
    this.inventoryService.listItemInventoryTransactions(itemId).subscribe(transactions => {
      if (transactions) {
        this.dataSource.data = transactions.sort(this.sortByTime);
      } else {
        this.dataSource.data = [];
      }
      this.loading = false;
    }, () => {
      this.dataSource.data = [];
      this.loading = false;
    });
  }

  loadTransactionByLocationId(locationId: string): void {
    this.loading = true;
    this.inventoryService.listLocationInventoryTransactions(locationId).subscribe(transactions => {
      if (transactions) {
        this.dataSource.data = transactions.sort(this.sortByTime);
      } else {
        this.dataSource.data = [];
      }
      this.loading = false;
    }, () => {
      this.dataSource.data = [];
      this.loading = false;
    });
  }

  getTime(time: Date): string {
    return moment(time).fromNow();
  }

  getItemName(id: string): string {
    const item = this.items.find((i) => i.id === id);
    if (item) {
      return item.name;
    }
    return id;
  }

  getLocationName(id: string): string {
    const location = this.locations.find((i) => i.id === id);
    if (location) {
      return location.name;
    }
    return id;
  }

  addTransaction() {
    const dialogRef = this.dialog.open(InventoryTransactionDialogComponent, {
      width: '450px',
      disableClose: true,
      data: {
        transaction: {
          item_id: this.itemId,
          location_id: this.locationId,
        },
        items: this.items,
        locations: this.locations,
      },
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result.event === 'submit') {
        this.loadData();
      }
    });
  }
}
