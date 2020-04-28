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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type initialBackendState struct {
	items                 map[string]*Item
	inventories           map[string]*Inventory
	inventoryTransactions map[string]*InventoryTransaction
	locations             map[string]*Location
}

type backendTester struct {
	resetBackend func(*testing.T) DatabaseBackend
	initBackend  func(*testing.T, initialBackendState) DatabaseBackend
}

func modelLess(a, b interface{}) bool {
	switch l := a.(type) {
	case *Inventory:
		r := b.(*Inventory)
		return l.ItemId < r.ItemId || l.ItemId == r.ItemId && l.LocationId < r.LocationId
	case *InventoryTransaction:
		return l.Id < b.(*InventoryTransaction).Id
	case *Item:
		return l.Id < b.(*Item).Id
	case *Location:
		return l.Id < b.(*Location).Id
	default:
		panic(fmt.Sprintf("unknown type: %v", a))
	}
}

func (bt *backendTester) testDeleteItem(t *testing.T) {
	ctx := context.Background()
	id := "item-id"
	item := Item{
		Id:          id,
		Name:        "name",
		Description: "description",
	}
	backend := bt.initBackend(t, initialBackendState{items: map[string]*Item{id: &item}})

	err := backend.DeleteItem(ctx, id)

	if err != nil {
		t.Fatalf("DeleteItem(%q) = %v, want nil", id, err)
	}
	if _, err = backend.GetItem(ctx, id); err == nil {
		t.Fatalf("after DeleteItem(%q), GetItem(%q) succeeded, want error", id, id)
	}
}

func (bt *backendTester) testDeleteItemNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	backend := bt.resetBackend(t)
	want := ItemNotFound(id)

	err := backend.DeleteItem(ctx, id)

	if err == nil {
		t.Fatalf("DeleteItem(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("DeleteItem(%q) returned %v, want %v", id, err, want)
	}
}

func (bt *backendTester) testDeleteLocation(t *testing.T) {
	ctx := context.Background()
	id := "location-id"
	location := Location{
		Id:        id,
		Name:      "name",
		Warehouse: "warehouse",
	}
	backend := bt.initBackend(t, initialBackendState{locations: map[string]*Location{id: &location}})

	err := backend.DeleteLocation(ctx, id)

	if err != nil {
		t.Fatalf("DeleteLocation(%v) = %v, want nil", id, err)
	}
	if _, err = backend.GetLocation(ctx, id); err == nil {
		t.Errorf("after backend.DeleteLocation(%v), backend.GetLocation(%v) is not nil", id, id)
	}
}

func (bt *backendTester) testDeleteLocationNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	backend := bt.resetBackend(t)
	want := LocationNotFound(id)

	err := backend.DeleteLocation(ctx, id)

	if err == nil {
		t.Fatalf("DeleteLocation(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("DeleteLocation(%q) returned %v, want %v", id, err, want)
	}
}

func (bt *backendTester) testGetInventoryTransaction(t *testing.T) {
	ctx := context.Background()
	id := "txn-id"
	txn := InventoryTransaction{
		Id:         id,
		ItemId:     "item-id",
		LocationId: "location-id",
		Action:     "action",
		Count:      100,
		Note:       "some note",
		Timestamp:  time.Now(),
		CreatedBy:  "someone",
	}
	backend := bt.initBackend(t, initialBackendState{inventoryTransactions: map[string]*InventoryTransaction{id: &txn}})

	got, err := backend.GetInventoryTransaction(ctx, id)

	if err != nil || !cmp.Equal(got, &txn, cmpopts.EquateApproxTime(time.Millisecond)) {
		t.Errorf("GetInventoryTransaction(%v) = %v, %v want %v, nil", id, got, err, &txn)
	}
}

func (bt *backendTester) testGetInventoryTransactionNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	backend := bt.resetBackend(t)
	want := InventoryTransactionNotFound(id)

	_, err := backend.GetInventoryTransaction(ctx, id)

	if err == nil {
		t.Fatalf("GetInventoryTransaction(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("GetInventoryTransaction(%q) returned %v, want %v", id, err, want)
	}
}

func (bt *backendTester) testGetItem(t *testing.T) {
	ctx := context.Background()
	id := "item-id"
	item := Item{
		Id:          id,
		Name:        "name",
		Description: "description",
	}
	backend := bt.initBackend(t, initialBackendState{items: map[string]*Item{item.Id: &item}})
	got, err := backend.GetItem(ctx, id)

	if err != nil || *got != item {
		t.Errorf("GetItem(%v) = %v, %v want %v, nil", id, got, err, &item)
	}
}

