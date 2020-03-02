package controller

import (
	"github.com/kabanero-io/events-operator/pkg/controller/eventconnections"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, eventconnections.Add)
}
