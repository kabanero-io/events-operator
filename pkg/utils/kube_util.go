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

package utils

import (
	"bytes"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	//"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"strings"
)

/* Kubernetes and Kabanero yaml constants*/
const (
	V1                         = "v1"
	V1ALPHA1                   = "v1alpha1"
	KABANEROIO                 = "kabanero.io"
	KABANEROS                  = "kabaneros"
	ANNOTATIONS                = "annotations"
	DATA                       = "data"
	URL                        = "url"
	USERNAME                   = "username"
	PASSWORD                   = "password"
	SECRETS                    = "secrets"
	SPEC                       = "spec"
	COLLECTIONS                = "collections"
	REPOSITORIES               = "repositories"
	ACTIVATEDEFAULTCOLLECTIONS = "activateDefaultCollections"
	METADATA                   = "metadata"

	maxLabelLength = 63  // max length of a label in Kubernetes
	maxNameLength  = 253 // max length of a name in Kubernetes
)

// NewKubeConfig Creates a new kube config
func NewKubeConfig(masterURL, kubeconfigPath string) (*rest.Config, error) {
	var cfg *rest.Config
	var err error

	if masterURL != "" && kubeconfigPath != "" {
		// running outside of Kube cluster
		klog.Infof("Starting Kabanero listener outside of cluster\n")
		klog.Infof("  masterURL: %s\n", masterURL)
		klog.Infof("  kubeconfig: %s\n", kubeconfigPath)
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	} else {
		// running inside the Kube cluster
		klog.Infof("Starting Kabanero webhook status controller inside cluster\n")
		cfg, err = rest.InClusterConfig()
	}

	return cfg, err
}

// NewKabConfig Creates a new kabanero client config
/*
func NewKabConfig(masterURL, kubeconfigPath string) (*rest.Config, error) {
	crdConfig, err := NewKubeConfig(masterURL, kubeconfigPath)
	if err != nil {
		return nil, err
	}
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: KABANEROIO, Version: "v1alpha2"}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	return crdConfig, nil
}
*/

// NewKubeClient Creates a new kube client
func NewKubeClient(kubeConfig *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(kubeConfig)
}

// NewDynamicClient Creates a new dynamic client
func NewDynamicClient(kubeConfig *rest.Config) (dynamic.Interface, error) {
	return dynamic.NewForConfig(kubeConfig)
}

/*
ToDomainName Convert a name to domain name format. The name must:
  - Start with [a-z0-9]. If not, "0" is prepended.
  - lower case. If not, lower case is used.
  - contain only '.', '-', and [a-z0-9]. If not, "." is used insteaad.
  - end with alpha numeric characters. Otherwise, '0' is appended
  - can't have consecutive '.'.  Consecutive ".." is substituted with ".".
Returns the empty string if the name is empty after conversion
*/
func ToDomainName(name string) string {
	maxLength := maxNameLength
	name = strings.ToLower(name)
	ret := bytes.Buffer{}
	chars := []byte(name)
	for i, ch := range chars {
		if i == 0 {
			// first character must be [a-z0-9]
			if (ch >= 'a' && ch <= 'z') ||
				(ch >= '0' && ch <= '9') {
				ret.WriteByte(ch)
			} else {
				ret.WriteByte('0')
				if isValidDomainNameChar(ch) {
					ret.WriteByte(ch)
				} else {
					ret.WriteByte('.')
				}
			}
		} else {
			if isValidDomainNameChar(ch) {
				ret.WriteByte(ch)
			} else {
				ret.WriteByte('.')
			}
		}
	}

	// change all ".." to ".
	retStr := ret.String()
	for strings.Index(retStr, "..") > 0 {
		retStr = strings.ReplaceAll(retStr, "..", ".")
	}

	strLen := len(retStr)
	if strLen == 0 {
		return retStr
	}
	if strLen > maxLength {
		strLen = maxLength
		retStr = retStr[0:strLen]
	}
	ch := retStr[strLen-1]
	if (ch >= 'a' && ch <= 'z') ||
		(ch >= '0' && ch <= '9') {
		// last char is alphanumeric
		return retStr
	}
	if strLen < maxLength-1 {
		//  append alphanumeric
		return retStr + "0"
	}
	// replace last char to be alphanumeric
	return retStr[0:strLen-2] + "0"
}

/*
ToLabelName Convert the  name part of  a label. The name must:
  - Start with [a-z0-9A-Z]. If not, "0" is prepended.
  - End with [a-z0-9A-Z]. If not, "0" is appended
  - Intermediate characters can only be: [a-z0-9A-Z] or '_', '-', and '.' If not, '.' is used.
  - be maximum maxLabelLength characters long
*/
func ToLabelName(name string) string {
	chars := []byte(name)
	ret := bytes.Buffer{}
	for i, ch := range chars {
		if i == 0 {
			// first character must be [a-z0-9]
			if (ch >= 'a' && ch <= 'z') ||
				(ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') {
				ret.WriteByte(ch)
			} else {
				ret.WriteByte('0')
				if isValidLabelChar(ch) {
					ret.WriteByte(ch)
				} else {
					ret.WriteByte('.')
				}
			}
		} else {
			if isValidLabelChar(ch) {
				ret.WriteByte(ch)
			} else {
				ret.WriteByte('.')
			}
		}
	}

	retStr := ret.String()
	strLen := len(retStr)
	if strLen == 0 {
		return retStr
	}
	if strLen > maxLabelLength {
		strLen = maxLabelLength
		retStr = retStr[0:strLen]
	}

	ch := retStr[strLen-1]
	if (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') {
		// last char is alphanumeric
		return retStr
	} else if strLen < maxLabelLength-1 {
		//  append alphanumeric
		return retStr + "0"
	} else {
		// replace last char to be alphanumeric
		return retStr[0:strLen-2] + "0"
	}
}

// ToLabel Convert a string to Kubernetes label format.
func ToLabel(input string) string {
	slashIndex := strings.Index(input, "/")
	var prefix, label string
	if slashIndex < 0 {
		prefix = ""
		label = input
	} else if slashIndex == len(input)-1 {
		prefix = input[0:slashIndex]
		label = ""
	} else {
		prefix = input[0:slashIndex]
		label = input[slashIndex+1:]
	}

	newPrefix := ToDomainName(prefix)
	newLabel := ToLabelName(label)
	ret := ""
	if newPrefix == "" {
		if newLabel == "" {
			// shouldn't happen
			newLabel = "nolabel"
		} else {
			ret = newLabel
		}
	} else if newLabel == "" {
		ret = newPrefix
	} else {
		ret = newPrefix + "/" + newLabel
	}
	return ret
}

/* @Return true if character is valid for a domain name */
func isValidDomainNameChar(ch byte) bool {
	return ch == '.' || ch == '-' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= '0' && ch <= '9')
}

func isValidLabelChar(ch byte) bool {
	return ch == '.' || ch == '-' || (ch == '_') ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9')
}
