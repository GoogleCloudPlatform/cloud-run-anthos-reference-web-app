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

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ItemsComponent } from './items/items.component';
import { HomeComponent } from './home/home.component';
import { ItemEditorComponent } from './item-editor/item-editor.component';
import { LocationsComponent } from './locations/locations.component';
import { ItemViewComponent } from './item-view/item-view.component';
import { UsersComponent } from './users/users.component';
import { LocationEditorComponent } from './location-editor/location-editor.component';
import { LocationViewComponent } from './location-view/location-view.component';
import { AlertsComponent } from './alerts/alerts.component';
import { LoginGuard } from './login.guard';
import { LoginComponent } from './login/login.component';

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
    pathMatch: 'full',
    canActivate: [LoginGuard],
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: 'items',
    component: ItemsComponent,
    canActivate: [LoginGuard],
  },
  {
    path: 'items/new',
    component: ItemEditorComponent,
    canActivate: [LoginGuard],
    data: {roles: ['admin']},
  },
  {
    path: 'items/:id',
    component: ItemViewComponent,
    canActivate: [LoginGuard],
  },
  {
    path: 'items/:id/edit',
    component: ItemEditorComponent,
    canActivate: [LoginGuard],
    data: {roles: ['admin']},
  },
  {
    path: 'locations',
    component: LocationsComponent,
    canActivate: [LoginGuard],
  },
  {
    path: 'locations/new',
    component: LocationEditorComponent,
    canActivate: [LoginGuard],
    data: {roles: ['admin']},
  },
  {
    path: 'locations/:id',
    component: LocationViewComponent,
    canActivate: [LoginGuard],
  },
  {
    path: 'locations/:id/edit',
    component: LocationEditorComponent,
    canActivate: [LoginGuard],
    data: {roles: ['admin']},
  },
  {
    path: 'users',
    component: UsersComponent,
    canActivate: [LoginGuard],
  },
  {
    path: 'users/:id',
    component: UsersComponent,
    canActivate: [LoginGuard],
  },
  {
    path: 'alerts',
    component: AlertsComponent,
    canActivate: [LoginGuard],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
