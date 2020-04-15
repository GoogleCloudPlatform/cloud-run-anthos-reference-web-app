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
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// InMemoryBackend is an in memory backend that implements the Backend interface
type InMemoryBackend struct {
	items                          map[string]*Item
	locations                      map[string]*Location
	inventoryByItemByLocationIndex map[string]map[string]*Inventory
	inventoryByLocationByItemIndex map[string]map[string]*Inventory
	inventoryTransactions          map[string]*InventoryTransaction
}

func NewInMemoryBackend() *InMemoryBackend {
	return &InMemoryBackend{
		items:                          make(map[string]*Item),
		locations:                      make(map[string]*Location),
		inventoryByItemByLocationIndex: make(map[string]map[string]*Inventory),
		inventoryByLocationByItemIndex: make(map[string]map[string]*Inventory),
		inventoryTransactions:          make(map[string]*InventoryTransaction),
	}
}

func (mb *InMemoryBackend) DeleteItem(ctx context.Context, id string) error {
	if _, ok := mb.items[id]; ok {
		delete(mb.items, id)
		return nil
	}
	return ItemNotFound(id)
}

func (mb *InMemoryBackend) DeleteLocation(ctx context.Context, id string) error {
	if _, ok := mb.locations[id]; ok {
		delete(mb.locations, id)
		return nil
	}
	return LocationNotFound(id)
}

func (mb *InMemoryBackend) GetInventoryTransaction(ctx context.Context, id string) (*InventoryTransaction, error) {
	if transaction, ok := mb.inventoryTransactions[id]; ok {
		return transaction, nil
	}
	return nil, InventoryTransactionNotFound(id)
}

func (mb *InMemoryBackend) GetItem(ctx context.Context, id string) (*Item, error) {
	if item, ok := mb.items[id]; ok {
		return item, nil
	}
	return nil, ItemNotFound(id)
}

func (mb *InMemoryBackend) GetLocation(ctx context.Context, id string) (*Location, error) {
	if loc, ok := mb.locations[id]; ok {
		return loc, nil
	}
	return nil, LocationNotFound(id)
}

func (mb *InMemoryBackend) ListItems(ctx context.Context) ([]*Item, error) {
	items := make([]*Item, 0, len(mb.items))
	for _, item := range mb.items {
		items = append(items, item)
	}
	return items, nil
}

func (mb *InMemoryBackend) ListItemInventory(ctx context.Context, id string) ([]*Inventory, error) {
	var inventories []*Inventory
	if itemInvs, ok := mb.inventoryByItemByLocationIndex[id]; ok {
		for _, inventory := range itemInvs {
			inventories = append(inventories, inventory)
		}
	}
	return inventories, nil
}

func (mb *InMemoryBackend) ListItemInventoryTransactions(ctx context.Context, id string) ([]*InventoryTransaction, error) {
	var txns []*InventoryTransaction
	for _, txn := range mb.inventoryTransactions {
		if txn.ItemId == id {
			txns = append(txns, txn)
		}
	}
	return txns, nil
}

func (mb *InMemoryBackend) ListInventoryTransactions(ctx context.Context) ([]*InventoryTransaction, error) {
	transactions := make([]*InventoryTransaction, 0, len(mb.inventoryTransactions))
	for _, transaction := range mb.inventoryTransactions {
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (mb *InMemoryBackend) ListLocations(ctx context.Context) ([]*Location, error) {
	locations := make([]*Location, 0, len(mb.locations))
	for _, location := range mb.locations {
		locations = append(locations, location)
	}
	return locations, nil
}

func (mb *InMemoryBackend) ListLocationInventory(ctx context.Context, id string) ([]*Inventory, error) {
	var inventories []*Inventory
	if locInvs, ok := mb.inventoryByLocationByItemIndex[id]; ok {
		for _, inventory := range locInvs {
			inventories = append(inventories, inventory)
		}
	}
	return inventories, nil
}

func (mb *InMemoryBackend) ListLocationInventoryTransactions(ctx context.Context, id string) ([]*InventoryTransaction, error) {
	var txns []*InventoryTransaction
	for _, txn := range mb.inventoryTransactions {
		if txn.LocationId == id {
			txns = append(txns, txn)
		}
	}
	return txns, nil
}

func (mb *InMemoryBackend) NewItem(ctx context.Context, inputItem *Item) (*Item, error) {
	item := &Item{}
	*item = *inputItem
	item.Id = uuid.New().String()
	mb.items[item.Id] = item
	return item, nil
}

// lookupInventory returns the inventory associated with the transaction
func (mb *InMemoryBackend) lookupInventory(itemID, locID string) *Inventory {
	// item and/or location may not have any inventory yet. Create index entries as needed.
	if _, found := mb.inventoryByItemByLocationIndex[itemID]; !found {
		mb.inventoryByItemByLocationIndex[itemID] = make(map[string]*Inventory)
	}
	if _, found := mb.inventoryByLocationByItemIndex[locID]; !found {
		mb.inventoryByLocationByItemIndex[locID] = make(map[string]*Inventory)
	}

	// Make sure the two indices return the same inventory (which could be nil)
	inv, exists := mb.inventoryByItemByLocationIndex[itemID][locID]
	locItemInv, _ := mb.inventoryByLocationByItemIndex[locID][itemID]
	if locItemInv != inv {
		msg := fmt.Sprintf("[item][location] index returned %p and [location][item] index returned %p", inv, locItemInv)
		log.Panicf("inventory data inconsistent for item %q and location %q: %s", itemID, locID, msg)
	}

	// Create an inventory entry if needed
	if !exists {
		inv = &Inventory{ItemId: itemID, LocationId: locID}
		mb.inventoryByItemByLocationIndex[itemID][locID] = inv
		mb.inventoryByLocationByItemIndex[locID][itemID] = inv
	}
	return inv
}

func (mb *InMemoryBackend) NewInventoryTransaction(ctx context.Context, inputTxn *InventoryTransaction) (*InventoryTransaction, error) {
	if _, ok := mb.items[inputTxn.ItemId]; !ok {
		return nil, ItemNotFound(inputTxn.ItemId)
	}
	if _, ok := mb.locations[inputTxn.LocationId]; !ok {
		return nil, LocationNotFound(inputTxn.LocationId)
	}

	transaction := &InventoryTransaction{}
	*transaction = *inputTxn
	transaction.Id = uuid.New().String()
	inv := mb.lookupInventory(transaction.ItemId, transaction.LocationId)
	if err := inv.applyTransaction(transaction); err != nil {
		return nil, err
	}
	mb.inventoryTransactions[transaction.Id] = transaction
	return transaction, nil
}

func (mb *InMemoryBackend) NewLocation(ctx context.Context, inputLocation *Location) (*Location, error) {
	location := &Location{}
	*location = *inputLocation
	location.Id = uuid.New().String()
	mb.locations[location.Id] = location
	return location, nil
}

func (mb *InMemoryBackend) UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	if _, ok := mb.items[item.Id]; ok {
		*mb.items[item.Id] = *item
		return mb.items[item.Id], nil
	}
	return nil, ItemNotFound(item.Id)
}

func (mb *InMemoryBackend) UpdateLocation(ctx context.Context, location *Location) (*Location, error) {
	if _, ok := mb.locations[location.Id]; ok {
		*mb.locations[location.Id] = *location
		return mb.locations[location.Id], nil
	}
	return nil, LocationNotFound(location.Id)
}
