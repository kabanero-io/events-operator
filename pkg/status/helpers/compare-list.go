package helpers

import (
	eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
)

func CompareList(list1 []eventsv1alpha1.EventStatusSummary, list2[]eventsv1alpha1.EventStatusSummary) bool {
	if len(list1) != len(list1) {
		return false
	}

	for index, elem := range list1{
		if !elem.Equals(&list2[index]) {
			return false
		}
	}
	return true
}
