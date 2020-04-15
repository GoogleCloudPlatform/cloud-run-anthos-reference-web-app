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
import { FirestoreService } from 'src/app/firestore/firestore.service';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss']
})
export class UsersComponent implements OnInit {
  displayedColumns: string[] = ['uid', 'displayName', 'email', 'role'];
  roleList: string[] = ['', 'worker', 'admin'];
  dataSource = new MatTableDataSource<any>();
  constructor(private fs: FirestoreService) {
    this.fs.getUsers().valueChanges().subscribe(
      data => {
        console.log(data);
        this.dataSource.data = data;
      }
    );
  }

  ngOnInit() {}
}
