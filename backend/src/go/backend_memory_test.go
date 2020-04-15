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

func TestDeleteItem(t *testing.T) {
	ctx := context.Background()
	id := "item-id"
	item := Item{
		Id:          id,
		Name:        "name",
		Description: "description",
	}
	mb := &InMemoryBackend{items: map[string]*Item{id: &item}}

	err := mb.DeleteItem(ctx, id)

	if err != nil {
		t.Fatalf("mb.DeleteItem(%v) = %v, want nil", id, err)
	}
	if err == nil && mb.items[id] != nil {
		t.Errorf("after mb.DeleteItem(%v), mb.items[%v] is not nil", id, id)
	}
}

func TestDeleteItemNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	mb := NewInMemoryBackend()
	want := ItemNotFound(id)

	err := mb.DeleteItem(ctx, id)

	if err == nil {
		t.Fatalf("mb.DeleteItem(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.DeleteItem(%q) returned %v, want %v", id, err, want)
	}
}

func TestDeleteLocation(t *testing.T) {
	ctx := context.Background()
	id := "location-id"
	location := Location{
		Id:        id,
		Name:      "name",
		Warehouse: "warehouse",
	}
	mb := &InMemoryBackend{locations: map[string]*Location{id: &location}}

	err := mb.DeleteLocation(ctx, id)

	if err != nil {
		t.Fatalf("mb.DeleteLocation(%v) = %v, want nil", id, err)
	}
	if mb.locations[id] != nil {
		t.Errorf("after mb.DeleteLocation(%v), mb.locations[%v] is not nil", id, id)
	}
}

func TestDeleteLocationNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	mb := NewInMemoryBackend()
	want := LocationNotFound(id)

	err := mb.DeleteLocation(ctx, id)

	if err == nil {
		t.Fatalf("mb.DeleteLocation(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.DeleteLocation(%q) returned %v, want %v", id, err, want)
	}
}

func TestGetInventoryTransaction(t *testing.T) {
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
	mb := InMemoryBackend{inventoryTransactions: map[string]*InventoryTransaction{id: &txn}}

	got, err := mb.GetInventoryTransaction(ctx, id)

	if got != &txn || err != nil {
		t.Errorf("mb.GetInventoryTransaction(%v) = %v, %v want %v, nil", id, got, err, &txn)
	}
}

func TestGetInventoryTransactionNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	mb := NewInMemoryBackend()
	want := InventoryTransactionNotFound(id)

	_, err := mb.GetInventoryTransaction(ctx, id)

	if err == nil {
		t.Fatalf("mb.GetInventoryTransaction(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.GetInventoryTransaction(%q) returned %v, want %v", id, err, want)
	}
}

func TestGetItem(t *testing.T) {
	ctx := context.Background()
	id := "item-id"
	item := Item{
		Id:          id,
		Name:        "name",
		Description: "description",
	}
	mb := InMemoryBackend{items: map[string]*Item{item.Id: &item}}
	got, err := mb.GetItem(ctx, id)

	if got != &item || err != nil {
		t.Errorf("mb.GetItem(%v) = %v, %v want %v, nil", id, got, err, &item)
	}
}

func TestGetItemNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	mb := NewInMemoryBackend()
	want := ItemNotFound(id)

	_, err := mb.GetItem(ctx, id)

	if err == nil {
		t.Fatalf("mb.GetItem(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.GetItem(%q) returned %v, want %v", id, err, want)
	}
}

func TestGetLocation(t *testing.T) {
	ctx := context.Background()
	id := "location-id"
	location := Location{
		Id:        id,
		Name:      "name",
		Warehouse: "warehouse",
	}
	mb := InMemoryBackend{locations: map[string]*Location{location.Id: &location}}
	got, err := mb.GetLocation(ctx, id)

	if got != &location || err != nil {
		t.Errorf("mb.GetLocation(%v) = %v, %v want %v, nil", id, got, err, &location)
	}
}

