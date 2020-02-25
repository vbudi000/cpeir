package controller

import (
	"github.ibm.com/CASE/cpeir/pkg/controller/cpeir"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cpeir.Add)
}
