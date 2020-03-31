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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/kabanero-io/events-operator/pkg/eventenv"
	"k8s.io/klog"
)

const (
	// MessageHeader is the message key containing the request's headers
	MessageHeader = "header"
	// MessageBody is the message key containing the request's payload
	MessageBody = "body"

	defaultHttpPort  int32 = 9080
	defaultHttpsPort int32 = 9443

	defaultTLSCertPath = "/etc/tls/tls.crt"
	defaultTLSKeyPath  = "/etc/tls/tls.key"
)

type listenerInfo struct {
	port    int32
	key     string
	handler eventenv.ListenerHandler
	env     *eventenv.EventEnv
	queue   Queue
}

type ListenerManagerDefault struct {
	listeners map[int32]*listenerInfo
	mutex     sync.Mutex
}

type queueElem struct {
	message map[string]interface{}
	url     *url.URL
	info    *listenerInfo
}

func NewDefaultListenerManager() eventenv.ListenerManager {
	return &ListenerManagerDefault{
		listeners: make(map[int32]*listenerInfo),
	}
}

// NewListener creates a new event listener on port 9080

func (listenerMgr *ListenerManagerDefault) NewListener(env *eventenv.EventEnv, key string, handler eventenv.ListenerHandler, options eventenv.ListenerOptions) error {
	if options.Port != 0 {
		options.Port = defaultHttpPort
	}
	port := options.Port

	klog.Infof("Starting new listener on port %v", port)

	listener := &listenerInfo{
		port:    port,
		key:     key,
		handler: handler,
		env:     env,
		queue:   NewQueue(),
	}

	if err := listenerMgr.addListener(port, listener); err != nil {
		return err
	}

	/* start listener thread */
	go func() {
		klog.Infof("Listener thread started for port %v", port)
		err := http.ListenAndServe(":"+strconv.Itoa(int(port)), listenerHandler(listener))
		if err != nil {
			klog.Errorf("Listener thread error for port %v, error: %v", port, err)
		}
		klog.Infof("Listener thread stopped for port %v", port)
	}()

	/* Start a new thread to process the queue */
	klog.Infof("Starting worker thread for listener on port %v", port)
	go processQueueWorker(listener.queue)

	return nil
}

func (listenerMgr *ListenerManagerDefault) NewListenerTLS(env *eventenv.EventEnv, key string, handler eventenv.ListenerHandler, options eventenv.ListenerOptions) error {
	if options.Port == 0 {
		options.Port = defaultHttpsPort
	}
	port := options.Port

	if options.TLSCertPath == "" {
		options.TLSCertPath = defaultTLSCertPath
	}

	if options.TLSKeyPath == "" {
		options.TLSKeyPath = defaultTLSKeyPath
	}

	klog.Infof("Starting TLS listener on port %v", port)

	if _, err := os.Stat(options.TLSCertPath); os.IsNotExist(err) {
		klog.Fatalf("TLS certificate '%s' not found: %v", options.TLSCertPath, err)
		return err
	}
	if _, err := os.Stat(options.TLSKeyPath); os.IsNotExist(err) {
		klog.Fatalf("TLS private key '%s' not found: %v", options.TLSKeyPath, err)
		return err
	}

	listener := &listenerInfo{
		port:    port,
		key:     key,
		handler: handler,
		env:     env,
		queue:   NewQueue(),
	}

	if err := listenerMgr.addListener(port, listener); err != nil {
		return err
	}

	/* start listener thread */
	go func() {
		klog.Infof("TLS Listener thread started for port %v", port)
		err := http.ListenAndServeTLS(":"+strconv.Itoa(int(port)), options.TLSCertPath, options.TLSKeyPath, listenerHandler(listener))
		if err != nil {
			klog.Infof("TLS Listener thread error for port %v, error: %v  ", port, err)
		}
		klog.Infof("TLS Listener thread ended for port %v", port)
	}()

	/* Start a new thread to process the queue */
	klog.Infof("Starting worker thread for listener on port %v", port)
	go processQueueWorker(listener.queue)

	return nil
}

/* Process events on queue */
func processQueueWorker(queue Queue) {
	klog.Info("Worker thread started to process messages.")
	for {
		interf := queue.Dequeue()
		qElem := interf.(*queueElem)
		klog.Infof("Worker thread processing url: %v, message: %v", qElem.url.String(), qElem.message)
		err := (qElem.info.handler)(qElem.info.env, qElem.message, qElem.info.key, qElem.url)
		if err != nil {
			klog.Errorf("Worker thread error: url: %v, error: %v", qElem.url.String(), err)
			return
		}
		klog.Infof("Worker thread completed processing url: %v", qElem.url.String())
	}
}

func (listenerMgr *ListenerManagerDefault) addListener(port int32, listener *listenerInfo) error {
	listenerMgr.mutex.Lock()
	defer listenerMgr.mutex.Unlock()

	if _, exists := listenerMgr.listeners[port]; exists {
		return fmt.Errorf("listener on port %v already exists", port)
	}

	listenerMgr.listeners[port] = listener
	return nil
}

/* Event listener */
func listenerHandler(listener *listenerInfo) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {

		header := req.Header
		klog.Infof("Received request. Header: %v", header)

		var bodyMap map[string]interface{}

		if req.Body != nil {
			defer req.Body.Close()
			bytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				klog.Errorf("Listener can not read body. Error: %v", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			klog.Infof("Listener received body: %v", string(bytes))
			err = json.Unmarshal(bytes, &bodyMap)
			if err != nil {
				klog.Errorf("Unable to unmarshal json body: %v", err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			klog.Info("Request did not have a body")
		}

		message := make(map[string]interface{})
		message[MessageHeader] = map[string][]string(header)
		message[MessageBody] = bodyMap

		listener.queue.Enqueue(&queueElem{
			message: message,
			url:     req.URL,
			info:    listener,
		})

		writer.WriteHeader(http.StatusOK)
	}
}
