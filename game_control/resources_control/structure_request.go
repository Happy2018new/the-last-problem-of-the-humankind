package resources_control

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// StructureRequestCallback 是简单的结构请求回调维护器
type StructureRequestCallback struct {
	callback utils.SyncMap[string, func(p *packet.StructureTemplateDataResponse)]
}

// NewStructureRequestCallback 创建并返回一个新的 StructureRequestCallback
func NewStructureRequestCallback() *StructureRequestCallback {
	return new(StructureRequestCallback)
}

// SetStructureRequestCallback 设置当收到请求 ID 为 request 的结构请求的响应后，
// 应当执行的回调函数 f。其中，p 指示服务器发送的针对此结构请求的响应体
func (c *StructureRequestCallback) SetStructureRequestCallback(
	request *packet.StructureTemplateDataRequest,
	f func(p *packet.StructureTemplateDataResponse),
) {
	c.callback.Store(request.StructureName, f)
}

// DeleteStructureRequestCallback 清除请求
// ID 为 structureInfo 的命令请求的回调函数。
// 此函数应当只在结构请求超时的时候被调用
func (c *StructureRequestCallback) DeleteStructureRequestCallback(structureName string) {
	c.callback.Delete(structureName)
}

// onStructureTemplateDataResponse ..
func (c *StructureRequestCallback) onStructureTemplateDataResponse(p *packet.StructureTemplateDataResponse) {
	cb, ok := c.callback.LoadAndDelete(p.StructureName)
	if ok {
		go cb(p)
	}
}