func (bt *backendTester) testGetItemNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	backend := bt.resetBackend(t)
	want := ItemNotFound(id)

	_, err := backend.GetItem(ctx, id)

	if err == nil {
		t.Fatalf("GetItem(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("GetItem(%q) returned %v, want %v", id, err, want)
	}
}

func (bt *backendTester) testGetLocation(t *testing.T) {
	ctx := context.Background()
	id := "location-id"
	location := Location{
		Id:        id,
		Name:      "name",
		Warehouse: "warehouse",
	}
	backend := bt.initBackend(t, initialBackendState{locations: map[string]*Location{location.Id: &location}})
	got, err := backend.GetLocation(ctx, id)

	if err != nil || *got != location {
		t.Errorf("GetLocation(%v) = %v, %v want %v, nil", id, got, err, &location)
	}
}

func (bt *backendTester) testGetLocationNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	backend := bt.resetBackend(t)
	want := LocationNotFound(id)

	_, err := backend.GetLocation(ctx, id)

	if err == nil {
		t.Fatalf("GetLocation(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("GetLocation(%q) returned %v, want %v", id, err, want)
	}
}

func (bt *backendTester) testListInventoryTransactions(t *testing.T) {
	txn1, txn2 := InventoryTransaction{Id: "txn1-id"}, InventoryTransaction{Id: "txn2-id"}
	cases := []struct {
		desc string
		init initialBackendState
		want []*InventoryTransaction
	}{
		{
			desc: "no txns",
			init: initialBackendState{},
			want: []*InventoryTransaction{},
		},
		{
			desc: "single txn",
			init: initialBackendState{inventoryTransactions: map[string]*InventoryTransaction{txn1.Id: &txn1}},
			want: []*InventoryTransaction{&txn1},
		},
		{
			desc: "multiple txns",
			init: initialBackendState{inventoryTransactions: map[string]*InventoryTransaction{txn1.Id: &txn1, txn2.Id: &txn2}},
			want: []*InventoryTransaction{&txn1, &txn2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			backend := bt.initBackend(t, tc.init)

			got, err := backend.ListInventoryTransactions(ctx)

			if err != nil {
				t.Fatalf("ListInventoryTransactions() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess)) {
				t.Errorf("ListInventoryTransactions() = %v want %v", got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testListItemInventory(t *testing.T) {
	item := Item{Id: "item-id"}
	location1, location2 := Location{Id: "loc-id-1"}, Location{Id: "loc-id-2"}
	inventory1 := Inventory{
		ItemId:      item.Id,
		LocationId:  location1.Id,
		Count:       20,
		LastUpdated: time.Now(),
	}
	inventory2 := Inventory{
		ItemId:      item.Id,
		LocationId:  location2.Id,
		Count:       50,
		LastUpdated: time.Now(),
	}
	cases := []struct {
		desc string
		init initialBackendState
		id   string
		want []*Inventory
	}{
		{
			desc: "no inventory",
			init: initialBackendState{},
			id:   item.Id,
			want: []*Inventory{},
		},
		{
			desc: "single location",
			init: initialBackendState{
				inventories: map[string]*Inventory{
					"inventory1": &inventory1,
				},
			},
			id:   item.Id,
			want: []*Inventory{&inventory1},
		},
		{
			desc: "multiple locations",
			init: initialBackendState{
				inventories: map[string]*Inventory{
					"inventory1": &inventory1,
					"inventory2": &inventory2,
				},
			},
			id:   item.Id,
			want: []*Inventory{&inventory1, &inventory2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			backend := bt.initBackend(t, tc.init)

			got, err := backend.ListItemInventory(ctx, tc.id)

			if err != nil {
				t.Fatalf("ListItemInventory(%v) returned unexpected err: %v", tc.id, err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateApproxTime(time.Millisecond)) {
				t.Errorf("ListItemInventory(%v) = %v want %v", tc.id, got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testListItemInventoryTransactions(t *testing.T) {
	item0, item1, item2 := Item{Id: "item0-id"}, Item{Id: "item1-id"}, Item{Id: "item2-id"}
	item1Txn1 := InventoryTransaction{
		Id:         "txn-id-1",
		ItemId:     item1.Id,
		LocationId: "location-id",
		Count:      100,
		Action:     "add",
	}
	item2Txn1 := InventoryTransaction{
		Id:         "txn-id-2",
		ItemId:     item2.Id,
		LocationId: "location-id",
		Count:      100,
		Action:     "update",
	}
	item2Txn2 := InventoryTransaction{
		Id:         "txn-id-3",
		ItemId:     item2.Id,
		LocationId: "location-id",
		Count:      100,
		Action:     "delete",
	}
	backend := bt.initBackend(t, initialBackendState{
		inventoryTransactions: map[string]*InventoryTransaction{
			item1Txn1.Id: &item1Txn1,
			item2Txn1.Id: &item2Txn1,
			item2Txn2.Id: &item2Txn2,
		},
	})
	cases := []struct {
		desc string
		id   string
		want []*InventoryTransaction
	}{
		{
			desc: "no matching transactions",
			id:   item0.Id,
			want: []*InventoryTransaction{},
		},
		{
			desc: "single matching transaction",
			id:   item1.Id,
			want: []*InventoryTransaction{&item1Txn1},
		},
		{
			desc: "multiple matching transactions",
			id:   item2.Id,
			want: []*InventoryTransaction{&item2Txn1, &item2Txn2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := backend.ListItemInventoryTransactions(ctx, tc.id)

			if err != nil {
				t.Fatalf("ListItemInventoryTransactions() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess)) {
				t.Errorf("ListItemInventoryTransactions() = %v want %v", got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testListItems(t *testing.T) {
	item1, item2 := Item{Id: "item1-id"}, Item{Id: "item2-id"}
	cases := []struct {
		desc string
		init initialBackendState
		want []*Item
	}{
		{
			desc: "no items",
			init: initialBackendState{},
			want: []*Item{},
		},
		{
			desc: "single item",
			init: initialBackendState{items: map[string]*Item{item1.Id: &item1}},
			want: []*Item{&item1},
		},
		{
			desc: "multiple items",
			init: initialBackendState{items: map[string]*Item{item1.Id: &item1, item2.Id: &item2}},
			want: []*Item{&item1, &item2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			backend := bt.initBackend(t, tc.init)

			got, err := backend.ListItems(ctx)

			if err != nil {
				t.Fatalf("ListItems() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess)) {
				t.Errorf("ListItems() = %v want %v", got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testListLocationInventory(t *testing.T) {
	loc0, loc1, loc2 := Location{Id: "loc0-id"}, Location{Id: "loc1-id"}, Location{Id: "loc2-id"}
	loc1Inv1 := Inventory{
		LocationId:  loc1.Id,
		ItemId:      "item-id",
		Count:       20,
		LastUpdated: time.Now(),
	}
	loc2Inv1 := Inventory{
		LocationId:  loc2.Id,
		ItemId:      "item-id",
		Count:       20,
		LastUpdated: time.Now(),
	}
	loc2Inv2 := Inventory{
		LocationId:  loc2.Id,
		ItemId:      "different-item-id",
		Count:       50,
		LastUpdated: time.Now(),
	}
	backend := bt.initBackend(t, initialBackendState{
		inventories: map[string]*Inventory{
			"loc1Inv1-id": &loc1Inv1,
			"loc2Inv1-id": &loc2Inv1,
			"loc2Inv2-id": &loc2Inv2,
		},
	})
	cases := []struct {
		desc string
		id   string
		want []*Inventory
	}{
		{
			desc: "no inventory",
			id:   loc0.Id,
			want: []*Inventory{},
		},
		{
			desc: "single item",
			id:   loc1.Id,
			want: []*Inventory{&loc1Inv1},
		},
		{
			desc: "multiple items",
			id:   loc2.Id,
			want: []*Inventory{&loc2Inv1, &loc2Inv2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := backend.ListLocationInventory(ctx, tc.id)

			if err != nil {
				t.Fatalf("ListLocationInventory(%v) returned unexpected err: %v", tc.id, err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateApproxTime(time.Millisecond)) {
				t.Errorf("ListLocationInventory(%v) = %v want %v", tc.id, got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testListLocationInventoryTransactions(t *testing.T) {
	loc0, loc1, loc2 := Location{Id: "loc0-id"}, Location{Id: "loc1-id"}, Location{Id: "loc2-id"}
	loc1Txn1 := InventoryTransaction{
		Id:         "txn-id-1",
		LocationId: loc1.Id,
		ItemId:     "item-id",
		Count:      100,
		Action:     "add",
	}
	loc2Txn1 := InventoryTransaction{
		Id:         "txn-id-2",
		LocationId: loc2.Id,
		ItemId:     "item-id",
		Count:      100,
		Action:     "update",
	}
	loc2Txn2 := InventoryTransaction{
		Id:         "txn-id-3",
		LocationId: loc2.Id,
		ItemId:     "item-id",
		Count:      100,
		Action:     "delete",
	}
	backend := bt.initBackend(t, initialBackendState{
		inventoryTransactions: map[string]*InventoryTransaction{
			"loc1Txn1-id": &loc1Txn1,
			"loc2Txn1-id": &loc2Txn1,
			"loc2Txn2-id": &loc2Txn2,
		},
	})
	cases := []struct {
		desc string
		id   string
		want []*InventoryTransaction
	}{
		{
			desc: "no matching transactions",
			id:   loc0.Id,
			want: []*InventoryTransaction{},
		},
		{
			desc: "single matching transaction",
			id:   loc1.Id,
			want: []*InventoryTransaction{&loc1Txn1},
		},
		{
			desc: "multiple matching transactions",
			id:   loc2.Id,
			want: []*InventoryTransaction{&loc2Txn1, &loc2Txn2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := backend.ListLocationInventoryTransactions(ctx, tc.id)

			if err != nil {
				t.Fatalf("ListLocationInventoryTransactions() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateApproxTime(time.Millisecond)) {
				t.Errorf("ListLocationInventoryTransactions() = %v want %v", got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testListLocations(t *testing.T) {
	location1, location2 := Location{Id: "location1-id"}, Location{Id: "location2-id"}
	cases := []struct {
		desc string
		init initialBackendState
		want []*Location
	}{
		{
			desc: "no locations",
			init: initialBackendState{},
			want: []*Location{},
		},
		{
			desc: "single location",
			init: initialBackendState{locations: map[string]*Location{location1.Id: &location1}},
			want: []*Location{&location1},
		},
		{
			desc: "multiple locations",
			init: initialBackendState{locations: map[string]*Location{location1.Id: &location1, location2.Id: &location2}},
			want: []*Location{&location1, &location2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			backend := bt.initBackend(t, tc.init)

			got, err := backend.ListLocations(ctx)

			if err != nil {
				t.Fatalf("ListLocations() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess)) {
				t.Errorf("ListLocations() = %v want %v", got, tc.want)
			}
		})
	}
}

func (bt *backendTester) testNewInventoryTransaction(t *testing.T) {
	// properties of a
	stockedItem := Item{Id: "stocked-item"}
	stockedLoc := Location{Id: "stocked-location"}
	const stockedInitCount = 100

	newItem := Item{Id: "new-item"}
	newLoc := Location{Id: "new-location"}
	cases := []struct {
		desc      string
		txn       *InventoryTransaction
		wantCount int64
	}{
		{
			desc: "add to existing inventory",
			txn: &InventoryTransaction{
				ItemId:     stockedItem.Id,
				LocationId: stockedLoc.Id,
				Action:     "ADD",
				Count:      20,
			},
			wantCount: stockedInitCount + 20,
		},
		{
			desc: "add new inventory",
			txn: &InventoryTransaction{
				ItemId:     newItem.Id,
				LocationId: newLoc.Id,
				Action:     "ADD",
				Count:      20,
			},
			wantCount: 20,
		},
		{
			desc: "recount existing inventory",
			txn: &InventoryTransaction{
				ItemId:     stockedItem.Id,
				LocationId: stockedLoc.Id,
				Action:     "RECOUNT",
				Count:      20,
			},
			wantCount: 20,
		},
		{
			desc: "recount existing inventory to 0",
			txn: &InventoryTransaction{
				ItemId:     stockedItem.Id,
				LocationId: stockedLoc.Id,
				Action:     "RECOUNT",
				Count:      0,
			},
			wantCount: 0,
		},
		{
			desc: "recount new inventory",
			txn: &InventoryTransaction{
				ItemId:     newItem.Id,
				LocationId: newLoc.Id,
				Action:     "RECOUNT",
				Count:      20,
			},
			wantCount: 20,
		},
		{
			desc: "remove from existing inventory",
			txn: &InventoryTransaction{
				ItemId:     stockedItem.Id,
				LocationId: stockedLoc.Id,
				Action:     "REMOVE",
				Count:      20,
			},
			wantCount: stockedInitCount - 20,
		},
		{
			desc: "remove from new inventory",
			txn: &InventoryTransaction{
				ItemId:     newItem.Id,
				LocationId: newLoc.Id,
				Action:     "REMOVE",
				Count:      20,
			},
			wantCount: -20,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		stockedInv := Inventory{
			LocationId:  stockedLoc.Id,
			ItemId:      stockedItem.Id,
			Count:       stockedInitCount,
			LastUpdated: time.Now(),
		}
		backend := bt.initBackend(t, initialBackendState{
			inventories: map[string]*Inventory{"stockedInv-id": &stockedInv},
			items:       map[string]*Item{newItem.Id: &newItem, stockedItem.Id: &stockedItem},
			locations:   map[string]*Location{newLoc.Id: &newLoc, stockedLoc.Id: &stockedLoc},
		})
		t.Run(tc.desc, func(t *testing.T) {

			got, err := backend.NewInventoryTransaction(ctx, tc.txn)
			if err != nil {
				t.Fatalf("NewInventoryTransaction(%v) returned unexpected err: %v", tc.txn, err)
			}
			if got.Id == "" {
				t.Errorf("NewInventoryTransaction(%v) did not generate InventoryTransaction.Id", tc.txn)
			}
			if !cmp.Equal(got, tc.txn, cmpopts.IgnoreFields(InventoryTransaction{}, "Id", "Timestamp")) {
				t.Errorf("NewInventoryTransaction(%v) = %v want %v (ignoring Id field)", tc.txn, got, tc.txn)
			}
			if v, _ := backend.GetInventoryTransaction(ctx, got.Id); v != got {
				t.Errorf("after backend.NewInventoryTransaction(%v), backend.GetInventoryTransaction(%v) = %v want %v", tc.txn, got.Id, v, got)
			}

			wantInv := &Inventory{ItemId: tc.txn.ItemId, LocationId: tc.txn.LocationId, Count: tc.wantCount}

			inv, err := backend.lookupInventory(ctx, tc.txn.ItemId, tc.txn.LocationId)
			if err != nil {
				t.Errorf("getting inventory for item: %v, location: %v, produced error: %v", tc.txn.ItemId, tc.txn.LocationId, err)
			}
			if !cmp.Equal(inv, wantInv, cmpopts.IgnoreFields(Inventory{}, "LastUpdated")) {
				t.Errorf("NewInventoryTransaction(%v) produced unexpected inventory state", tc.txn)
				t.Errorf("[item %q, location %q] inventory = %v want %v", tc.txn.ItemId, tc.txn.LocationId, inv, wantInv)
			}
		})
	}
}

func (bt *backendTester) testNewInventoryTransactionNotFoundErrors(t *testing.T) {
	existingLoc := Location{Id: "existing-location-id"}
	existingItem := Item{Id: "existing-item-id"}
	backend := bt.initBackend(t, initialBackendState{
		items: map[string]*Item{
			existingItem.Id: &existingItem,
		},
		locations: map[string]*Location{
			existingLoc.Id: &existingLoc,
		},
	})
	cases := []struct {
		desc string
		txn  *InventoryTransaction
		want *ResourceNotFound
	}{
		{
			desc: "item not found",
			txn: &InventoryTransaction{
				ItemId:     "bad-item-id",
				LocationId: existingLoc.Id,
			},
			want: ItemNotFound("bad-item-id"),
		},
		{
			desc: "location not found",
			txn: &InventoryTransaction{
				ItemId:     existingItem.Id,
				LocationId: "bad-loc-id",
			},
			want: LocationNotFound("bad-loc-id"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			_, err := backend.NewInventoryTransaction(ctx, tc.txn)
			if err == nil {
				t.Fatalf("NewInventoryTransaction(%q) succeeded, want error", tc.txn)
			}
			if nf, ok := err.(*ResourceNotFound); !ok || nf.id != tc.want.id || nf.collection != tc.want.collection {
				t.Errorf("NewInventoryTransaction(%q) returned %v, want %v", tc.txn, err, tc.want)
			}
		})
	}
}

func (bt *backendTester) testNewItem(t *testing.T) {
	ctx := context.Background()
	backend := bt.resetBackend(t)
	item := Item{
		Id:          "id-to-be-replaced-by-uuid",
		Name:        "name",
		Description: "description",
	}

	got, err := backend.NewItem(ctx, &item)
	if err != nil {
		t.Fatalf("NewItem(%v) returned unexpected err: %v", item, err)
	}
	if got.Id == "" {
		t.Errorf("NewItem(%v) did not generate Item.Id", item)
	}
	if !cmp.Equal(got, &item, cmpopts.IgnoreFields(Item{}, "Id")) {
		t.Errorf("NewItem(%v) = %v want %v (ignoring Id field)", item, got, item)
	}
	if v, _ := backend.GetItem(ctx, got.Id); !cmp.Equal(v, got) {
		t.Errorf("after backend.NewItem(%v), backend.GetItem(%v) = %v want %v", item, got.Id, v, got)
	}
}

func (bt *backendTester) testNewLocation(t *testing.T) {
	ctx := context.Background()
	backend := bt.resetBackend(t)
	location := Location{
		Name:      "name",
		Warehouse: "warehouse",
	}

	got, err := backend.NewLocation(ctx, &location)

	if err != nil {
		t.Fatalf("NewLocation(%v) returned unexpected err: %v", location, err)
	}
	if got.Id == "" {
		t.Errorf("NewLocation(%v) did not generate Location.Id", location)
	}
	if !cmp.Equal(got, &location, cmpopts.IgnoreFields(Location{}, "Id")) {
		t.Errorf("NewLocation(%v) = %v want %v (ignoring Id field)", location, got, location)
	}
	if v, _ := backend.GetLocation(ctx, got.Id); !cmp.Equal(v, got) {
		t.Errorf("after backend.NewLocation(%v), backend.GetLocation(%v) = %v want %v", location, got.Id, v, got)
	}
}

func (bt *backendTester) testUpdateItem(t *testing.T) {
	ctx := context.Background()
	id := "item-id"
	backend := bt.initBackend(t, initialBackendState{
		items: map[string]*Item{
			id: {
				Id:          id,
				Name:        "old-name",
				Description: "old-description",
			},
		},
	})
	item := &Item{
		Id:          id,
		Name:        "updated-name",
		Description: "updated-description",
	}

	got, err := backend.UpdateItem(ctx, item)

	if err != nil {
		t.Fatalf("UpdateItem(%v) returned unexpected err: %v", item, err)
	}
	if !cmp.Equal(got, item) {
		t.Errorf("UpdateItem(%v) = %v want %v", item, got, item)
	}
	if got, _ := backend.GetItem(ctx, id); !cmp.Equal(got, item) {
		t.Errorf("after backend.UpdateItem(%v), backend.GetItem(%v) = %v want %v", item, id, got, item)
	}
}

func (bt *backendTester) testUpdateItemNotFound(t *testing.T) {
	ctx := context.Background()
	item := &Item{Id: "not-found-id"}
	backend := bt.resetBackend(t)
	want := ItemNotFound(item.Id)

	_, err := backend.UpdateItem(ctx, item)

	if err == nil {
		t.Fatalf("UpdateItem(%q) succeeded, want error", item)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("UpdateItem(%q) returned %v, want %v", item, err, want)
	}
}

func (bt *backendTester) testUpdateLocation(t *testing.T) {
	ctx := context.Background()
	id := "location-id"
	backend := bt.initBackend(t, initialBackendState{
		locations: map[string]*Location{
			id: {
				Id:        id,
				Name:      "old-name",
				Warehouse: "old-warehouse",
			},
		},
	})
	location := &Location{
		Id:        id,
		Name:      "updated-name",
		Warehouse: "updated-warehouse",
	}

	got, err := backend.UpdateLocation(ctx, location)

	if err != nil {
		t.Fatalf("UpdateLocation(%v) returned unexpected err: %v", location, err)
	}
	if !cmp.Equal(got, location) {
		t.Errorf("UpdateLocation(%v) = %v want %v", location, got, location)
	}
	if got, _ = backend.GetLocation(ctx, id); !cmp.Equal(got, location) {
		t.Errorf("after backend.UpdateLocation(%v), backend.GetLocation(%v) = %v want %v", location, id, got, location)
	}
}

func (bt *backendTester) testUpdateLocationNotFound(t *testing.T) {
	location := &Location{Id: "not-found-id"}
	backend := bt.resetBackend(t)
	ctx := context.Background()
	want := LocationNotFound(location.Id)

	_, err := backend.UpdateLocation(ctx, location)

	if err == nil {
		t.Fatalf("UpdateLocation(%q) succeeded, want error", location)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("UpdateLocation(%q) returned %v, want %v", location, err, want)
	}
}
