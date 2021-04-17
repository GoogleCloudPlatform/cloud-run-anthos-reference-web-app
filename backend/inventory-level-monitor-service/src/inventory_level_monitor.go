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
	"errors"
	"fmt"
	client "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/api-client"
	state "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/inventory-state-service/src"
	optional "github.com/antihax/optional"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
	"time"
)

// Creates low and high total inventory alerts for items.
type InventoryLevelMonitor struct {
	apiClient *client.APIClient
}

// Creates a new InventoryLevelMonitor
func NewInventoryLevelMonitor(apiHostname string) *InventoryLevelMonitor {
	cfg := client.NewConfiguration()
	cfg.BasePath = fmt.Sprintf("http://%s/api", apiHostname)
	apiClient := client.NewAPIClient(cfg)
	return &InventoryLevelMonitor{apiClient}
}

// Creates an alert out of the ItemInventoryState in the event.
func (m InventoryLevelMonitor) Monitor(ctx context.Context, event cloudevents.Event) error {
	log.Printf("received event: %s\n", event)
	state := &state.ItemInventoryState{}
	if err := event.DataAs(state); err != nil {
		log.Printf("could not get ItemInventoryState from event data: %s\n", err.Error())
		return err
	}

	if state.ItemId == "" {
		err := errors.New("cannot create alert from ItemInventoryState without an ItemId")
		log.Printf("%s\n", err.Error())
		return err
	}
	alertOpts := client.NewAlertOpts{
		Alert: optional.NewInterface(client.Alert{
			ItemId:    state.ItemId,
			Text:      fmt.Sprintf("%s total inventory for item: %d.", state.Classification, state.TotalCount),
			Timestamp: time.Now(),
		}),
	}
	_, _, err := m.apiClient.AlertApi.NewAlert(ctx, &alertOpts)
	return err
}
