package main

import (
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/pterm/pterm"
)

var (
	c         *client.Client
	resources *resources_control.Resources
	api       *game_interface.GameInterface
)

func main() {
	tA := time.Now()

	SystemTestingLogin()
	defer func() {
		c.Conn().Close()
		time.Sleep(time.Second)
	}()

	SystemTestingCommands()
	SystemTestingStructureRequest()
	SystemTestingUUIDSafeString()
	SystemTestingSetblock()
	SystemTestingReplaceitem()
	SystemTestingQuerytarget()
	SystemTestingStructrueBackup()
	SystemTestingBotClick()
	SystemTestingContainer()
	SystemTestingItemStackOperation()
	SystemTestingItemCopy()
	SystemTestingItemTransition()

	pterm.Success.Printfln("System Testing: ALL PASS (Time used = %v)", time.Since(tA))
}
