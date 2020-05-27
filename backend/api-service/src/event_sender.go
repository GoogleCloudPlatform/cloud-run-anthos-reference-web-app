package service

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type EventSender interface {
	SendInventoryTransactionEvent(inventoryTransaction InventoryTransaction) error
}

// BrokerEventSender sends cloudevents to a Broker for API actions
type BrokerEventSender struct {
	client cloudevents.Client
}

func NewBrokerEventSender(eventBrokerHostname string) (*BrokerEventSender, error) {
	p, err := cloudevents.NewHTTP(cloudevents.WithTarget(eventBrokerHostname))
	if err != nil {
		return nil, err
	}
	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow())
	if err != nil {
		return nil, err
	}
	return &BrokerEventSender{c}, nil
}

// Send an event for the given InventoryTransaction
func (s *BrokerEventSender) SendInventoryTransactionEvent(inventoryTransaction InventoryTransaction) error {
	event := inventoryTransactionToEvent(inventoryTransaction)
	if result := s.client.Send(context.Background(), event); !cloudevents.IsACK(result) {
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
