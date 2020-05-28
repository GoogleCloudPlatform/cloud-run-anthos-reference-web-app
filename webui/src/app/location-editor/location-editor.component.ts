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
import { InventoryService } from 'api-client';
import { Router, ActivatedRoute } from '@angular/router';
import { FormGroup, FormControl } from '@angular/forms';

@Component({
  selector: 'app-location-editor',
  templateUrl: './location-editor.component.html',
  styleUrls: ['./location-editor.component.scss']
})
export class LocationEditorComponent implements OnInit {
  isNew = true;
  loading = false;
  submitting = false;
  locationId: string | null = null;

  locationForm = new FormGroup({
    name: new FormControl(''),
    warehouse: new FormControl(''),
    id: new FormControl(''),
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private inventoryService: InventoryService,
  ) { }

  ngOnInit() {
    this.route.paramMap.subscribe(pmap => {
      this.locationId = pmap.get('id');
      if (this.locationId) {
        this.isNew = false;
        this.loading = true;
        this.inventoryService.getLocation(this.locationId).subscribe(location => {
          this.locationForm.setValue(
            Object.assign({
              name: '',
              warehouse: '',
            }, location)
          );
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    });
  }

  onSubmit() {
    const location = this.locationForm.value;
    this.submitting = true;
    if (this.isNew) {
      this.inventoryService.newLocation(location).subscribe((newLocation) => {
        this.router.navigate(['/locations', newLocation.id]);
        this.submitting = false;
      }, () => {
        this.submitting = false;
      });
    } else {
      this.inventoryService.updateLocation(location.id, location).subscribe((newLocation) => {
        this.router.navigate(['/locations', newLocation.id]);
        this.submitting = false;
      }, () => {
        this.submitting = false;
      });
    }
  }

  onCancel() {
    if (this.isNew) {
      this.router.navigate(['/locations']);
    } else {
      this.router.navigate(['/locations', this.locationId]);
    }
  }

}
