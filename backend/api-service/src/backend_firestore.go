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
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FirestoreBackend is a backend that talks to Firestore and implements the
// Backend interface
type FirestoreBackend struct {
	projectID string
}

func NewFirestoreBackend(projectID string) *FirestoreBackend {
	return &FirestoreBackend{projectID}
}

func (fb *FirestoreBackend) NewClient(ctx context.Context) (*firestore.Client, error) {
	// Return a new client
	client, err := firestore.NewClient(ctx, fb.projectID)
	if err != nil {
		log.Printf("error creating firestore client: %v", err)
	}
	return client, err
}

func (fb *FirestoreBackend) deleteDoc(ctx context.Context, path, id string) error {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection(path).Doc(id).Delete(ctx, firestore.Exists)
	if err != nil && status.Code(err) == codes.NotFound {
		return &ResourceNotFound{collection: path, id: id}
	}
	return err
}

func (fb *FirestoreBackend) DeleteItem(ctx context.Context, id string) error {
	return fb.deleteDoc(ctx, "items", id)
}

func (fb *FirestoreBackend) DeleteLocation(ctx context.Context, id string) error {
	return fb.deleteDoc(ctx, "locations", id)
}

func (fb *FirestoreBackend) DeleteAlert(ctx context.Context, id string) error {
	return fb.deleteDoc(ctx, "alerts", id)
}

func (fb *FirestoreBackend) getDoc(ctx context.Context, path, id string) (*firestore.DocumentSnapshot, error) {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := client.Collection(path).Doc(id).Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return nil, &ResourceNotFound{collection: path, id: id}
	}
	return doc, err
}

func (fb *FirestoreBackend) GetInventoryTransaction(ctx context.Context, id string) (*InventoryTransaction, error) {
	doc, err := fb.getDoc(ctx, "inventoryTransactions", id)
	if err != nil {
		return nil, err
	}
	transaction := &InventoryTransaction{}
	err = doc.DataTo(transaction)
	return transaction, err
}

func (fb *FirestoreBackend) GetItem(ctx context.Context, id string) (*Item, error) {
	doc, err := fb.getDoc(ctx, "items", id)
	if err != nil {
		return nil, err
	}
	item := &Item{}
	err = doc.DataTo(item)
	return item, err
}

func (fb *FirestoreBackend) GetLocation(ctx context.Context, id string) (*Location, error) {
	doc, err := fb.getDoc(ctx, "locations", id)
	if err != nil {
		return nil, err
	}
	location := &Location{}
	err = doc.DataTo(location)
	return location, err
}

type queryFilter struct {
	path  string
	op    string
	value interface{}
}

func (fb *FirestoreBackend) listDocs(ctx context.Context, path string, filters ...queryFilter) ([]*firestore.DocumentSnapshot, error) {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	q := client.Collection(path).Query
	for _, f := range filters {
		q = q.Where(f.path, f.op, f.value)
	}

	return q.Documents(ctx).GetAll()
}

