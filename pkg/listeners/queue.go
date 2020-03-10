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
    "k8s.io/klog"
    "container/list"
    "sync"
)

type Queue  interface {
    Enqueue(elem interface{})
    Dequeue() interface{}
    Len() int
}

type queueImpl struct {
    cond *sync.Cond
    list *list.List
}

func NewQueue() Queue{
    return &queueImpl {
        cond: sync.NewCond(&sync.Mutex{}),
        list: list.New(),
    }
}

func (qImpl *queueImpl) Enqueue(elem interface{}) {
    klog.Info("Enqueue called")
    qImpl.cond.L.Lock()
    defer qImpl.cond.L.Unlock()
    qImpl.list.PushBack(elem)

    /* wake anyone waiting to dequeue */
    qImpl.cond.Signal()
}

func (qImpl *queueImpl) Dequeue() interface{} {
    klog.Info("Dequeue called")
    qImpl.cond.L.Lock()
    defer qImpl.cond.L.Unlock()

    /* wait until there is something in the queue */
    for qImpl.list.Len() == 0 {
         qImpl.cond.Wait()
    }

    return qImpl.list.Remove(qImpl.list.Front())
}

func (qImpl *queueImpl) Len() int {
    qImpl.cond.L.Lock()
    defer qImpl.cond.L.Unlock()

    return qImpl.list.Len()
}
