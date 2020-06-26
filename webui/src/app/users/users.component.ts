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
import { MatTableDataSource } from '@angular/material/table';
import { UserService, User } from 'user-svc-client';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss']
})
export class UsersComponent implements OnInit {
  displayedColumns: string[] = ['name', 'email', 'role'];
  roleList: string[] = ['', 'worker', 'admin'];
  dataSource = new MatTableDataSource<any>();
  userRoles: {[name: string]: string} = {};
  loading = false;

  constructor(private userService: UserService) {
  }

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    this.loading = true;
    this.userService.listUsers().subscribe(
      data => {
        this.loading = false;
        this.dataSource.data = data;
        this.initUserRoles();
      },
      (e) => {
        this.loading = false;
      }
    );
  }
  initUserRoles() {
    this.dataSource.data.forEach(user => {
      if (!user.uid) {
        return;
      }
      if (user.customClaims && user.customClaims.role) {
        this.userRoles[user.uid] = user.customClaims.role;
      } else {
        this.userRoles[user.uid] = '';
      }
    });
  }

  onRoleChange(user: User) {
    if (user.uid) {
      const newRole = this.userRoles[user.uid];
      if (newRole === '' || newRole === 'worker' || newRole === 'admin') {
        this.userService.updateUser(user.uid, newRole).subscribe(
          () => this.loadData(),
          () => this.initUserRoles(),
        );
        return;
      }
    }
    // reset userRoles on failure to wipe temporary values.
    this.initUserRoles();
  }
}
