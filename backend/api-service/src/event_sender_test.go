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
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/client/test"
	cecontext "github.com/cloudevents/sdk-go/v2/context"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

const id = "123"
const eventBrokerHostname = "http://event-broker-hostname"
const incorrectTargetErrorStringForTest = "Incorrect target for test;"

var inventoryTransaction = InventoryTransaction{
	Id:     id,
	ItemId: "456",
	Action: "ADD",
}

// A Sender whose Send call fails if the context does not have eventBrokerHostname as the target
type TargetCheckingSender struct{}

func (t TargetCheckingSender) Send(ctx context.Context, m binding.Message, transformers ...binding.Transformer) (err error) {
	target := cecontext.TargetFrom(ctx)
	if target.String() != eventBrokerHostname {
		return fmt.Errorf(incorrectTargetErrorStringForTest+" expected %s, got %s", eventBrokerHostname, target.String())
	}
	return nil
}

func TestNewBrokerEventSender(t *testing.T) {
	sender, err := NewBrokerEventSender(eventBrokerHostname)
	if err != nil {
		t.Fatalf("Expected success; got %s", err)
	}
	if (sender == nil) {
		t.Fatalf("Expected a sender")
	}

	s, ok := sender.(*brokerEventSender)
	if !ok {
		t.Fatalf("Expected to get a brokerEventSender %s, %s", s, sender)
	}
	if s.client == nil {
		t.Errorf("Expected a client")
	}
	if eventBrokerHostname != s.hostname {
		t.Errorf("Expected %s; got %s", eventBrokerHostname, s.hostname)
	}
}

func TestSendInventoryTransactionEventSuccess(t *testing.T) {
	expectedEvent := cloudevents.NewEvent()
	expectedEvent.SetID(id)
	expectedEvent.SetSource("api-service")
	expectedEvent.SetType("service.InventoryTransaction")
	expectedEvent.SetData(cloudevents.ApplicationJSON, inventoryTransaction)

	mockClient, eventCh := test.NewMockSenderClient(t, 1)
	brokerEventSender := brokerEventSender{mockClient, eventBrokerHostname}
	err := brokerEventSender.SendInventoryTransactionEvent(inventoryTransaction)
	if err != nil {
		t.Errorf("Expected success; got %s", err)
	}
	got := <-eventCh
	if !cmp.Equal(got, expectedEvent) {
		t.Errorf("Expected = %s; got %s", expectedEvent, got)
	}
}

func TestSendInventoryTransactionEventError(t *testing.T) {
	mockClient, _ := client.New(TargetCheckingSender{})
	brokerEventSender := brokerEventSender{mockClient, "bad host name"}
	err := brokerEventSender.SendInventoryTransactionEvent(inventoryTransaction)
	if err == nil {
		t.Fatalf("Expected an error.")
	}
	if !strings.HasPrefix(err.Error(), incorrectTargetErrorStringForTest) {
		t.Errorf("Expected a different error; got %s", err)
	}
}

func TestSendInventoryTransactionEventCheckHostname(t *testing.T) {
	mockClient, _ := client.New(TargetCheckingSender{})
	brokerEventSender := brokerEventSender{mockClient, eventBrokerHostname}
	err := brokerEventSender.SendInventoryTransactionEvent(inventoryTransaction)
	if err != nil {
		t.Errorf("Expected success; got %s", err)
	}
}
