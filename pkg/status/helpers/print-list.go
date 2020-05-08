package helpers

import (
	"fmt"
	eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
)

func PrintList(list []eventsv1alpha1.EventStatusSummary) {
	fmt.Printf("List: ")
	for _, entry := range list {
		fmt.Printf("%v ", entry.Operation)
	}
	fmt.Printf("\n")
}
