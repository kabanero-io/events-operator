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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestListener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Listener Suite")
}

var _ = Describe("TestListener", func() {
	Context("TestAddListener", func() {
		var lm *ListenerManagerDefault
		BeforeEach(func() {
			lm = &ListenerManagerDefault{
				listeners: make(map[int32]*listenerInfo),
			}
		})

		info := &listenerInfo{
			port: 9080,
		}

		It("should add a listener on Port without any error", func() {
			err := lm.addListener(9080, info)
			Expect(err).Should(BeNil())
		})

		It("should fail when trying to add a listener on port in use", func() {
			err := lm.addListener(9080, info)
			Expect(err).Should(BeNil())
			err = lm.addListener(9080, info)
			Expect(err).Should(Not(BeNil()))
		})
	})
})


