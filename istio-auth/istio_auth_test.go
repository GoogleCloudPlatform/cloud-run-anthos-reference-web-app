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

// These tests verify the Istio authorization policies.
// Note that no actual requests will succeed as the "Host:" header is never set.
// However, we can still check for 404s for requests that managed to pass the auth policy.

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var app *firebase.App
var tokens = map[string]string{"admin": "", "worker": "", "other": ""}

var host = "http://" + os.Getenv("HOST_IP")
var apiKey = os.Getenv("API_KEY")

var methods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

var paths = []string{
	"/api",
	"/api/alerts",
	"/api/alerts/id",
	"/api/inventoryTransactions",
	"/api/inventoryTransactions/id",
	"/api/items",
	"/api/items/id",
	"/api/items/id/inventory",
	"/api/items/id/inventoryTransactions",
	"/api/locations",
	"/api/locations/id",
	"/api/locations/id/inventory",
	"/api/locations/id/inventoryTransactions",
	"/api/users",
	"/api/users/id",
}

func checkResponse(t *testing.T, method, path, token string, want int) {
	client := http.DefaultClient
	req, err := http.NewRequest(method, host+path, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
		return
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("%v %v returned unexpected err: %v", method, path, err)
		return
	}
	if resp.StatusCode != want {
		t.Errorf("%v %v returned unexpected code %v, want %v", method, path, resp.StatusCode, want)
	}
}

// Tests api-origin-auth and require-valid-token policies
func TestDenyUnauthenticatedUsers(t *testing.T) {
	var users = []struct {
		desc  string
		token string
	}{
		{desc: "unauthenticated user"},
		{desc: "bad token", token: "badToken"},
	}
	want := http.StatusUnauthorized

	for _, u := range users {
		token := u.token
		t.Run(u.desc, func(t *testing.T) {
			for _, m := range methods {
				for _, p := range paths {
					checkResponse(t, m, p, token, want)
				}
			}
		})
	}
}

// Tests api-allow-get policy
func TestAllowGetPolicy(t *testing.T) {
	m := http.MethodGet
	want := http.StatusNotFound
	for _, u := range []string{"admin", "worker", "other"} {
		t.Run(u, func(t *testing.T) {
			token := tokens[u]
			for _, p := range paths {
				checkResponse(t, m, p, token, want)
			}
		})
	}
}

// Tests api-allow-admin and api-allow-workers policy
func TestRoleBasedPolicies(t *testing.T) {
	for _, u := range []string{"admin", "worker", "other"} {
		token := tokens[u]
		defaultWant := http.StatusForbidden
		if u == "admin" {
			defaultWant = http.StatusNotFound
		}
		t.Run(u, func(t *testing.T) {
			for _, m := range []string{http.MethodDelete, http.MethodPost, http.MethodPut} {
				for _, p := range paths {
					want := defaultWant
					if p == "/api/inventoryTransactions" && u == "worker" && m == "POST" {
						want = http.StatusNotFound
					}
					checkResponse(t, m, p, token, want)
				}
			}
		})
	}
}

func generateIDToken(client *auth.Client, role string) string {
	uid := role + "@example.com"
	claims := map[string]interface{}{"role": role}

	// Create a custom token:
	// https://firebase.google.com/docs/auth/admin/create-custom-tokens#create_custom_tokens_using_the_firebase_admin_sdk
	token, err := client.CustomTokenWithClaims(context.Background(), uid, claims)
	if err != nil {
		log.Fatalf("error minting custom token: %v", err)
	}

	// Exchange for an ID token:
	// https://firebase.google.com/docs/reference/rest/auth/#section-verify-custom-token
	payload := map[string]interface{}{"token": token, "returnSecureToken": true}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("error encoding json: %v", err)
	}

	resp, err := http.Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key="+apiKey, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("error getting ID token: %v", err)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v", body)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unable to sign in with custom token (status %v): %v", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatalf("error decoding json: %v", err)
	}
	return payload["idToken"].(string)
}

func TestMain(m *testing.M) {
	var err error
	app, err = firebase.NewApp(context.Background(), &firebase.Config{ServiceAccountID: os.Getenv("FIREBASE_SA")})
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}

	for r := range tokens {
		tokens[r] = generateIDToken(client, r)
	}

	os.Exit(m.Run())
}
