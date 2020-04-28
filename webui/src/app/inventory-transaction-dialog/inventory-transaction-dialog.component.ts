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

import { Component, OnInit, Inject } from '@angular/core';
import { InventoryTransaction, InventoryService, Item, Location } from 'api-client';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

export interface InventoryTransactionDialogData {
  transaction: InventoryTransaction;
  items: Item[];
  locations: Location[];
}
@Component({
  selector: 'app-inventory-transaction-dialog',
  templateUrl: './inventory-transaction-dialog.component.html',
  styleUrls: ['./inventory-transaction-dialog.component.scss']
})
export class InventoryTransactionDialogComponent implements OnInit {
  inventoryTransactionForm = new FormGroup({
    item_id: new FormControl('', Validators.required),
    location_id: new FormControl('', Validators.required),
    action: new FormControl('', Validators.required),
    count: new FormControl('', Validators.required),
    note: new FormControl(''),
  });
  items: Item[] = [];
  locations: Location[] = [];
  submitting = false;

  constructor(
    private dialogRef: MatDialogRef<InventoryTransactionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) private data: InventoryTransactionDialogData,
    private inventoryService: InventoryService,
  ) {
    this.inventoryTransactionForm.patchValue(Object.assign( {action: 'ADD'}, data.transaction));
  }

  ngOnInit() {
    this.inventoryService.listItems().subscribe(items => {
      this.items = items;
    });
    this.inventoryService.listLocations().subscribe(locations => {
      this.locations = locations;
    });
  }

  onCancel(): void {
    this.dialogRef.close({event: 'cancel'});
  }

  onSubmit(): void {
    const data = this.inventoryTransactionForm.value;
    this.submitting = true;
    this.inventoryService.newInventoryTransaction(data).subscribe(() => {
      this.dialogRef.close({event: 'submit'});
      this.submitting = false;
    }, () => {
      this.submitting = false;
    });
  }
}
