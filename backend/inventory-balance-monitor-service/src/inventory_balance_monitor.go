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
	"fmt"
	client "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/api-client"
	optional "github.com/antihax/optional"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

var inventoryRatioThreshold int64 = 5
var alertText = "Inventory for item in Warehouse \"%s\" (%d) is %d times more than in " +
	"Warehouse \"%s\" (%d). Consider rebalancing the inventory."

// Creates inventory balance alerts for items.
type InventoryBalanceMonitor struct {
	apiClient *client.APIClient
}

// Creates a new InventoryBalanceMonitor
func NewInventoryBalanceMonitor(apiHostname string) *InventoryBalanceMonitor {
	cfg := client.NewConfiguration()
	cfg.BasePath = fmt.Sprintf("http://%s/api", apiHostname)
	apiClient := client.NewAPIClient(cfg)
	return &InventoryBalanceMonitor{apiClient}
}

// Checks the inventory balance across warehouses for the item in the transaction in the event,
// and creates an alert if necessary.
// Alerts are created when the warehouse with the most inventory has 5x or more inventory
// than the warehouse with the least inventory. Warehouses that have an inventory of 0
// for an item are not looked at.
func (m InventoryBalanceMonitor) Monitor(ctx context.Context, event cloudevents.Event) error {
	log.Printf("received event: %s\n", event)
	inventoryTransaction := &client.InventoryTransaction{}
	if err := event.DataAs(inventoryTransaction); err != nil {
		log.Printf("could not get InventoryTransaction from event data: %s\n", err.Error())
		return err
	}

	warehouseToItemCount, err := m.getWarehouseToInventoryMap(inventoryTransaction.ItemId)
	if err != nil {
		return err
	}

	// The item may not be in any warehouses if its inventory is down to 0 everywhere.
	// We don't need to check inventory balance in this case.
	if len(warehouseToItemCount) == 0 {
		return nil
	}
	most := m.getWarehouseWithMostInventory(warehouseToItemCount)
	least := m.getWarehouseWithLeastInventory(warehouseToItemCount)
	ratio := warehouseToItemCount[most] / warehouseToItemCount[least]
	if ratio >= inventoryRatioThreshold {
		alertOpts := client.NewAlertOpts{
			Alert: optional.NewInterface(client.Alert{
				ItemId:    inventoryTransaction.ItemId,
				Text:      fmt.Sprintf(alertText, most, warehouseToItemCount[most], ratio, least, warehouseToItemCount[least]),
				Timestamp: time.Now(),
			}),
		}
		_, _, err := m.apiClient.AlertApi.NewAlert(ctx, &alertOpts)
		return err
	}
	return nil
}

func (m InventoryBalanceMonitor) getWarehouseWithMostInventory(warehouseToItemCount map[string]int64) string {
	var (
		warehouse     string
		mostInventory int64
	)
	for w, count := range warehouseToItemCount {
		if warehouse == "" || count > mostInventory {
			warehouse = w
			mostInventory = count
		}
	}
	return warehouse
}

func (m InventoryBalanceMonitor) getWarehouseWithLeastInventory(warehouseToItemCount map[string]int64) string {
	var (
		warehouse      string
		leastInventory int64
	)
	for w, count := range warehouseToItemCount {
		if warehouse == "" || count < leastInventory {
			warehouse = w
			leastInventory = count
		}
	}
	return warehouse
}

func (m InventoryBalanceMonitor) getWarehouseToInventoryMap(itemId string) (map[string]int64, error) {
	var (
		g             errgroup.Group
		inventoryList []client.Inventory
		locationList  []client.Location
	)

	g.Go(func() error {
		var err error
		inventoryList, _, err = m.apiClient.InventoryApi.ListItemInventory(context.Background(), itemId)
		return err
	})
	g.Go(func() error {
		var err error
		locationList, _, err = m.apiClient.InventoryApi.ListLocations(context.Background())
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	locationIdToWarehouse := make(map[string]string)
	for _, location := range locationList {
		locationIdToWarehouse[location.Id] = location.Warehouse
	}

	warehouseToItemCount := make(map[string]int64)
	for _, itemInv := range inventoryList {
		if itemInv.Count != 0 {
			warehouse := locationIdToWarehouse[itemInv.LocationId]
			warehouseToItemCount[warehouse] = warehouseToItemCount[warehouse] + itemInv.Count
		}
	}
	return warehouseToItemCount, nil
}
