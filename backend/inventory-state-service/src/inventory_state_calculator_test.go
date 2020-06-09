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

package state

import (
	"encoding/json"
	"fmt"
	client "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/api-client"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const itemId = "123"
const numInLocationA = 10
const numInLocationB = 15
const total = numInLocationA + numInLocationB

var calculator InventoryStateCalculator

func TestMain(m *testing.M) {
	os.Exit(setUp(m))
}

func setUp(m *testing.M) int {
	ts := httptest.NewServer(newFakeInventoryApi())
	defer ts.Close()
	cfg := client.NewConfiguration()
	cfg.Host = ts.Listener.Addr().String()
	apiClient := client.NewAPIClient(cfg)
	calculator = InventoryStateCalculator{
		apiClient: apiClient,
	}
	return m.Run()
}

func TestGetItemInventoryStateSuccess(t *testing.T) {
	inventoryTransaction := client.InventoryTransaction{
		ItemId: itemId,
		Action: "ADD",
	}
	state, err := calculator.GetItemInventoryState(&inventoryTransaction)
	if err != nil {
		t.Fatalf("Expected success; got %s", err)
	}

	expectedState := ItemInventoryState{
		ItemId:     itemId,
		TotalCount: total,
	}
	if !cmp.Equal(state, &expectedState) {
		t.Errorf("Expected = %v; got %v", expectedState, state)
	}
}

func TestGetItemInventoryStateError(t *testing.T) {
	inventoryTransaction := client.InventoryTransaction{
		ItemId: "BadItemId",
		Action: "ADD",
	}
	_, err := calculator.GetItemInventoryState(&inventoryTransaction)
	if err == nil {
		t.Fatalf("Expected an error due to an incorrect item id.")
	}
}

func newFakeInventoryApi() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/items/%s/inventory", itemId), func(w http.ResponseWriter, r *http.Request) {

		var result = []client.Inventory{
			client.Inventory{
				ItemId:     itemId,
				LocationId: "a",
				Count:      numInLocationA,
			},
			client.Inventory{
				ItemId:     itemId,
				LocationId: "b",
				Count:      numInLocationB,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
		}
	})
	return mux
}
