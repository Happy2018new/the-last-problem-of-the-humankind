package main

import (
	"time"
	"fmt"

	"github.com/pterm/pterm"
)

func SystemTestingStructureRequest() {
	tA := time.Now()

	structure, err := api.StructureRequest().GetStructure([3]int32{0, 0, 0}, [3]int32{1, 1, 1})
	_, _, err := structure.GetBlock([3]int32{0, 0, 0})
	if err != nil {
		panic("SystemTestingStructureRequest: Failed")
	}

	pterm.Success.Printfln("SystemTestingStructureRequest: PASS (Time used = %v)", time.Since(tA))
}
