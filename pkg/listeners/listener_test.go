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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/kabanero-io/events-operator/pkg/eventenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestListener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Listener Suite")
}

var _ = Describe("TestListener", func() {

	var (
		lm 			*ListenerManagerDefault
		info 		*listenerInfo
		rec 		*httptest.ResponseRecorder
		handler		http.HandlerFunc
	)

	BeforeEach(func() {
		lm = &ListenerManagerDefault{
			listeners: make(map[int32]*listenerInfo),
		}

		info = &listenerInfo{
			port: 9080,
			key: "test-key",
			handler: func(env *eventenv.EventEnv, message map[string]interface{}, key string, url *url.URL) error {
				return nil
			},
			env: nil,
			queue: NewQueue(),
		}

		rec = httptest.NewRecorder()
		handler = listenerHandler(info)
	})

	Context("TestAddListener", func() {
		It("should add a listener on port without any error", func() {
			err := lm.addListener(9080, info)
			Expect(err).Should(BeNil())
		})

		It("should fail when trying to add a listener on port in used", func() {
			err := lm.addListener(9080, info)
			Expect(err).Should(BeNil())
			err = lm.addListener(9080, info)
			Expect(err).Should(Not(BeNil()))
		})

	})

	Context("TestListenerHandler", func() {
		It("should create a request with an empty body and receive OK status", func() {
			req, err := http.NewRequest("GET", "https://localhost/test-url", nil)
			Expect(err).Should(BeNil())
			rec = httptest.NewRecorder()
			handler(rec, req)
			res := rec.Result()
			Expect(res.StatusCode).Should(Equal(http.StatusOK))
		})

		It("should create a request with an invalid JSON payload and receive BadRequest status", func() {
			payload := "{invalid: json}"
			req, err := http.NewRequest("GET", "https://localhost/test-url", strings.NewReader(payload))
			Expect(err).Should(BeNil())
			rec = httptest.NewRecorder()
			handler(rec, req)
			res := rec.Result()
			Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))
		})

		It("should create a request with a valid JSON payload and receive OK status", func() {
			payload := `{"data": "hello world"}`
			req, err := http.NewRequest("GET", "https://localhost/test-url", strings.NewReader(payload))
			Expect(err).Should(BeNil())
			rec = httptest.NewRecorder()
			handler(rec, req)
			res := rec.Result()
			Expect(res.StatusCode).Should(Equal(http.StatusOK))
		})

	})
})


