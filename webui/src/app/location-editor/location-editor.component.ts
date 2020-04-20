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
import { Location } from '@angular/common';
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

  locationForm = new FormGroup({
    name: new FormControl(''),
    warehouse: new FormControl(''),
    id: new FormControl(''),
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private inventoryService: InventoryService,
    private location: Location,
  ) { }

  ngOnInit() {
    const locationId = this.route.snapshot.paramMap.get('id');
    if (locationId) {
      this.isNew = false;
      this.loading = true;
      this.inventoryService.getLocation(locationId).subscribe(location => {
        this.locationForm.setValue(location);
        this.loading = false;
      }, () => {
        this.loading = false;
      });
    }
  }

  onSubmit() {
    const location = this.locationForm.value;
    if (this.isNew) {
      this.inventoryService.newLocation(location).subscribe((newLocation) => {
        this.router.navigate(['/locations', newLocation.id]);
      });
    } else {
      this.inventoryService.updateLocation(location.id, location).subscribe((newLocation) => {
        this.router.navigate(['/locations', newLocation.id]);
      });
    }
  }

  onCancel() {
    this.location.back();
  }

}
