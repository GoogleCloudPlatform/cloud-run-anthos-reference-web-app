package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewInventoryTransactionBadRequests(t *testing.T) {
	cases := []struct {
		desc string
		txn  InventoryTransaction
		msg  string
	}{
		{
			desc: "missing action field",
			txn:  InventoryTransaction{ItemId: "iid", LocationId: "lid"},
			msg:  "required field: action",
		},
		{
			desc: "missing item_id field",
			txn:  InventoryTransaction{Action: "ADD", LocationId: "lid"},
			msg:  "required field: item_id",
		},
		{
			desc: "missing location_id field",
			txn:  InventoryTransaction{Action: "ADD", ItemId: "iid"},
			msg:  "required field: location_id",
		},
		{
			desc: "bad action type",
			txn:  InventoryTransaction{Action: "bad-action", ItemId: "iid", LocationId: "lid"},
			msg:  "Unknown action",
		},
	}

	for _, tc := range cases {
		s := InventoryApiService{}
		t.Run(tc.desc, func(t *testing.T) {
			r := httptest.NewRecorder()
			err := s.NewInventoryTransaction(tc.txn, r)

			if err != nil {
				t.Errorf("s.NewInventoryTransaction(%v) returned unexpected error: %v", tc.txn, err)
			}
			if r.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("status code: %v, want: %v", r.Result().StatusCode, http.StatusBadRequest)
			}
			if !strings.Contains(r.Body.String(), tc.msg) {
				t.Errorf("response body %q does not contain %q", r.Body.String(), tc.msg)
			}
		})
	}
}

func TestNewItemMissingNameField(t *testing.T) {
	item := Item{}
	msg := "required field: name"

	s := InventoryApiService{}
	r := httptest.NewRecorder()
	err := s.NewItem(item, r)

	if err != nil {
		t.Errorf("s.NewItem(%v) returned unexpected error: %v", item, err)
	}
	if r.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("status code: %v, want: %v", r.Result().StatusCode, http.StatusBadRequest)
	}
	if !strings.Contains(r.Body.String(), msg) {
		t.Errorf("response body %q does not contain %q", r.Body.String(), msg)
	}
}

func TestNewLocationBadRequests(t *testing.T) {
	cases := []struct {
		desc string
		loc  Location
		msg  string
	}{
		{
			desc: "missing warehouse field",
			loc:  Location{Name: "name"},
			msg:  "required field: warehouse",
		},
		{
			desc: "missing name field",
			loc:  Location{Warehouse: "warehouse"},
			msg:  "required field: name",
		},
	}

	for _, tc := range cases {
		s := InventoryApiService{}
		t.Run(tc.desc, func(t *testing.T) {
			r := httptest.NewRecorder()
			err := s.NewLocation(tc.loc, r)

			if err != nil {
				t.Errorf("s.NewLocation(%v) returned unexpected error: %v", tc.loc, err)
			}
			if r.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("status code: %v, want: %v", r.Result().StatusCode, http.StatusBadRequest)
			}
			if !strings.Contains(r.Body.String(), tc.msg) {
				t.Errorf("response body %q does not contain %q", r.Body.String(), tc.msg)
			}
		})
	}
}

func TestUpdateItemBadRequests(t *testing.T) {
	id := "item-id"
	cases := []struct {
		desc string
		item Item
		msg  string
	}{
		{
			desc: "mismatched item id",
			item: Item{Id: "bad" + id, Name: "name"},
			msg:  "Mismatched path id",
		},
		{
			desc: "missing name field",
			item: Item{Id: "item-id"},
			msg:  "required field: name",
		},
	}

	for _, tc := range cases {
		s := InventoryApiService{}
		t.Run(tc.desc, func(t *testing.T) {
			r := httptest.NewRecorder()
			err := s.UpdateItem(id, tc.item, r)

			if err != nil {
				t.Errorf("s.UpdateItem(%v, %v) returned unexpected error: %v", id, tc.item, err)
			}
			if r.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("status code: %v, want: %v", r.Result().StatusCode, http.StatusBadRequest)
			}
			if !strings.Contains(r.Body.String(), tc.msg) {
				t.Errorf("response body %q does not contain %q", r.Body.String(), tc.msg)
			}
		})
	}
}

func TestUpdateLocationBadRequests(t *testing.T) {
	id := "location-id"
	cases := []struct {
		desc string
		loc  Location
		msg  string
	}{
		{
			desc: "mismatched item id",
			loc:  Location{Id: "bad" + id, Name: "name", Warehouse: "warehouse"},
			msg:  "Mismatched path id",
		},
		{
			desc: "missing warehouse field",
			loc:  Location{Id: id, Name: "name"},
			msg:  "required field: warehouse",
		},
		{
			desc: "missing name field",
			loc:  Location{Id: id, Warehouse: "warehouse"},
			msg:  "required field: name",
		},
	}

	for _, tc := range cases {
		s := InventoryApiService{}
		t.Run(tc.desc, func(t *testing.T) {
			r := httptest.NewRecorder()
			err := s.UpdateLocation(id, tc.loc, r)

			if err != nil {
				t.Errorf("s.UpdateLocation(%v, %v) returned unexpected error: %v", id, tc.loc, err)
			}
			if r.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("status code: %v, want: %v", r.Result().StatusCode, http.StatusBadRequest)
			}
			if !strings.Contains(r.Body.String(), tc.msg) {
				t.Errorf("response body %q does not contain %q", r.Body.String(), tc.msg)
			}
		})
	}
}
