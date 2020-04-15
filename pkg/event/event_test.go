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

package event_test

import (
	"github.com/kabanero-io/events-operator/pkg/event"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEvent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Suite")
}

var _ = Describe("TestEvent", func() {
	Context("TestEnqueueHandler", func() {
		queue := event.NewQueue()
		handler := event.EnqueueHandler(queue)

		It("should receive an OK status for a request with no body", func() {
			req, err := http.NewRequest("GET", "https://localhost/test-url", nil)
			Expect(err).Should(BeNil())
			rec := httptest.NewRecorder()
			handler(rec, req)
			Expect(rec.Result().StatusCode).Should(Equal(http.StatusOK))
		})

		It("should receive a StatusBadRequest status for a request with an invalid JSON payload", func() {
			payload := "{invalid: json}"
			req, err := http.NewRequest("GET", "https://localhost/test-url", strings.NewReader(payload))
			Expect(err).Should(BeNil())
			rec := httptest.NewRecorder()
			handler(rec, req)
			Expect(rec.Result().StatusCode).Should(Equal(http.StatusBadRequest))
		})

		It("should receive an OK  status for a request with a valid JSON payload", func() {
			payload := `{"data": "hello world"}`
			req, err := http.NewRequest("GET", "https://localhost/test-url", strings.NewReader(payload))
			Expect(err).Should(BeNil())
			rec := httptest.NewRecorder()
			handler(rec, req)
			Expect(rec.Result().StatusCode).Should(Equal(http.StatusOK))
		})

	})
})



