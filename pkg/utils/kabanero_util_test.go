// +build full_test

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

package utils_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kabanero-io/events-operator/pkg/utils"
	"github.com/kabanero-io/kabanero-operator/pkg/apis/kabanero/v1alpha2"
	"io"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	fakerest "k8s.io/client-go/rest/fake"
	"net/http"
	"testing"
)

func defaultHeaders() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return header
}

func bytesBody(bodyBytes []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(bodyBytes))
}

func TestGetTriggerInfo(t *testing.T) {
	const (
		triggerID     = "incubator"
		triggerSHA256 = "0123456789abcdef"
		triggerURL    = "https://example.com/incubator.trigger.tar.gz"
	)

	kabanero := v1alpha2.KabaneroList{
		Items: []v1alpha2.Kabanero{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Kabanero",
					APIVersion: "v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "kabanero-no-trigger",
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Kabanero",
					APIVersion: "v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "kabanero-with-trigger",
				},
				Spec: v1alpha2.KabaneroSpec{
					Triggers: []v1alpha2.TriggerSpec{
						{
							Id:     triggerID,
							Sha256: triggerSHA256,
							Https:  v1alpha2.HttpsProtocolFile{
								Url: triggerURL,
							},
						},
					},
				},
			},
		},
	}

	resourcePaths := map[string]v1alpha2.KabaneroList{
		"/apis/kabanero.io/v1alpha2/namespaces/kabanero/kabaneros": kabanero,
	}

	fakeReqHandler := func(req *http.Request) (*http.Response, error) {
		kab, isKabPath := resourcePaths[req.URL.Path]

		if !isKabPath {
			return nil, fmt.Errorf("unexpected request for URL %q with method %q", req.URL.String(), req.Method)
		}

		switch req.Method {
		case "GET":
			res, err := json.Marshal(kab)
			if err != nil {
				return nil, err
			}
			return &http.Response{StatusCode: http.StatusOK, Header: defaultHeaders(), Body: bytesBody(res)}, nil
		default:
			return nil, fmt.Errorf("unexpected request for URL %q with method %q", req.URL.String(), req.Method)
		}
	}

	fakeClient := &fakerest.RESTClient{
		Client:               fakerest.CreateHTTPClient(fakeReqHandler),
		NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		GroupVersion:         schema.GroupVersion{},
		VersionedAPIPath:     "/apis/kabanero.io/v1alpha2",
	}

	url, sha256, err := utils.GetTriggerInfo(fakeClient, "kabanero")
	if err != nil {
		t.Errorf("failed to get trigger info: %v", err)
		t.Fail()
	}

	if url.String() != triggerURL {
		t.Errorf("trigger URL is incorrect; got %s, want %s", url, triggerURL)
		t.Fail()
	}

	if sha256 != triggerSHA256 {
		t.Errorf("trigger checksum is incorrect; got %s, want %s", sha256, triggerSHA256)
		t.Fail()
	}

	t.Logf("received expected url and checksum: url -> %s, sha256 -> %s", url.String(), sha256)
}
