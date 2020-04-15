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

import "fmt"

type ResourceNotFound struct {
	collection string
	id         string
}

func (e ResourceNotFound) Error() string {
	return fmt.Sprintf("resource %q not found in collection %q", e.id, e.collection)
}

func ItemNotFound(id string) *ResourceNotFound {
	return &ResourceNotFound{collection: "items", id: id}
}

func LocationNotFound(id string) *ResourceNotFound {
	return &ResourceNotFound{collection: "locations", id: id}
}

func InventoryTransactionNotFound(id string) *ResourceNotFound {
	return &ResourceNotFound{collection: "inventoryTransactions", id: id}
}

type ResourceConflict struct {
	collection string
	id         string
}

func (e ResourceConflict) Error() string {
	return fmt.Sprintf("concurrent transaction ongoing conflicting with resource %q in collection %q", e.id, e.collection)
}
