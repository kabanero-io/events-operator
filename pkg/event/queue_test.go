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

package event

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Queue Suite")
}

var _ = Describe("TestQueue", func() {
	var queue Queue

	BeforeEach(func() {
		queue = NewQueue()
	})

	Context("QueueFunctions", func() {
		It("it should be able to enqueue data", func() {
			By("adding data from an array and checking length", func() {
				Expect(queue.Len()).Should(BeZero())
				data := []int{1, 2, 3, 4, 5}
				for i := range data {
					queue.Enqueue(i)
				}
				Expect(len(data)).Should(Equal(queue.Len()))
			})

		})

		It("It should be able to dequeue data", func() {
			By("adding data from an array and checking content", func() {
				data := []int{1, 2, 3, 4, 5}
				for i := range data {
					queue.Enqueue(data[i])
				}
				res := [5]int{}
				i := 0
				for queue.Len() != 0 {
					val := queue.Dequeue()
					res[i] = val.(int)
					i++
				}
				Expect(res).Should(Equal([5]int{1, 2, 3, 4, 5}))
			})
		})

	})

	Context("TestBlockingDequeue", func() {
		It("should be tested using concurrency with go functions", func() {
			const val = 5
			ch := make(chan int, 1)

			go func() {
				fmt.Print("Waiting for element to be enqueued\n")
				e := queue.Dequeue().(int)
				ch <- e
				close(ch)
			}()

			time.Sleep(100 * time.Millisecond)
			queue.Enqueue(val)
			e := <-ch
			Expect(e).Should(Equal(val))
		})
	})

})
