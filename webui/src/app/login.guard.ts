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

import { Injectable } from '@angular/core';
import { Router, CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable, from, of, forkJoin } from 'rxjs';
import { AngularFireAuth } from '@angular/fire/auth';
import { InventoryService } from 'api-client';
import { map, mergeMap } from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { UserService } from 'user-svc-client';

@Injectable({
  providedIn: 'root'
})
export class LoginGuard implements CanActivate {

  userRole = '';

  constructor(
    private afAuth: AngularFireAuth,
    private router: Router,
    private inventoryService: InventoryService,
    private userService: UserService,
    private snackBar: MatSnackBar,
  ) {
  }

  canActivate( next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
    return true;
    // return (this.afAuth.user.pipe(
    //   mergeMap(user => {
    //     if (!user) {
    //       this.router.navigate(['/login'], { queryParams: { returnUrl: state.url }});
    //       this.snackBar.open('Please login first', '', { duration: 2000, });
    //       return of(false);
    //     }

    //     return forkJoin({
    //       idTokenResult: from(user.getIdTokenResult()),
    //       idToken: from(user.getIdToken())
    //     }).pipe(
    //       map( ({idToken, idTokenResult}) => {
    //         if (!idToken) {
    //           this.router.navigate(['/login'], { queryParams: { returnUrl: state.url }});
    //           this.snackBar.open('Failed to retrieve user token.', '', { duration: 5000, });
    //           return false;
    //         }
    //         if (idTokenResult && idTokenResult.claims.role) {
    //           this.userRole = idTokenResult.claims.role;
    //         }
    //         if (next.data.roles && next.data.roles.indexOf(this.userRole) < 0) {
    //           this.snackBar.open(`You don't have permission to go there.`, '', { duration: 5000, });
    //           return false;
    //         }
    //         const headers = new HttpHeaders({
    //           Authorization: 'Bearer ' + idToken,
    //         });

    //         this.inventoryService.defaultHeaders = headers;
    //         this.userService.defaultHeaders = headers;
    //         return true;
    //       })
    //     );
    //   })
    // ));
  }

}
