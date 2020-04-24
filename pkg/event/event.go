package event

import (
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"net/url"

	"encoding/json"
)

const (
	// MessageHeader is the message key containing the request's headers
	MessageHeader = "header"
	// MessageBody is the message key containing the request's payload
	MessageBody = "body"
)

// Event contains the destination URL, headers, and a body
type Event struct {
	URL    *url.URL
    RemoteAddr string
	Header map[string][]string
	Body   map[string]interface{}
}

// Types of events
type Type int

const (
	TypeOther Type = iota
	TypeGitHub
)

// A handler that responds to an event
type Handler func(event *Event) error

/* Event listener listens for REST requests and enqueues a message consisting of the request's headers and payloads. */
func EnqueueHandler(queue Queue) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		klog.Infof("Received request. Header: %v", r.Header)

		var bodyMap map[string]interface{}

		if r.Body != nil {
			bytes, err := ioutil.ReadAll(r.Body)
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

		queue.Enqueue(&Event{
			URL:    r.URL,
            RemoteAddr: r.RemoteAddr,
			Header: r.Header,
			Body:   bodyMap,
		})

		writer.WriteHeader(http.StatusOK)
	}
}

/* ProcessQueueWorker processes events on the Queue */
func ProcessQueueWorker(queue Queue, handler Handler) {
	klog.Info("Worker thread started to process messages.")
	for {
		event := queue.Dequeue().(*Event)
		// TODO: Remove this later or only include when very verbose logging is enabled
		klog.Infof("Worker thread processing url: %s, header: %v, body: %v", event.URL, event.Header, event.Body)
		err := handler(event)
		if err != nil {
			klog.Errorf("Worker thread error: url: %s, error: %v", event.URL, err)
			continue
		}
		klog.Infof("Worker thread completed processing url: %s", event.URL)
	}
}
