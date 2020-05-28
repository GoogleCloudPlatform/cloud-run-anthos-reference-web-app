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
import { FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { InventoryService } from 'api-client';

@Component({
  selector: 'app-item-editor',
  templateUrl: './item-editor.component.html',
  styleUrls: ['./item-editor.component.scss']
})
export class ItemEditorComponent implements OnInit {
  isNew = true;
  loading = false;
  submitting = false;
  itemId: string | null = null;

  itemForm = new FormGroup({
    name: new FormControl(''),
    description: new FormControl(''),
    id: new FormControl(''),
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private inventoryService: InventoryService,
  ) { }

  ngOnInit() {
    this.route.paramMap.subscribe(pmap => {
      this.itemId = pmap.get('id');
      if (this.itemId) {
        this.isNew = false;
        this.loading = true;
        this.inventoryService.getItem(this.itemId).subscribe(item => {
          this.itemForm.setValue(
            Object.assign({
              name: '',
              description: '',
            }, item)
          );
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    });
  }

  onSubmit() {
    const item = this.itemForm.value;
    this.submitting = true;
    if (this.isNew) {
      this.inventoryService.newItem(item).subscribe((newItem) => {
        this.router.navigate(['/items', newItem.id]);
        this.submitting = false;
      }, () => {
        this.submitting = false;
      });
    } else {
      this.inventoryService.updateItem(item.id, item).subscribe((newItem) => {
        this.router.navigate(['/items', newItem.id]);
        this.submitting = false;
      }, () => {
        this.submitting = false;
      });
    }
  }

  onCancel() {
    if (this.isNew) {
      this.router.navigate(['/items']);
    } else {
      this.router.navigate(['/items', this.itemId]);
    }
  }
}
