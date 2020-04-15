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

import "net/http"

// EncodeJSONStatus calls EncodeJSONResponse with the status and message wrapped in a Status message.
func EncodeJSONStatus(status int, message string, w http.ResponseWriter) error {
	response := &Status{
		Status:  int32(status),
		Message: message,
	}
	return EncodeJSONResponse(response, &status, w)
}
