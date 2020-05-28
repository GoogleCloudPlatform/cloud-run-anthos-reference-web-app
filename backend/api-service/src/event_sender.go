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
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// EventSender sends cloudevents
type EventSender interface {
	SendInventoryTransactionEvent(inventoryTransaction InventoryTransaction) error
}

type brokerEventSender struct {
	client cloudevents.Client
	hostname string
}

// Creates a new brokerEventSender
func NewBrokerEventSender(eventBrokerHostname string) (EventSender, error) {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &brokerEventSender{c, eventBrokerHostname}, nil
}

// Send an event for the given InventoryTransaction
func (s brokerEventSender) SendInventoryTransactionEvent(inventoryTransaction InventoryTransaction) error {
	event := inventoryTransactionToEvent(inventoryTransaction)
	ctx := cloudevents.ContextWithTarget(context.Background(), s.hostname)
	if result := s.client.Send(ctx, event); !cloudevents.IsACK(result) {
		return result
	}
	return nil
}

func inventoryTransactionToEvent(inventoryTransaction InventoryTransaction) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID(inventoryTransaction.Id)
	event.SetSource("api-service")
	event.SetType("service.InventoryTransaction")
	event.SetData(cloudevents.ApplicationJSON, inventoryTransaction)
	return event
}