func (fb *FirestoreBackend) ListItems(ctx context.Context) ([]*Item, error) {
	docs, err := fb.listDocs(ctx, "items")
	if err != nil {
		return nil, err
	}

	items := make([]*Item, 0, len(docs))
	for _, doc := range docs {
		item := &Item{}
		if err = doc.DataTo(item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (fb *FirestoreBackend) ListLocations(ctx context.Context) ([]*Location, error) {
	docs, err := fb.listDocs(ctx, "locations")
	if err != nil {
		return nil, err
	}

	locations := make([]*Location, 0, len(docs))
	for _, doc := range docs {
		location := &Location{}
		if err = doc.DataTo(location); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (fb *FirestoreBackend) ListAlerts(ctx context.Context) ([]*Alert, error) {
	docs, err := fb.listDocs(ctx, "alerts")
	if err != nil {
		return nil, err
	}

	alerts := make([]*Alert, 0, len(docs))
	for _, doc := range docs {
		alert := &Alert{}
		if err = doc.DataTo(alert); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func (fb *FirestoreBackend) listInventories(ctx context.Context, filters ...queryFilter) ([]*Inventory, error) {
	docs, err := fb.listDocs(ctx, "inventories", filters...)
	if err != nil {
		return nil, err
	}

	invs := make([]*Inventory, 0, len(docs))
	for _, doc := range docs {
		inv := &Inventory{}
		if err = doc.DataTo(inv); err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, nil
}

func (fb *FirestoreBackend) listInventoryTransactions(ctx context.Context, filters ...queryFilter) ([]*InventoryTransaction, error) {
	docs, err := fb.listDocs(ctx, "inventoryTransactions", filters...)
	if err != nil {
		return nil, err
	}

	txns := make([]*InventoryTransaction, 0, len(docs))
	for _, doc := range docs {
		txn := &InventoryTransaction{}
		if err = doc.DataTo(txn); err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}
	return txns, nil
}

func (fb *FirestoreBackend) ListItemInventory(ctx context.Context, itemId string) ([]*Inventory, error) {
	return fb.listInventories(ctx, queryFilter{"ItemId", "==", itemId})
}

func (fb *FirestoreBackend) ListLocationInventory(ctx context.Context, locationId string) ([]*Inventory, error) {
	return fb.listInventories(ctx, queryFilter{"LocationId", "==", locationId})
}

func (fb *FirestoreBackend) ListInventoryTransactions(ctx context.Context) ([]*InventoryTransaction, error) {
	return fb.listInventoryTransactions(ctx)
}

func (fb *FirestoreBackend) ListItemInventoryTransactions(ctx context.Context, itemId string) ([]*InventoryTransaction, error) {
	return fb.listInventoryTransactions(ctx, queryFilter{"ItemId", "==", itemId})
}

func (fb *FirestoreBackend) ListLocationInventoryTransactions(ctx context.Context, locationId string) ([]*InventoryTransaction, error) {
	return fb.listInventoryTransactions(ctx, queryFilter{"LocationId", "==", locationId})
}

func (fb *FirestoreBackend) NewInventoryTransaction(ctx context.Context, invTxn *InventoryTransaction) (*InventoryTransaction, error) {
	itemId, locId := invTxn.ItemId, invTxn.LocationId
	if _, err := fb.getDoc(ctx, "items", itemId); err != nil {
		return nil, err
	}
	if _, err := fb.getDoc(ctx, "locations", locId); err != nil {
		return nil, err
	}
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// Lookup the inventory id
	invs := client.Collection("inventories")
	var invRef *firestore.DocumentRef

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Find the inventory
		q := invs.Where("ItemId", "==", itemId).Where("LocationId", "==", locId)
		docs, err := tx.Documents(q).GetAll()
		if err != nil {
			return fmt.Errorf("error querying inventories collection: %v", err)
		}

		if len(docs) > 2 {
			log.Panicf("Found multiple inventories for item %q and location %q", invTxn.ItemId, invTxn.LocationId)
		}

		if len(docs) == 0 {
			inv := &Inventory{ItemId: itemId, LocationId: locId, LastUpdated: time.Now()}
			invRef = invs.NewDoc()
			return tx.Create(invRef, inv)
		}
		invRef = docs[0].Ref
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Fetch the inventory
		inv := Inventory{}
		doc, err := tx.Get(invRef)
		if err != nil {
			return err
		}
		if err = doc.DataTo(&inv); err != nil {
			return err
		}

		// Update the inventory
		if err := inv.applyTransaction(invTxn); err != nil {
			return err
		}
		if err := tx.Set(invRef, inv); err != nil {
			return err
		}

		// Create the inventory transaction itself
		dref := client.Collection("inventoryTransactions").NewDoc()
		invTxn.Id = dref.ID
		return tx.Create(dref, invTxn)
	})

	return invTxn, err
}

func (fb *FirestoreBackend) NewItem(ctx context.Context, item *Item) (*Item, error) {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	dref := client.Collection("items").NewDoc()
	item.Id = dref.ID
	_, err = dref.Create(ctx, item)
	return item, err
}

func (fb *FirestoreBackend) NewLocation(ctx context.Context, location *Location) (*Location, error) {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	dref := client.Collection("locations").NewDoc()
	location.Id = dref.ID
	_, err = dref.Create(ctx, location)
	return location, err
}

func (fb *FirestoreBackend) NewAlert(ctx context.Context, alert *Alert) (*Alert, error) {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	dref := client.Collection("alerts").NewDoc()
	alert.Id = dref.ID
	_, err = dref.Create(ctx, alert)
	return alert, err
}

func (fb *FirestoreBackend) update(ctx context.Context, path, id string, value interface{}) error {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return err
	}
	dref := client.Collection(path).Doc(id)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(dref); err != nil {
			if status.Code(err) == codes.NotFound {
				return &ResourceNotFound{collection: path, id: id}
			}
			return err
		}
		return tx.Set(dref, value)
	})

	return err
}

func (fb *FirestoreBackend) UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	err := fb.update(ctx, "items", item.Id, item)
	return item, err
}

func (fb *FirestoreBackend) UpdateLocation(ctx context.Context, location *Location) (*Location, error) {
	err := fb.update(ctx, "locations", location.Id, location)
	return location, err
}

func (fb *FirestoreBackend) lookupInventory(ctx context.Context, itemID, locationID string) (*Inventory, error) {
	client, err := fb.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// Lookup the inventory id
	invs := client.Collection("inventories")
	var inv *Inventory
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Find the inventory
		q := invs.Where("ItemId", "==", itemID).Where("LocationId", "==", locationID)
		docs, err := tx.Documents(q).GetAll()
		if err != nil {
			return err
		}

		if len(docs) > 2 {
			log.Panicf("Found multiple inventories for item %q and location %q", itemID, locationID)
		}

		if len(docs) == 0 {
			inv = &Inventory{ItemId: itemID, LocationId: locationID, LastUpdated: time.Now()}
			return tx.Create(invs.NewDoc(), inv)
		}

		return docs[0].DataTo(inv)
	})
	return inv, err
}
