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

package main

import (
	"context"
	"flag"
	client "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/api-client"
	state "github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/inventory-state-service/src"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
)

var calculator *state.InventoryStateCalculator

// Transforms an InventoryTransaction event to an ItemInventoryState event
func publishState(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	log.Printf("received event: %s\n", event)
	inventoryTransaction := &client.InventoryTransaction{}
	if err := event.DataAs(inventoryTransaction); err != nil {
		log.Printf("could not get InventoryTransaction from event data: %s\n", err.Error())
		return nil, err
	}

	s, err := calculator.GetItemInventoryState(inventoryTransaction)
	if err != nil {
		log.Printf("could not calculate inventory state: %s\n", err.Error())
		return nil, err
	}

	e := itemInventoryStateToEvent(inventoryTransaction.Id, *s)
	return &e, nil
}

func itemInventoryStateToEvent(transactionId string, s state.ItemInventoryState) cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetID(transactionId)
	e.SetSource("inventory-state-service")
	e.SetType("state.ItemInventoryState")
	e.SetData(cloudevents.ApplicationJSON, s)
	e.SetExtension("iteminventoryclassification", s.Classification)
	return e
}

func main() {
	log.Printf("server started")
	backendClusterHostName := flag.String("backend_cluster_host_name", "", "cluster hostname of the api service")
	flag.Parse()
	if *backendClusterHostName == "" {
		log.Fatal("backend_cluster_host_name must be set")
	}
	calculator = state.NewInventoryStateCalculator(*backendClusterHostName)

	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create cloudevents client, %s", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), publishState))
}
