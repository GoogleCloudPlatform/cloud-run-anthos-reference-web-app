package service

import (
	"context"
	"flag"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
)

// EventSender sends cloud events containing an InventoryTransaction
type EventSender struct {
	eventBrokerHostname string
}

// NewInventoryApiService creates a default api service
func NewEventSender() EventSender {
	eventBrokerHostname := flag.String("EVENT_BROKER_HOSTNAME", "", "local hostname of the event broker for eventing")
	flag.Parse()
	return EventSender{*eventBrokerHostname}
}

// Send an event containing an InventoryTransaction
func (s *EventSender) SendEvent(inventoryTransaction InventoryTransaction) {
	ctx := cloudevents.ContextWithTarget(context.Background(), s.eventBrokerHostname)
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Printf("failed to create cloudevents protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow())
	if err != nil {
		log.Printf("failed to create cloudevents client: %v", err)
	}
	// Create an Event.
	event := cloudevents.NewEvent()
	event.SetSource("api-service")
	event.SetType("service.InventoryTransaction")
	event.SetData(cloudevents.ApplicationJSON, inventoryTransaction)

	if result := c.Send(ctx, event); !cloudevents.IsACK(result) {
		log.Printf("failed to send event: %v", result)
	}
}