func TestGetLocationNotFound(t *testing.T) {
	ctx := context.Background()
	id := "not-found-id"
	mb := NewInMemoryBackend()
	want := LocationNotFound(id)

	_, err := mb.GetLocation(ctx, id)

	if err == nil {
		t.Fatalf("mb.GetLocation(%q) succeeded, want error", id)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.GetLocation(%q) returned %v, want %v", id, err, want)
	}
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

func TestListItems(t *testing.T) {
	item1, item2 := Item{Id: "item1-id"}, Item{Id: "item2-id"}
	cases := []struct {
		desc string
		mb   *InMemoryBackend
		want []*Item
	}{
		{
			desc: "no items",
			mb:   NewInMemoryBackend(),
			want: []*Item{},
		},
		{
			desc: "single item",
			mb:   &InMemoryBackend{items: map[string]*Item{item1.Id: &item1}},
			want: []*Item{&item1},
		},
		{
			desc: "multiple items",
			mb:   &InMemoryBackend{items: map[string]*Item{item1.Id: &item1, item2.Id: &item2}},
			want: []*Item{&item1, &item2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.mb.ListItems(ctx)

			if err != nil {
				t.Fatalf("tc.mb.ListItems() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("tc.mb.ListItems() = %v want %v", got, tc.want)
			}
		})
	}
}

func TestListItemInventory(t *testing.T) {
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
		mb   *InMemoryBackend
		id   string
		want []*Inventory
	}{
		{
			desc: "no inventory",
			mb:   NewInMemoryBackend(),
			id:   item.Id,
			want: []*Inventory{},
		},
		{
			desc: "single location",
			mb: &InMemoryBackend{
				inventoryByItemByLocationIndex: map[string]map[string]*Inventory{
					item.Id: map[string]*Inventory{
						location1.Id: &inventory1,
					},
				},
			},
			id:   item.Id,
			want: []*Inventory{&inventory1},
		},
		{
			desc: "multiple locations",
			mb: &InMemoryBackend{
				inventoryByItemByLocationIndex: map[string]map[string]*Inventory{
					item.Id: map[string]*Inventory{
						location1.Id: &inventory1,
						location2.Id: &inventory2,
					},
				},
			},
			id:   item.Id,
			want: []*Inventory{&inventory1, &inventory2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.mb.ListItemInventory(ctx, tc.id)

			if err != nil {
				t.Fatalf("tc.mb.ListItemInventory(%v) returned unexpected err: %v", tc.id, err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("tc.mb.ListItemInventory(%v) = %v want %v", tc.id, got, tc.want)
			}
		})
	}
}

func TestListItemInventoryTransactions(t *testing.T) {
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
	mb := InMemoryBackend{
		inventoryTransactions: map[string]*InventoryTransaction{
			item1Txn1.Id: &item1Txn1,
			item2Txn1.Id: &item2Txn1,
			item2Txn2.Id: &item2Txn2,
		},
	}
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
			got, err := mb.ListItemInventoryTransactions(ctx, tc.id)

			if err != nil {
				t.Fatalf("mb.ListItemInventoryTransactions() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("mb.ListItemInventoryTransactions() = %v want %v", got, tc.want)
			}
		})
	}
}

func TestListInventoryTransactions(t *testing.T) {
	txn1, txn2 := InventoryTransaction{Id: "txn1-id"}, InventoryTransaction{Id: "txn2-id"}
	cases := []struct {
		desc string
		mb   *InMemoryBackend
		want []*InventoryTransaction
	}{
		{
			desc: "no txns",
			mb:   NewInMemoryBackend(),
			want: []*InventoryTransaction{},
		},
		{
			desc: "single txn",
			mb:   &InMemoryBackend{inventoryTransactions: map[string]*InventoryTransaction{txn1.Id: &txn1}},
			want: []*InventoryTransaction{&txn1},
		},
		{
			desc: "multiple txns",
			mb:   &InMemoryBackend{inventoryTransactions: map[string]*InventoryTransaction{txn1.Id: &txn1, txn2.Id: &txn2}},
			want: []*InventoryTransaction{&txn1, &txn2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.mb.ListInventoryTransactions(ctx)

			if err != nil {
				t.Fatalf("tc.mb.ListInventoryTransactions() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("tc.mb.ListInventoryTransactions() = %v want %v", got, tc.want)
			}
		})
	}
}

func TestListLocations(t *testing.T) {
	location1, location2 := Location{Id: "location1-id"}, Location{Id: "location2-id"}
	cases := []struct {
		desc string
		mb   *InMemoryBackend
		want []*Location
	}{
		{
			desc: "no locations",
			mb:   NewInMemoryBackend(),
			want: []*Location{},
		},
		{
			desc: "single location",
			mb:   &InMemoryBackend{locations: map[string]*Location{location1.Id: &location1}},
			want: []*Location{&location1},
		},
		{
			desc: "multiple locations",
			mb:   &InMemoryBackend{locations: map[string]*Location{location1.Id: &location1, location2.Id: &location2}},
			want: []*Location{&location1, &location2},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.mb.ListLocations(ctx)

			if err != nil {
				t.Fatalf("tc.mb.ListLocations() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("tc.mb.ListLocations() = %v want %v", got, tc.want)
			}
		})
	}
}

func TestListLocationInventory(t *testing.T) {
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
	mb := InMemoryBackend{
		inventoryByLocationByItemIndex: map[string]map[string]*Inventory{
			loc1.Id: map[string]*Inventory{
				loc1Inv1.ItemId: &loc1Inv1,
			},
			loc2.Id: map[string]*Inventory{
				loc2Inv1.ItemId: &loc2Inv1,
				loc2Inv2.ItemId: &loc2Inv2,
			},
		},
	}
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
			got, err := mb.ListLocationInventory(ctx, tc.id)

			if err != nil {
				t.Fatalf("mb.ListLocationInventory(%v) returned unexpected err: %v", tc.id, err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("mb.ListLocationInventory(%v) = %v want %v", tc.id, got, tc.want)
			}
		})
	}
}

