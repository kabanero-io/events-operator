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
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue()

	data := []int{1, 2, 3, 4, 5}

	for i := range data {
		q.Enqueue(i)
	}

	if q.Len() != len(data) {
		t.Errorf("expected queue to have %d items but has %d items", len(data), q.Len())
	}

	for i := range data {
		e := q.Dequeue().(int)
		if i != e {
			t.Errorf("expected dequeued element to be %d; got %d", i, e)
		}
	}

	if q.Len() != 0 {
		t.Errorf("queue should be empty but has %d elements: ", q.Len())
		for q.Len() != 0 {
			e := q.Dequeue().(int)
			t.Errorf("\t%d", e)
		}
	}
}

func TestBlockingDequeue(t *testing.T) {
	q := NewQueue()
	ch := make(chan int, 1)
	const val= 5
	go func() {
		t.Log("Waiting for element to be enqueued")
		e := q.Dequeue().(int)
		ch <- e
	}()
	time.Sleep(100 * time.Millisecond)
	q.Enqueue(val)

	e := <-ch
	if e != val {
		t.Errorf("expected value %d; got %d", val, e)
	}
}
