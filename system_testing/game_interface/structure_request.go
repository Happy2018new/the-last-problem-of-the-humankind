package main

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

func SystemTestingStructureRequest() {
	tA := time.Now()
    
    // 获取 (0, 0, 0) 处 3x3x3 结构
    origin := [3]int32{0, 0, 0}
    size := [3]int32{3, 3, 3}
    structure, err := api.StructureRequest().GetStructure(origin, size)
    if err != nil {
		panic(fmt.Sprintf("SystemTestingStructureRequest: %v", err))
    }
    
    // 读取 (1,1,1) 相对位置的方块
    relPos := [3]int32{1, 1, 1}
    block, entityData, err := structure.GetBlock(relPos)
    if err != nil {
		panic(fmt.Sprintf("SystemTestingStructureRequest: %v", err))
    }
    
    pterm.Success.Printfln("Block ID: %s", block.Name)
    pterm.Success.Printfln("Block States: %s", block.States)
    pterm.Success.Printfln("Block Val: %d", block.Val)
    
    if entityData != nil {
    	pterm.Success.Printfln("Block Entity Data %v", entityData)
    }
	pterm.Success.Printfln("SystemTestingStructureRequest: PASS (Time used = %v)", time.Since(tA))
}