func TestListLocationInventoryTransactions(t *testing.T) {
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
	mb := InMemoryBackend{
		inventoryTransactions: map[string]*InventoryTransaction{
			loc1Txn1.Id: &loc1Txn1,
			loc2Txn1.Id: &loc2Txn1,
			loc2Txn2.Id: &loc2Txn2,
		},
	}
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
			got, err := mb.ListLocationInventoryTransactions(ctx, tc.id)

			if err != nil {
				t.Fatalf("mb.ListLocationInventoryTransactions() returned unexpected err: %v", err)
			}
			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(modelLess), cmpopts.EquateEmpty()) {
				t.Errorf("mb.ListLocationInventoryTransactions() = %v want %v", got, tc.want)
			}
		})
	}
}

func TestNewItem(t *testing.T) {
	ctx := context.Background()
	mb := NewInMemoryBackend()
	item := Item{
		Id:          "id-to-be-replaced-by-uuid",
		Name:        "name",
		Description: "description",
	}

	got, err := mb.NewItem(ctx, &item)
	if err != nil {
		t.Fatalf("mb.NewItem(%v) returned unexpected err: %v", item, err)
	}
	if got.Id == "" {
		t.Errorf("mb.NewItem(%v) did not generate Item.Id", item)
	}
	if !cmp.Equal(got, &item, cmpopts.IgnoreFields(Item{}, "Id")) {
		t.Errorf("mb.NewItem(%v) = %v want %v (ignoring Id field)", item, got, item)
	}
	if v := mb.items[got.Id]; v != got {
		t.Errorf("after mb.NewItem(%v), mb.items[%v] = %v want %v", item, got.Id, v, got)
	}
}

func TestNewInventoryTransaction(t *testing.T) {
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
		mb := InMemoryBackend{
			inventoryByItemByLocationIndex: map[string]map[string]*Inventory{
				stockedItem.Id: map[string]*Inventory{stockedLoc.Id: &stockedInv},
			},
			inventoryByLocationByItemIndex: map[string]map[string]*Inventory{
				stockedLoc.Id: map[string]*Inventory{stockedItem.Id: &stockedInv},
			},
			inventoryTransactions: make(map[string]*InventoryTransaction),
			items:                 map[string]*Item{newItem.Id: &newItem, stockedItem.Id: &stockedItem},
			locations:             map[string]*Location{newLoc.Id: &newLoc, stockedLoc.Id: &stockedLoc},
		}
		t.Run(tc.desc, func(t *testing.T) {

			got, err := mb.NewInventoryTransaction(ctx, tc.txn)
			if err != nil {
				t.Fatalf("mb.NewInventoryTransaction(%v) returned unexpected err: %v", tc.txn, err)
			}
			if got.Id == "" {
				t.Errorf("mb.NewInventoryTransaction(%v) did not generate InventoryTransaction.Id", tc.txn)
			}
			if !cmp.Equal(got, tc.txn, cmpopts.IgnoreFields(InventoryTransaction{}, "Id", "Timestamp")) {
				t.Errorf("mb.NewInventoryTransaction(%v) = %v want %v (ignoring Id field)", tc.txn, got, tc.txn)
			}
			if v := mb.inventoryTransactions[got.Id]; v != got {
				t.Errorf("after mb.NewInventoryTransaction(%v), mb.inventoryTransactions[%v] = %v want %v", tc.txn, got.Id, v, got)
			}

			wantInv := &Inventory{ItemId: tc.txn.ItemId, LocationId: tc.txn.LocationId, Count: tc.wantCount}

			inv := mb.inventoryByItemByLocationIndex[tc.txn.ItemId][tc.txn.LocationId]
			invByLoc := mb.inventoryByLocationByItemIndex[tc.txn.LocationId][tc.txn.ItemId]
			if inv != invByLoc {
				t.Errorf("mb.NewInventoryTransaction(%v) produced inconsistent index inventory values", tc.txn)
				t.Errorf("byItem index value: %v, byLoc index value: %v", inv, invByLoc)
			}
			if inv == nil {
				t.Errorf("mb.NewInventoryTransaction(%v) produced nil inventory", tc.txn)
			}
			if !cmp.Equal(inv, wantInv, cmpopts.IgnoreFields(Inventory{}, "LastUpdated")) {
				t.Errorf("mb.NewInventoryTransaction(%v) produced unexpected inventory state", tc.txn)
				t.Errorf("[item %q, location %q] inventory = %v want %v", tc.txn.ItemId, tc.txn.LocationId, inv, wantInv)
			}
		})
	}

}

