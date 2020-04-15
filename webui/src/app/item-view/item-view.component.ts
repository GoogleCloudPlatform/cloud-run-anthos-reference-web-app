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
import { ActivatedRoute, Router } from '@angular/router';
import { Item, InventoryService } from 'api-client';

@Component({
  selector: 'app-item-view',
  templateUrl: './item-view.component.html',
  styleUrls: ['./item-view.component.scss']
})
export class ItemViewComponent implements OnInit {
  item: Item = null;
  loading = false;

  constructor(
    private inventoryService: InventoryService,
    private route: ActivatedRoute,
    private router: Router,
  ) { }

  ngOnInit() {
    const itemId = this.route.snapshot.paramMap.get('id');
    if (itemId) {
      this.getItem(itemId);
    }
  }

  getItem(itemId: string): void {
    if (itemId) {
      this.loading = true;
      this.inventoryService.getItem(itemId).subscribe(item => {
        this.item = item;
        console.log(item);
        this.loading = false;
      }, (error) => {
        console.error(error);
        this.item = null;
        this.loading = false;
      });
    }
  }

  handleDelete(): void {
    if (this.item) {
      this.inventoryService.deleteItem(this.item.id).subscribe(() => {
        this.router.navigate(['/items']);
      });
    }
  }
}
