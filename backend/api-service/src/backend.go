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
	"time"
)

var supportedTransactionActions = []string{"ADD", "REMOVE", "RECOUNT"}

// applyTransaction applies the transaction to the inventory
func (i *Inventory) applyTransaction(txn *InventoryTransaction) error {
	switch txn.Action {
	case "ADD":
		i.Count += txn.Count
	case "REMOVE":
		i.Count -= txn.Count
	case "RECOUNT":
		i.Count = txn.Count
	default:
		return fmt.Errorf("unknown action: %s", txn.Action)
	}
	txn.Timestamp = time.Now()
	i.LastUpdated = txn.Timestamp
	return nil
}

type DatabaseBackend interface {
	DeleteItem(ctx context.Context, id string) error
	DeleteLocation(ctx context.Context, id string) error
	DeleteAlert(ctx context.Context, id string) error

	GetInventoryTransaction(ctx context.Context, id string) (*InventoryTransaction, error)
	GetItem(ctx context.Context, id string) (*Item, error)
	GetLocation(ctx context.Context, id string) (*Location, error)

	ListItems(ctx context.Context) ([]*Item, error)
	ListItemInventory(ctx context.Context, itemId string) ([]*Inventory, error)
	ListItemInventoryTransactions(ctx context.Context, itemId string) ([]*InventoryTransaction, error)
	ListInventoryTransactions(ctx context.Context) ([]*InventoryTransaction, error)
	ListLocations(ctx context.Context) ([]*Location, error)
	ListLocationInventory(ctx context.Context, locationId string) ([]*Inventory, error)
	ListLocationInventoryTransactions(ctx context.Context, locationId string) ([]*InventoryTransaction, error)
	ListAlerts(ctx context.Context) ([]*Alert, error)

	NewItem(ctx context.Context, item *Item) (*Item, error)
	NewInventoryTransaction(ctx context.Context, transaction *InventoryTransaction) (*InventoryTransaction, error)
	NewLocation(ctx context.Context, location *Location) (*Location, error)
	NewAlert(ctx context.Context, alert *Alert) (*Alert, error)

	UpdateItem(ctx context.Context, item *Item) (*Item, error)
	UpdateLocation(ctx context.Context, location *Location) (*Location, error)

	lookupInventory(ctx context.Context, itemID, locationID string) (*Inventory, error)
}
