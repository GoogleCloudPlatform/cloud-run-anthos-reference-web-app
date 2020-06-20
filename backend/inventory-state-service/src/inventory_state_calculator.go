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
	"context"
	"fmt"
	client "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/api-client"
)

type InventoryStateCalculator struct {
	apiClient *client.APIClient
}

// Creates a new InventoryStateCalculator
func NewInventoryStateCalculator(apiHostname string) *InventoryStateCalculator {
	cfg := client.NewConfiguration()
	cfg.BasePath = fmt.Sprintf("http://%s/api", apiHostname)
	apiClient := client.NewAPIClient(cfg)
	return &InventoryStateCalculator{apiClient}
}

func (r InventoryStateCalculator) GetItemInventoryState(i *client.InventoryTransaction) (*ItemInventoryState, error) {
	inventoryList, _, err := r.apiClient.InventoryApi.ListItemInventory(context.Background(), i.ItemId)
	if err != nil {
		return nil, err
	}
	var total int64 = 0
	for _, itemInv := range inventoryList {
		total += itemInv.Count
	}

	s := ItemInventoryState{
		ItemId:         i.ItemId,
		TotalCount:     total,
		Classification: getItemInventoryClassification(total),
	}
	return &s, nil
}

func getItemInventoryClassification(total int64) ItemInventoryClassification {
	switch {
	case total < 100:
		return Low
	case total >= 1000:
		return High
	default:
		return Normal
	}
}
