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
package service

import (
	"testing"
)

var inMemoryBackendTester = backendTester{
	resetBackend: func(t *testing.T) DatabaseBackend {
		t.Helper()
		return NewInMemoryBackend()
	},
	initBackend: func(t *testing.T, state initialBackendState) DatabaseBackend {
		t.Helper()
		mb := NewInMemoryBackend()
		if state.inventoryTransactions != nil {
			mb.inventoryTransactions = state.inventoryTransactions
		}
		if state.items != nil {
			mb.items = state.items
		}
		if state.locations != nil {
			mb.locations = state.locations
		}

		for _, inv := range state.inventories {
			if mb.inventoryByItemByLocationIndex[inv.ItemId] == nil {
				mb.inventoryByItemByLocationIndex[inv.ItemId] = make(map[string]*Inventory)
			}
			if mb.inventoryByLocationByItemIndex[inv.LocationId] == nil {
				mb.inventoryByLocationByItemIndex[inv.LocationId] = make(map[string]*Inventory)
			}
			mb.inventoryByItemByLocationIndex[inv.ItemId][inv.LocationId] = inv
			mb.inventoryByLocationByItemIndex[inv.LocationId][inv.ItemId] = inv
		}
		return mb
	},
}

func TestIMBDeleteItem(t *testing.T) {
	inMemoryBackendTester.testDeleteItem(t)
}

func TestIMBDeleteItemNotFound(t *testing.T) {
	inMemoryBackendTester.testDeleteItemNotFound(t)
}

func TestIMBDeleteLocation(t *testing.T) {
	inMemoryBackendTester.testDeleteLocation(t)
}

func TestIMBDeleteLocationNotFound(t *testing.T) {
	inMemoryBackendTester.testDeleteLocationNotFound(t)
}

func TestIMBGetInventoryTransaction(t *testing.T) {
	inMemoryBackendTester.testGetInventoryTransaction(t)
}

func TestIMBGetInventoryTransactionNotFound(t *testing.T) {
	inMemoryBackendTester.testGetInventoryTransactionNotFound(t)
}

func TestIMBGetItem(t *testing.T) {
	inMemoryBackendTester.testGetItem(t)
}

func TestIMBGetItemNotFound(t *testing.T) {
	inMemoryBackendTester.testGetItemNotFound(t)
}

func TestIMBGetLocation(t *testing.T) {
	inMemoryBackendTester.testGetLocation(t)
}

func TestIMBGetLocationNotFound(t *testing.T) {
	inMemoryBackendTester.testGetLocationNotFound(t)
}

func TestIMBListItems(t *testing.T) {
	inMemoryBackendTester.testListItems(t)
}

func TestIMBListItemInventory(t *testing.T) {
	inMemoryBackendTester.testListItemInventory(t)
}

func TestIMBListItemInventoryTransactions(t *testing.T) {
	inMemoryBackendTester.testListInventoryTransactions(t)
}

func TestIMBListInventoryTransactions(t *testing.T) {
	inMemoryBackendTester.testListInventoryTransactions(t)
}

func TestIMBListLocations(t *testing.T) {
	inMemoryBackendTester.testListLocations(t)
}

func TestIMBListLocationInventory(t *testing.T) {
	inMemoryBackendTester.testListLocations(t)
}

func TestIMBListLocationInventoryTransactions(t *testing.T) {
	inMemoryBackendTester.testListLocationInventoryTransactions(t)
}

func TestIMBNewItem(t *testing.T) {
	inMemoryBackendTester.testNewItem(t)
}

func TestIMBNewInventoryTransaction(t *testing.T) {
	inMemoryBackendTester.testNewInventoryTransaction(t)
}

func TestIMBNewInventoryTransactionNotFoundErrors(t *testing.T) {
	inMemoryBackendTester.testNewInventoryTransactionNotFoundErrors(t)
}

func TestIMBNewLocation(t *testing.T) {
	inMemoryBackendTester.testNewLocation(t)
}

func TestIMBUpdateItem(t *testing.T) {
	inMemoryBackendTester.testUpdateItem(t)
}

func TestIMBUpdateItemNotFound(t *testing.T) {
	inMemoryBackendTester.testUpdateItemNotFound(t)
}

func TestIMBUpdateLocation(t *testing.T) {
	inMemoryBackendTester.testUpdateLocation(t)
}

func TestIMBUpdateLocationNotFound(t *testing.T) {
	inMemoryBackendTester.testUpdateLocationNotFound(t)
}
