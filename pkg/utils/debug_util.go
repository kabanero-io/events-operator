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
	"encoding/json"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

func printEvent(event watch.Event) {
	klog.Infof("event type %s, object type is %T\n", event.Type, event.Object)
	printEventObject(event.Object, "    ")
}

func printEventObject(obj interface{}, indent string) {
	switch obj.(type) {
	case *unstructured.Unstructured:
		var unstructuredObj = obj.(*unstructured.Unstructured)
		// printObject(unstructuredObj.Object, indent)
		printUnstructuredJSON(unstructuredObj.Object, indent)
		klog.Infof("\n")
	default:
		klog.Infof("%snot Unstructured: type: %T val: %s\n", indent, obj, obj)
	}
}

func printUnstructuredJSON(obj interface{}, indent string) {
	data, err := json.MarshalIndent(obj, "", indent)
	if err != nil {
		klog.Fatalf("JSON Marshaling failed %s", err)
	}
	klog.Infof("%s\n", data)
}

func printObject(obj interface{}, indent string) {
	nextIndent := indent + "    "
	switch obj.(type) {
	case int:
		klog.Infof("%d", obj.(int))
	case bool:
		klog.Infof("%t", obj.(bool))
	case float64:
		klog.Infof("%f", obj.(float64))
	case string:
		klog.Infof("%s", obj.(string))
	case []interface{}:
		var arr = obj.([]interface{})
		for index, elem := range arr {
			klog.Infof("\n%sindex:%d, type %T, ", indent, index, elem)
			printObject(elem, nextIndent)
		}
	case map[string]interface{}:
		var objMap = obj.(map[string]interface{})
		for label, val := range objMap {
			klog.Infof("\n%skey: %s type: %T| ", indent, label, val)
			printObject(val, nextIndent)
		}
	default:
		klog.Infof("\n%stype: %T val: %s", indent, obj, obj)
	}
}

func printPods(pods *v1.PodList) {
	for _, pod := range pods.Items {
		klog.Infof("%s", pod.ObjectMeta.Name)
	}
}
