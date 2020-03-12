// +build unit_test

/*
Copyright 2020 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package listeners

import (
	"github.com/kabanero-io/events-operator/pkg/eventenv"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAddListener(t *testing.T) {
	lm := &ListenerManagerDefault{
		listeners: make(map[int]*listenerInfo),
	}

	info := &listenerInfo{
		port: 9080,
		key: "test-key",
		handler: func(env *eventenv.EventEnv, message map[string]interface{}, key string, url *url.URL) error {
			return nil
		},
		env: nil,
		queue: NewQueue(),
	}

	if err := lm.addListener(9080, info); err != nil {
		t.Fatalf("unable to add initial listener on port 9080: %v", err)
	}

	if err := lm.addListener(9080, info); err == nil {
		t.Fatal("adding second listener on port 9080 should've failed")
	}
}

func TestListenerHandler(t *testing.T) {
	info := &listenerInfo{
		port: 9080,
		key: "test-key",
		handler: func(env *eventenv.EventEnv, message map[string]interface{}, key string, url *url.URL) error {
			return nil
		},
		env: nil,
		queue: NewQueue(),
	}

	handler := listenerHandler(info)

	// Test request with no body
	req, err := http.NewRequest("GET", "https://localhost/test-url", nil)
	if err != nil {
		t.Fatalf("could not create request with no body: %v", err)
	}
	rec := httptest.NewRecorder()
	handler(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 'OK'; got %v", res.Status)
	}

	// Test request with invalid JSON payload
	payload := "{invalid: json}"
	req, err = http.NewRequest("GET", "https://localhost/test-url", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("could not create request with invalid JSON payload: %v", err)
	}
	rec = httptest.NewRecorder()
	handler(rec, req)
	res = rec.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 'Bad Request; got %v", res.Status)
	}

	// Test request with valid JSON
	payload = `{"data": "hello world"}`
	req, err = http.NewRequest("GET", "https://localhost/test-url", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("could not create request with valid JSON payload: %v", err)
	}
	rec = httptest.NewRecorder()
	handler(rec, req)
	res = rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 'OK'; got %v", res.Status)
	}
}
