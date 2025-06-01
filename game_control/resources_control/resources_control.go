package resources_control

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

// BotBasicInfo 记载机器人的基本信息
type BotBasicInfo struct {
	BotName         string // 机器人名称
	XUID            string // 机器人 XUID
	EntityUniqueID  int64  // 机器人唯一 ID
	EntityRuntimeID uint64 // 机器人运行时 ID
}

type Resources struct {
	// client 是连接到租赁服的基本客户端
	client *client.Client
	// commands 存放所有命令请求的回调函数
	commands *CommandRequestCallback
	// inventory 持有机器人已经拥有或打开的库存
	inventory *Inventories
	// itemStack 管理物品堆栈操作请求
	itemStack *ItemStackOperationManager
	// container 维护机器人的容器资源，
	// 处理其占用和释放，以及一些持久化数据
	container *ContainerManager
	// listener 是一个可撤销的简单数据包监听器实现
	listener *PacketListener
}

// NewResourcesControl 基于 client 创建一个新的资源中心。
// 它应当在机器人连接到租赁服后立即被调用，且最多调用一次。
//
// 需要注意的是，client.Conn().ReadPacket 不应继续被使用，
// 否则可能会出现未知的竞态条件问题，因为资源管理器本身也会
// 不断的读取数据包并依此更新其自身的资源数据
func NewResourcesControl(client *client.Client) *Resources {
	resourcesControl := &Resources{
		client:    client,
		commands:  NewCommandRequestCallback(),
		inventory: NewInventories(),
		itemStack: NewItemStackOperationManager(),
		container: NewContainerManager(),
		listener:  NewPacketListener(),
	}
	go resourcesControl.listenPacket()
	return resourcesControl
}

// listenPacket ..
func (r *Resources) listenPacket() {
	for {
		pk := <-r.client.CachedPacket()
		if pk == nil {
			break
		}
		r.handlePacket(pk)
	}
	for {
		pk, err := r.client.Conn().ReadPacket()
		if err != nil {
			return
		}
		r.handlePacket(pk)
	}
}

// ClientBasicInfo 返回机器人的基本信息
func (r *Resources) ClientBasicInfo() BotBasicInfo {
	return BotBasicInfo{
		BotName:         r.client.Conn().IdentityData().DisplayName,
		XUID:            r.client.Conn().IdentityData().XUID,
		EntityUniqueID:  r.client.Conn().GameData().EntityUniqueID,
		EntityRuntimeID: r.client.Conn().GameData().EntityRuntimeID,
	}
}

// WritePacket 返回向租赁服发送数据包的函数
func (r *Resources) WritePacket() func(p packet.Packet) error {
	return r.client.Conn().WritePacket
}

// Commands 返回命令请求的相关资源
func (r *Resources) Commands() *CommandRequestCallback {
	return r.commands
}

// Inventories 返回库存的相关资源
func (r *Resources) Inventories() *Inventories {
	return r.inventory
}

// ItemStackOperation 返回物品堆栈操作请求的相关资源
func (r *Resources) ItemStackOperation() *ItemStackOperationManager {
	return r.itemStack
}

// Container 返回容器的相关资源
func (r *Resources) Container() *ContainerManager {
	return r.container
}

// PacketListener 返回数据包监听的有关实现
func (r *Resources) PacketListener() *PacketListener {
	return r.listener
}
