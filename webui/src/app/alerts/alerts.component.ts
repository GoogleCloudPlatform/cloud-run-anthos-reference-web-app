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
import { Alert, AlertService, Item, InventoryService } from 'api-client';
import * as moment from 'moment';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'app-alerts',
  templateUrl: './alerts.component.html',
  styleUrls: ['./alerts.component.scss']
})
export class AlertsComponent implements OnInit {
  displayedColumns: string[] = ['actions', 'text', 'item', 'time'];
  dataSource = new MatTableDataSource<Alert>();
  items: Item[] = [];
  loading = false;

  constructor(
    private alertService: AlertService,
    private inventoryService: InventoryService,
  ) { }

  ngOnInit() {
    this.getAlerts();
    this.inventoryService.listItems().subscribe(items => {
      this.items = items;
    });
  }

  sortByTime(a: Alert, b: Alert): number {
    return moment(b.timestamp).unix() - moment(a.timestamp).unix();
  }

  handleDismiss(alert: Alert): void {
    if (alert.id !== undefined) {
      this.alertService.deleteAlert(alert.id).subscribe(() => {
        this.dataSource.data.splice(this.dataSource.data.indexOf(alert), 1);
        this.dataSource._updateChangeSubscription();
      });
    }
  }

  getAlerts(): void {
    this.loading = true;
    this.alertService.listAlerts().subscribe(alerts => {
      if (alerts) {
        this.dataSource.data = alerts.sort(this.sortByTime);
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
}
