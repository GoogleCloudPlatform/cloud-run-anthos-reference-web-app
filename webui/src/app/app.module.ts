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

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { AngularFireModule } from '@angular/fire';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { BASE_PATH, ApiModule, Configuration, ConfigurationParameters } from 'api-client';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatRadioModule } from '@angular/material/radio';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ReactiveFormsModule } from '@angular/forms';

import { NgxAuthFirebaseUIModule } from 'ngx-auth-firebaseui';
import { ItemsComponent } from './items/items.component';
import { ItemEditorComponent } from './item-editor/item-editor.component';
import { LocationsComponent } from './locations/locations.component';
import { UsersComponent } from './users/users.component';
import { environment } from 'src/environments/environment';
import { ItemViewComponent } from './item-view/item-view.component';
import { InventoryTransactionsComponent } from './inventory-transactions/inventory-transactions.component';
import { firebaseConfig } from '../../firebaseConfig';
import { LocationEditorComponent } from './location-editor/location-editor.component';
import { LocationViewComponent } from './location-view/location-view.component';
import { InventoryTransactionDialogComponent } from './inventory-transaction-dialog/inventory-transaction-dialog.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';

export function apiConfigFactory(): Configuration {
  const params: ConfigurationParameters = {
    // set configuration parameters here.
  };
  return new Configuration(params);
}

export function appNameFactory() {
  return 'CRfA Canonical Web App';
}

@NgModule({
  declarations: [
    AppComponent,
    ItemsComponent,
    ItemEditorComponent,
    LocationsComponent,
    ItemViewComponent,
    InventoryTransactionsComponent,
    UsersComponent,
    LocationEditorComponent,
    LocationViewComponent,
    InventoryTransactionDialogComponent,
    HomeComponent,
    LoginComponent,
  ],
  entryComponents: [
    InventoryTransactionDialogComponent,
  ],
  imports: [
    ApiModule,
    HttpClientModule,
    BrowserModule,
    AppRoutingModule,
    MatToolbarModule,
    MatButtonModule,
    MatCardModule,
    MatMenuModule,
    MatTableModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatSnackBarModule,
    BrowserAnimationsModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatRadioModule,
    AngularFireModule.initializeApp(firebaseConfig),
    NgxAuthFirebaseUIModule.forRoot(firebaseConfig, appNameFactory,
    {
      authGuardFallbackURL: '/login',
      authGuardLoggedInURL: '/',
    }),
  ],
  providers: [
    { provide: BASE_PATH, useValue: environment.API_BASE_PATH },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
