package controller

import (
	"github.com/events-operator/pkg/controller/eventmediator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, eventmediator.Add)
}
