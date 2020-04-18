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
import {Location} from '@angular/common';
import { FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { InventoryService } from 'api-client';

@Component({
  selector: 'app-item-editor',
  templateUrl: './item-editor.component.html',
  styleUrls: ['./item-editor.component.scss']
})
export class ItemEditorComponent implements OnInit {
  isNew = true;
  loading = false;

  itemForm = new FormGroup({
    name: new FormControl(''),
    description: new FormControl(''),
    id: new FormControl(''),
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private inventoryService: InventoryService,
    private location: Location,
  ) { }

  ngOnInit() {
    const itemId = this.route.snapshot.paramMap.get('id');
    if (itemId) {
      this.isNew = false;
      this.loading = true;
      this.inventoryService.getItem(itemId).subscribe(item => {
        this.itemForm.setValue(item);
        this.loading = false;
      }, () => {
        this.loading = false;
      });
    }
  }

  onSubmit() {
    const item = this.itemForm.value;
    if (this.isNew) {
      this.inventoryService.newItem(item).subscribe((newItem) => {
        this.router.navigate(['/items', newItem.id]);
      });
    } else {
      this.inventoryService.updateItem(item.id, item).subscribe((newItem) => {
        this.router.navigate(['/items', newItem.id]);
      });
    }
  }

  onCancel() {
    this.location.back();
  }
}
