// +build unit_test

package listeners

import (
	"fmt"
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
			if url.Scheme != "https" {
				return fmt.Errorf("insecure http requests are not accepted")
			}
			return nil
		},
		env: nil,
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
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 'Bad Request'; got %v", res.Status)
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
	defer res.Body.Close()
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
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		t.Errorf("expected status 'Accepted'; got %v", res.Status)
	}

	// Test error checking of handler
	payload = `{"data": "hello world"}`
	req, err = http.NewRequest("GET", "http://localhost/test-url", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	rec = httptest.NewRecorder()
	handler(rec, req)
	res = rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 'Internal Server Error'; got %v", res.Status)
	}
}
