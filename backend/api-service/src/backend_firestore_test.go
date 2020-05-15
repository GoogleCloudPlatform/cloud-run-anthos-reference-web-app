// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build emulator

package service

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
)

const testProjectId = "foo"

func clearFirestoreBackend(t *testing.T) DatabaseBackend {
	t.Helper()
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, testProjectId)
	if err != nil {
		t.Fatalf("error creating firestore client: %v", err)
	}

	crefs, err := client.Collections(ctx).GetAll()
	if err != nil {
		t.Fatalf("error getting all collections: %v", err)
	}
	for _, cref := range crefs {
		drefs, err := cref.DocumentRefs(ctx).GetAll()
		if err != nil {
			t.Fatalf("error getting all documents in a collection: %v", err)
		}
		for _, dref := range drefs {
			_, err := dref.Delete(ctx)
			if err != nil {
				t.Fatalf("error deleting document: %v", err)
			}
		}
	}
	return NewFirestoreBackend(testProjectId)
}

var firestoreBackendTester = backendTester{
	resetBackend: func(t *testing.T) DatabaseBackend {
		return clearFirestoreBackend(t)
	},
	initBackend: func(t *testing.T, state initialBackendState) DatabaseBackend {
		t.Helper()
		backend := clearFirestoreBackend(t)
		ctx := context.Background()
		client, err := firestore.NewClient(ctx, testProjectId)
		if err != nil {
			t.Fatalf("error creating firestore client: %v", err)
		}
		for id, data := range state.inventories {
			dref := client.Collection("inventories").Doc(id)
			_, err := dref.Create(ctx, data)
			if err != nil {
				t.Fatalf("error creating doc: %v", err)
			}
		}
		for id, data := range state.inventoryTransactions {
			dref := client.Collection("inventoryTransactions").Doc(id)
			_, err := dref.Create(ctx, data)
			if err != nil {
				t.Fatalf("error creating doc: %v", err)
			}
		}
		for id, data := range state.items {
			dref := client.Collection("items").Doc(id)
			_, err := dref.Create(ctx, data)
			if err != nil {
				t.Fatalf("error creating doc: %v", err)
			}
		}
		for id, data := range state.locations {
			dref := client.Collection("locations").Doc(id)
			_, err := dref.Create(ctx, data)
			if err != nil {
				t.Fatalf("error creating doc: %v", err)
			}
		}
		for id, data := range state.alerts {
			dref := client.Collection("alerts").Doc(id)
			_, err := dref.Create(ctx, data)
			if err != nil {
				t.Fatalf("error creating doc: %v", err)
			}
		}

		return backend
	},
}

func TestFSDeleteItem(t *testing.T) {
	firestoreBackendTester.testDeleteItem(t)
}

func TestFSDeleteItemNotFound(t *testing.T) {
	firestoreBackendTester.testDeleteItemNotFound(t)
}

func TestFSDeleteLocation(t *testing.T) {
	firestoreBackendTester.testDeleteLocation(t)
}

func TestFSDeleteLocationNotFound(t *testing.T) {
	firestoreBackendTester.testDeleteLocationNotFound(t)
}

func TestFSDeleteAlert(t *testing.T) {
	firestoreBackendTester.testDeleteAlert(t)
}

func TestFSDeleteAlertNotFound(t *testing.T) {
	firestoreBackendTester.testDeleteAlertNotFound(t)
}

func TestFSGetInventoryTransaction(t *testing.T) {
	firestoreBackendTester.testGetInventoryTransaction(t)
}

func TestFSGetInventoryTransactionNotFound(t *testing.T) {
	firestoreBackendTester.testGetInventoryTransactionNotFound(t)
}

func TestFSGetItem(t *testing.T) {
	firestoreBackendTester.testGetItem(t)
}

func TestFSGetItemNotFound(t *testing.T) {
	firestoreBackendTester.testGetItemNotFound(t)
}

func TestFSGetLocation(t *testing.T) {
	firestoreBackendTester.testGetLocation(t)
}

func TestFSGetLocationNotFound(t *testing.T) {
	firestoreBackendTester.testGetLocationNotFound(t)
}

func TestFSListItems(t *testing.T) {
	firestoreBackendTester.testListItems(t)
}

func TestFSListItemInventory(t *testing.T) {
	firestoreBackendTester.testListItemInventory(t)
}

func TestFSListItemInventoryTransactions(t *testing.T) {
	firestoreBackendTester.testListInventoryTransactions(t)
}

func TestFSListInventoryTransactions(t *testing.T) {
	firestoreBackendTester.testListInventoryTransactions(t)
}

func TestFSListLocations(t *testing.T) {
	firestoreBackendTester.testListLocations(t)
}

func TestFSListLocationInventory(t *testing.T) {
	firestoreBackendTester.testListLocations(t)
}

func TestFSListLocationInventoryTransactions(t *testing.T) {
	firestoreBackendTester.testListLocationInventoryTransactions(t)
}

func TestFSListAlerts(t *testing.T) {
	firestoreBackendTester.testListAlerts(t)
}

func TestFSNewItem(t *testing.T) {
	firestoreBackendTester.testNewItem(t)
}

func TestFSNewInventoryTransaction(t *testing.T) {
	// Running this test against the Firestore emulator raises an Unknown RPC
	// error when querying the "inventories" collection not reproducible against
	// an actual Firestore database:
	t.Skip("this test is known to fail due to unknown reasons")
	firestoreBackendTester.testNewInventoryTransaction(t)
}

func TestFSNewInventoryTransactionNotFoundErrors(t *testing.T) {
	firestoreBackendTester.testNewInventoryTransactionNotFoundErrors(t)
}

func TestFSNewLocation(t *testing.T) {
	firestoreBackendTester.testNewLocation(t)
}

func TestFSNewAlert(t *testing.T) {
	firestoreBackendTester.testNewAlert(t)
}

func TestFSUpdateItem(t *testing.T) {
	firestoreBackendTester.testUpdateItem(t)
}

func TestFSUpdateItemNotFound(t *testing.T) {
	firestoreBackendTester.testUpdateItemNotFound(t)
}

func TestFSUpdateLocation(t *testing.T) {
	firestoreBackendTester.testUpdateLocation(t)
}

func TestFSUpdateLocationNotFound(t *testing.T) {
	firestoreBackendTester.testUpdateLocationNotFound(t)
}