func TestNewInventoryTransactionNotFoundErrors(t *testing.T) {
	existingLoc := Location{Id: "existing-location-id"}
	existingItem := Item{Id: "existing-item-id"}
	mb := InMemoryBackend{
		items: map[string]*Item{
			existingItem.Id: &existingItem,
		},
		locations: map[string]*Location{
			existingLoc.Id: &existingLoc,
		},
	}
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
			_, err := mb.NewInventoryTransaction(ctx, tc.txn)
			if err == nil {
				t.Fatalf("mb.NewInventoryTransaction(%q) succeeded, want error", tc.txn)
			}
			if nf, ok := err.(*ResourceNotFound); !ok || nf.id != tc.want.id || nf.collection != tc.want.collection {
				t.Errorf("mb.NewInventoryTransaction(%q) returned %v, want %v", tc.txn, err, tc.want)
			}
		})
	}
}

func TestNewLocation(t *testing.T) {
	ctx := context.Background()
	mb := NewInMemoryBackend()
	location := Location{
		Name:      "name",
		Warehouse: "warehouse",
	}

	got, err := mb.NewLocation(ctx, &location)

	if err != nil {
		t.Fatalf("mb.NewLocation(%v) returned unexpected err: %v", location, err)
	}
	if got.Id == "" {
		t.Errorf("mb.NewLocation(%v) did not generate Location.Id", location)
	}
	if !cmp.Equal(got, &location, cmpopts.IgnoreFields(Location{}, "Id")) {
		t.Errorf("mb.NewLocation(%v) = %v want %v (ignoring Id field)", location, got, location)
	}
	if v := mb.locations[got.Id]; v != got {
		t.Errorf("after mb.NewLocation(%v), mb.locations[%v] = %v want %v", location, got.Id, v, got)
	}
}

func TestUpdateItem(t *testing.T) {
	ctx := context.Background()
	id := "item-id"
	mb := InMemoryBackend{
		items: map[string]*Item{
			id: &Item{
				Id:          id,
				Name:        "old-name",
				Description: "old-description",
			},
		},
	}
	item := &Item{
		Id:          id,
		Name:        "updated-name",
		Description: "updated-description",
	}

	got, err := mb.UpdateItem(ctx, item)

	if err != nil {
		t.Fatalf("mb.UpdateItem(%v) returned unexpected err: %v", item, err)
	}
	if !cmp.Equal(got, item) {
		t.Errorf("mb.UpdateItem(%v) = %v want %v", item, got, item)
	}
	if mb.items[id] != got {
		t.Errorf("after mb.UpdateItem(%v), mb.items[%v] = %v want %v", item, id, mb.items[id], got)
	}
}

func TestUpdateItemNotFound(t *testing.T) {
	ctx := context.Background()
	item := &Item{Id: "not-found-id"}
	mb := NewInMemoryBackend()
	want := ItemNotFound(item.Id)

	_, err := mb.UpdateItem(ctx, item)

	if err == nil {
		t.Fatalf("mb.UpdateItem(%q) succeeded, want error", item)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.UpdateItem(%q) returned %v, want %v", item, err, want)
	}
}

func TestUpdateLocation(t *testing.T) {
	ctx := context.Background()
	id := "location-id"
	mb := InMemoryBackend{
		locations: map[string]*Location{
			id: &Location{
				Id:        id,
				Name:      "old-name",
				Warehouse: "old-warehouse",
			},
		},
	}
	location := &Location{
		Id:        id,
		Name:      "updated-name",
		Warehouse: "updated-warehouse",
	}

	got, err := mb.UpdateLocation(ctx, location)

	if err != nil {
		t.Fatalf("mb.UpdateLocation(%v) returned unexpected err: %v", location, err)
	}
	if !cmp.Equal(got, location) {
		t.Errorf("mb.UpdateLocation(%v) = %v want %v", location, got, location)
	}
	if mb.locations[id] != got {
		t.Errorf("after mb.UpdateLocation(%v), mb.locations[%v] = %v want %v", location, id, mb.locations[id], got)
	}
}

func TestUpdateLocationNotFound(t *testing.T) {
	location := &Location{Id: "not-found-id"}
	mb := NewInMemoryBackend()
	ctx := context.Background()
	want := LocationNotFound(location.Id)

	_, err := mb.UpdateLocation(ctx, location)

	if err == nil {
		t.Fatalf("mb.UpdateLocation(%q) succeeded, want error", location)
	}
	if nf, ok := err.(*ResourceNotFound); !ok || nf.id != want.id || nf.collection != want.collection {
		t.Errorf("mb.UpdateLocation(%q) returned %v, want %v", location, err, want)
	}
}
