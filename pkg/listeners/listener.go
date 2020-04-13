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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"k8s.io/klog"
)

const (
	defaultHttpPort  int32 = 9080
	defaultHttpsPort int32 = 9443

	defaultTLSCertPath = "/etc/tls/tls.crt"
	defaultTLSKeyPath  = "/etc/tls/tls.key"
)

type ListenerManager interface {
	/* Create a new listener. */
	NewListener(handler http.Handler, options ListenerOptions) error

	/* Create a new TLS listener with TLS. Call the handler on every message received */
	NewListenerTLS(handler http.Handler, options ListenerOptions) error
}

type listenerInfo struct {
	port int32
}

type ListenerOptions struct {
	Port        int32
	TLSCertPath string
	TLSKeyPath  string
}

type ListenerManagerDefault struct {
	listeners map[int32]*listenerInfo
	mutex     sync.Mutex
}

// NewDefaultListenerManager creates a new ListenerManager
func NewDefaultListenerManager() ListenerManager {
	return &ListenerManagerDefault{
		listeners: make(map[int32]*listenerInfo),
	}
}

// NewListener creates a new event listener
func (listenerMgr *ListenerManagerDefault) NewListener(handler http.Handler, options ListenerOptions) error {
	if options.Port != 0 {
		options.Port = defaultHttpPort
	}
	port := options.Port

	klog.Infof("Starting new listener on Port %v", port)

	listener := &listenerInfo{
		port: port,
	}

	if err := listenerMgr.addListener(port, listener); err != nil {
		return err
	}

	/* start listener thread */
	go func() {
		klog.Infof("Listener thread started for Port %v", port)
		err := http.ListenAndServe(":"+strconv.Itoa(int(port)), handler)
		if err != nil {
			klog.Errorf("Listener thread error for Port %v, error: %v", port, err)
		}
		klog.Infof("Listener thread stopped for Port %v", port)
	}()

	return nil
}

// NewListener creates a new HTTPS event listener
func (listenerMgr *ListenerManagerDefault) NewListenerTLS(handler http.Handler, options ListenerOptions) error {
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

	klog.Infof("Starting TLS listener on Port %v", port)

	if _, err := os.Stat(options.TLSCertPath); os.IsNotExist(err) {
		klog.Fatalf("TLS certificate '%s' not found: %v", options.TLSCertPath, err)
		return err
	}
	if _, err := os.Stat(options.TLSKeyPath); os.IsNotExist(err) {
		klog.Fatalf("TLS private key '%s' not found: %v", options.TLSKeyPath, err)
		return err
	}

	listener := &listenerInfo{
		port: port,
	}

	if err := listenerMgr.addListener(port, listener); err != nil {
		return err
	}

	/* start listener thread */
	go func() {
		klog.Infof("TLS Listener thread started for Port %v", port)
		err := http.ListenAndServeTLS(":"+strconv.Itoa(int(port)), options.TLSCertPath, options.TLSKeyPath, handler)
		if err != nil {
			klog.Infof("TLS Listener thread error for Port %v, error: %v  ", port, err)
		}
		klog.Infof("TLS Listener thread ended for Port %v", port)
	}()

	return nil
}

func (listenerMgr *ListenerManagerDefault) addListener(port int32, listener *listenerInfo) error {
	listenerMgr.mutex.Lock()
	defer listenerMgr.mutex.Unlock()

	if _, exists := listenerMgr.listeners[port]; exists {
		return fmt.Errorf("listener on Port %v already exists", port)
	}

	listenerMgr.listeners[port] = listener
	return nil
}
