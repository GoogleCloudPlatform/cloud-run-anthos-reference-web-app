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
	state "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/inventory-state-service/src"
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
var inventoryLevelMonitor InventoryLevelMonitor

func TestMain(m *testing.M) {
	os.Exit(setUp(m))
}

func setUp(m *testing.M) int {
	ts := httptest.NewServer(newFakeAlertApi())
	defer ts.Close()
	cfg := client.NewConfiguration()
	cfg.Host = ts.Listener.Addr().String()
	apiClient = client.NewAPIClient(cfg)
	inventoryLevelMonitor = InventoryLevelMonitor{
		apiClient: apiClient,
	}
	return m.Run()
}

func TestMonitorSuccess(t *testing.T) {
	itemId := "123"
	state := state.ItemInventoryState{
		ItemId:         itemId,
		TotalCount:     9,
		Classification: "Low",
	}
	err := inventoryLevelMonitor.Monitor(context.Background(), createEventForTest(state))
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
		ItemId:    itemId,
		Text:      "Low total inventory for item: 9.",
		Timestamp: time.Now(),
	}
	if !cmp.Equal(alertList[0], expectedAlert, cmpopts.EquateApproxTime(time.Minute)) {
		t.Errorf("Expected = %v; got %v", alertList[0], expectedAlert)
	}
}

func TestMonitorErrorBadDataInEvent(t *testing.T) {
	event := createEventForTest(state.ItemInventoryState{})
	event.SetData("badformat", map[string]string{"bad": "data"})
	err := inventoryLevelMonitor.Monitor(context.Background(), event)
	if err == nil {
		t.Fatalf("Expected an error due to bad data in the event")
	}
}

func TestMonitorErrorMissingId(t *testing.T) {
	err := inventoryLevelMonitor.Monitor(context.Background(), createEventForTest(state.ItemInventoryState{}))
	if err == nil {
		t.Fatalf("Expected an error due to a missing item id in ItemInventoryState")
	}
}

func createEventForTest(s state.ItemInventoryState) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID("event id for test")
	event.SetSource("event source for test")
	event.SetType("event type for test")
	event.SetData(cloudevents.ApplicationJSON, s)
	return event
}

func newFakeAlertApi() http.Handler {
	mux := http.NewServeMux()
	var alerts []client.Alert
	mux.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			alert := &client.Alert{}
			if err := json.NewDecoder(r.Body).Decode(alert); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			alerts = append(alerts, *alert)
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(alerts); err != nil {
				http.Error(w, fmt.Sprintf("can't encode inventory: %v", err), http.StatusInternalServerError)
			}
		}
	})
	return mux
}
