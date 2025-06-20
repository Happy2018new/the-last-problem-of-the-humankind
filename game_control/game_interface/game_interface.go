package game_interface

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// ResourcesWrapper 是基于资源中心包装的机器人资源
type ResourcesWrapper struct {
	resources_control.BotInfo
	*resources_control.Resources
}

// GameInterface 实现了机器人与租赁服的高级交互，
// 例如基本的命令收发或高级的容器操作
type GameInterface struct {
	wrapper               *ResourcesWrapper
	commands              *Commands
	structureBackup       *StructureBackup
	querytarget           *Querytarget
	setblock              *SetBlock
	replaceitem           *Replaceitem
	botClick              *BotClick
	itemStackOperation    *ItemStackOperation
	containerOpenAndClose *ContainerOpenAndClose
	itemCopy              *ItemCopy
	itemTransition        *ItemTransition
	structureRequest      *StructureRequest
}

// NewResourcesWrapper 基于 resources 创建一个新的游戏交互器
func NewResourcesWrapper(resources *resources_control.Resources) *ResourcesWrapper {
	return &ResourcesWrapper{
		BotInfo:   resources.BotInfo(),
		Resources: resources,
	}
}

// NewGameInterface 基于 resources 创建一个新的游戏交互器
func NewGameInterface(resources *resources_control.Resources) *GameInterface {
	result := new(GameInterface)

	result.wrapper = NewResourcesWrapper(resources)
	result.commands = NewCommands(result.wrapper)
	result.structureBackup = NewStructureBackup(result.commands)
	result.querytarget = NewQuerytarget(result.commands)
	result.setblock = NewSetBlock(result.commands)
	result.replaceitem = NewReplaceitem(result.commands)
	result.botClick = NewBotClick(result.wrapper, result.commands, result.setblock)
	result.itemStackOperation = NewItemStackOperation(result.wrapper)
	result.containerOpenAndClose = NewContainerOpenAndClose(result.wrapper, result.commands, result.botClick)
	result.itemCopy = NewItemCopy(result.containerOpenAndClose, result.commands, result.itemStackOperation, result.structureBackup)
	result.itemTransition = NewItemTransition(result.wrapper, result.itemStackOperation)
	result.structureRequest = NewStructureRequest(result.wrapper)

	return result
}

// GetBotInfo 返回机器人的基本信息
func (g *GameInterface) GetBotInfo() resources_control.BotInfo {
	return g.wrapper.BotInfo
}

// PacketListener 返回一个可撤销的数据包监听实现
func (g *GameInterface) PacketListener() *resources_control.PacketListener {
	return g.wrapper.PacketListener()
}

// Resources 返回底层的资源管理器
func (g *GameInterface) Resources() *resources_control.Resources {
	return g.wrapper.Resources
}

// Commands 返回机器人在 MC 命令在收发上的相关实现
func (g *GameInterface) Commands() *Commands {
	return g.commands
}

// StructureBackup 返回机器人在结构备份和恢复上的相关实现
func (g *GameInterface) StructureBackup() *StructureBackup {
	return g.structureBackup
}

// Querytarget 返回机器人在 querytarget 命令上的相关实现
func (g *GameInterface) Querytarget() *Querytarget {
	return g.querytarget
}

// SetBlock 返回机器人在方块放置 (MC 命令的方式) 上的相关实现
func (g *GameInterface) SetBlock() *SetBlock {
	return g.setblock
}

// Replaceitem 返回机器人在 Replaceitem 命令上的简单包装
func (g *GameInterface) Replaceitem() *Replaceitem {
	return g.replaceitem
}

// BotClick 返回机器人在点击操作上的相关实现。
//
// 由于点击操作与机器人手持物品强相关，
// BotClick 也集成了切换手持物品的实现。
//
// 另外，考虑到 Pick Block 操作的语义也与
// 点击方块 有关，因此其也被集成在此，尽管
// 它使用了完全不同的数据包
func (g *GameInterface) BotClick() *BotClick {
	return g.botClick
}

// ItemStackOperation 返回机器人在物品堆栈操作请求上的相关实现
func (g *GameInterface) ItemStackOperation() *ItemStackOperation {
	return g.itemStackOperation
}

// ContainerOpenAndClose 返回机器人在容器打开和关闭上的相关实现
func (g *GameInterface) ContainerOpenAndClose() *ContainerOpenAndClose {
	return g.containerOpenAndClose
}

// ItemCopy 返回机器人在物品复制上的相关实现
func (g *GameInterface) ItemCopy() *ItemCopy {
	return g.itemCopy
}

// ItemTransition 返回机器人在物品状态转移上的实现
func (g *GameInterface) ItemTransition() *ItemTransition {
	return g.itemTransition
}

// StructureRequest 返回机器人在结构请求上的相关实现
func (g *GameInterface) StructureRequest() *StructureRequest {
	return g.structureRequest
}
