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

package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	client "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/api-client"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var apiClient *client.APIClient
var inventoryBalanceMonitor InventoryBalanceMonitor
var generatedAlerts []client.Alert

var unbalancedItemId = "111"
var balancedItemId = "222"
var zeroInvItemId = "333"

func TestMain(m *testing.M) {
	os.Exit(setUp(m))
}

func setUp(m *testing.M) int {
	ts := httptest.NewServer(newFakeApi())
	defer ts.Close()
	cfg := client.NewConfiguration()
	cfg.Host = ts.Listener.Addr().String()
	apiClient = client.NewAPIClient(cfg)
	inventoryBalanceMonitor = InventoryBalanceMonitor{
		apiClient: apiClient,
	}
	return m.Run()
}

func TestMonitorSuccessUnbalancedInventory(t *testing.T) {
	clearAlerts()
	transaction := client.InventoryTransaction{
		ItemId: unbalancedItemId,
	}
	err := inventoryBalanceMonitor.Monitor(context.Background(), createEventForTest(transaction))
	if err != nil {
		t.Fatalf("Expected success; got %s", err)
	}
	// Read alerts back, check that there's a single one and that it's correct.
	alertList, _, err := apiClient.AlertApi.ListAlerts(context.Background())
	if err != nil {
		t.Fatalf("Expected success; got %s", err)
	}
	if len(alertList) != 1 {
		t.Fatalf("Expected a single alert; got %s", alertList)
	}

	expectedAlert := client.Alert{
		ItemId: unbalancedItemId,
		Text: "Inventory for item in Warehouse \"B\" (600) is 40 times more than in " +
			"Warehouse \"A\" (15). Consider rebalancing the inventory.",
		Timestamp: time.Now(),
	}
	if !cmp.Equal(alertList[0], expectedAlert, cmpopts.EquateApproxTime(time.Minute)) {
		t.Errorf("Expected = %v; got %v", alertList[0], expectedAlert)
	}
}

func TestNoAlertForItemWithZeroInventory(t *testing.T) {
	checkSucessAndNoAlert(zeroInvItemId, t)
}

func TestNoAlertForItemWithBalancedInventory(t *testing.T) {
	checkSucessAndNoAlert(balancedItemId, t)
}

func checkSucessAndNoAlert(itemId string, t *testing.T) {
	clearAlerts()
	transaction := client.InventoryTransaction{ItemId: itemId}
	err := inventoryBalanceMonitor.Monitor(context.Background(), createEventForTest(transaction))
	if err != nil {
		t.Fatalf("Expected success; got %s", err)
	}
	alertList, _, err := apiClient.AlertApi.ListAlerts(context.Background())
	if err != nil {
		t.Fatalf("Expected success; got %s", err)
	}
	if len(alertList) != 0 {
		t.Fatalf("Expected no alerts; got %s", alertList)
	}
}

func TestMonitorErrorBadDataInEvent(t *testing.T) {
	event := createEventForTest(client.InventoryTransaction{})
	event.SetData("badformat", map[string]string{"bad": "data"})
	err := inventoryBalanceMonitor.Monitor(context.Background(), event)
	if err == nil {
		t.Fatalf("Expected an error due to bad data in the event")
	}
}

func createEventForTest(t client.InventoryTransaction) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID("event id for test")
	event.SetSource("event source for test")
	event.SetType("event type for test")
	event.SetData(cloudevents.ApplicationJSON, t)
	return event
}

func clearAlerts() {
	generatedAlerts = nil
}

func newFakeApi() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			alert := &client.Alert{}
			if err := json.NewDecoder(r.Body).Decode(alert); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			generatedAlerts = append(generatedAlerts, *alert)
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(generatedAlerts); err != nil {
				http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
			}
		}
	})

	// unbalancedItemId has inventory in 5 locations across 3 warehouses, one of them is 0.
	mux.HandleFunc(fmt.Sprintf("/api/items/%s/inventory", unbalancedItemId), func(w http.ResponseWriter, r *http.Request) {
		var result = []client.Inventory{
			client.Inventory{ItemId: unbalancedItemId, LocationId: "a1", Count: 10},
			client.Inventory{ItemId: unbalancedItemId, LocationId: "a2", Count: 5},
			client.Inventory{ItemId: unbalancedItemId, LocationId: "b1", Count: 100},
			client.Inventory{ItemId: unbalancedItemId, LocationId: "b2", Count: 500},
			client.Inventory{ItemId: unbalancedItemId, LocationId: "c1", Count: 0},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
		}
	})

	// zeroInvItemId has an inventory of zero in 2 different warehouses
	mux.HandleFunc(fmt.Sprintf("/api/items/%s/inventory", zeroInvItemId), func(w http.ResponseWriter, r *http.Request) {
		var result = []client.Inventory{
			client.Inventory{ItemId: zeroInvItemId, LocationId: "b2", Count: 0},
			client.Inventory{ItemId: zeroInvItemId, LocationId: "c1", Count: 0},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
		}
	})

	// balancedItemId has inventory in 3 different warehouses
	mux.HandleFunc(fmt.Sprintf("/api/items/%s/inventory", balancedItemId), func(w http.ResponseWriter, r *http.Request) {
		var result = []client.Inventory{
			client.Inventory{ItemId: balancedItemId, LocationId: "a1", Count: 100},
			client.Inventory{ItemId: balancedItemId, LocationId: "b1", Count: 200},
			client.Inventory{ItemId: balancedItemId, LocationId: "c1", Count: 300},
			client.Inventory{ItemId: balancedItemId, LocationId: "d1", Count: 0},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/locations", func(w http.ResponseWriter, r *http.Request) {
		var result = []client.Location{
			client.Location{Id: "a1", Warehouse: "A"},
			client.Location{Id: "a2", Warehouse: "A"},
			client.Location{Id: "b1", Warehouse: "B"},
			client.Location{Id: "b2", Warehouse: "B"},
			client.Location{Id: "c1", Warehouse: "C"},
			client.Location{Id: "d1", Warehouse: "D"},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
		}
	})

	return mux
}
