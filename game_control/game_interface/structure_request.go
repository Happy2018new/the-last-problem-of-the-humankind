package game_interface

import (
	"fmt"
	"strconv"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

// StructureRequest 是基于 ResourcesWrapper 实现的 MC 结构请求器
type StructureRequest struct {
	api *ResourcesWrapper
}

// StructureTemplate 表示解析后的结构数据
type StructureTemplate struct {
	Origin         [3]int32                      // 结构的世界原点坐标 (x, y, z)
	Size           [3]int32                      // 结构尺寸 (sizex, sizey, sizez)
	BlockIndices   []int32                        // 方块索引矩阵 (3D->1D 映射)
	Palette        []Block                   // 方块调色板
	BlockEntities  map[string]map[string]interface{}  // 方块实体数据 (相对坐标 -> NBT)
}

// Block 表示方块
type Block struct {
	Name    string                 // 方块ID (e.g., "minecraft:stone")
	States  map[string]interface{} // 方块状态字典
	Val     int32                  // 方块特殊值
	Version int32                  // 方块版本
}

// NewStructureRequest 基于 api 创建并返回一个新的 ResourcesWrapper
func NewStructureRequest(api *ResourcesWrapper) *StructureRequest {
	return &StructureRequest{api: api}
}

// SendStructureTemplateRequestWithResp 是用于
// 请求 request 代表的结构请求并获取与之对应的响应体
func (s *StructureRequest) SendStructureTemplateRequestWithResp(request *packet.StructureTemplateDataRequest) (
	resp *packet.StructureTemplateDataResponse,
	err error,
) {
	api := s.api
	channel := make(chan struct{})

	api.Resources.StructureRequest().SetStructureRequestCallback(
		request,
		func(p *packet.StructureTemplateDataResponse) {
			resp = p
			close(channel)
		},
	)
	err = api.WritePacket(request)
	if err != nil {
		return nil, fmt.Errorf("SendStructureTemplateRequestWithResp: %v", err)
	}

	<-channel
	return resp, nil
}

// GetStructure 获取坐标为 origin 处
// 且长度为 size 结构数据
//
// 未加载的区域也可获取结构
// size 的长宽高应小于 64
func (s *StructureRequest) GetStructure(origin [3]int32, size [3]int32) (*StructureTemplate, error) {
	api := s.api
	// 获取机器人的EntityUniqueID
	botUniqueID := api.BotInfo.EntityUniqueID
	
	// 创建结构请求
	request := &packet.StructureTemplateDataRequest{
		Position:      origin,
		StructureName: "mystructure:StructureExport",
		Settings: protocol.StructureSettings{
			PaletteName: "default",
			IgnoreEntities: true,
			IgnoreBlocks: false,
			Size: size,
			Offset: [3]int32{0, 0, 0},
			LastEditingPlayerUniqueID: botUniqueID,
			Rotation: 0,
			Mirror: 0,
			Integrity: 100,
			Seed: 0,
			AllowNonTickingChunks: false,
		},
		RequestType: packet.StructureTemplateRequestExportFromSave,
	}

	// 发送请求并获取响应
	resp, err := s.SendStructureTemplateRequestWithResp(request)
	if err != nil {
		return nil, fmt.Errorf("获取结构数据失败: %v", err)
	}

	// 解析结构响应
	return parseStructureResponse(origin, resp)
}

// parseStructureResponse 解析结构响应数据
func parseStructureResponse(origin [3]int32, resp *packet.StructureTemplateDataResponse) (*StructureTemplate, error) {
	// 从响应中提取NBT数据
	nbtData := resp.StructureTemplate

	// 安全地解析结构尺寸
	var size [3]int32
	switch v := nbtData["size"].(type) {
	case []int32:
		if len(v) >= 3 {
			size = [3]int32{v[0], v[1], v[2]}
		} else {
			return nil, fmt.Errorf("invalid size data length")
		}
	case []int16: // 处理可能的int16类型
		if len(v) >= 3 {
			size = [3]int32{int32(v[0]), int32(v[1]), int32(v[2])}
		} else {
			return nil, fmt.Errorf("invalid size data length")
		}
	default:
		return nil, fmt.Errorf("unexpected size type: %T", v)
	}

	// 解析方块索引
	blockIndices := make([]int32, 0)
	if structure, ok := nbtData["structure"].(map[string]interface{}); ok {
		// structure -> block_indices
		if indicesList, ok := structure["block_indices"].([]interface{}); ok && len(indicesList) >= 1 {
			// 取第一个列表（非-1索引）
			if layer0, ok := indicesList[0].([]int32); ok {
				// 类型为 []int32
				blockIndices = append(blockIndices, layer0...)
			} else if layer0, ok := indicesList[0].([]interface{}); ok {
				// 类型为 []interface{}
				for _, idx := range layer0 {
					switch v := idx.(type) {
					case int64:
						blockIndices = append(blockIndices, int32(v))
					case int32:
						blockIndices = append(blockIndices, v)
					case int:
						blockIndices = append(blockIndices, int32(v))
					case int16:
						blockIndices = append(blockIndices, int32(v))
					}
				}
			} else if layer0, ok := indicesList[0].([]int64); ok {
				// 类型为 []int64
				for _, v := range layer0 {
					blockIndices = append(blockIndices, int32(v))
				}
			} else if layer0, ok := indicesList[0].([]int); ok {
				// 类型为 []int
				for _, v := range layer0 {
					blockIndices = append(blockIndices, int32(v))
				}
			}
		}
	}

	// 调色板解析
	palette := make([]Block, 0)
	if structure, ok := nbtData["structure"].(map[string]interface{}); ok {
		if palettes, ok := structure["palette"].(map[string]interface{}); ok {
			if defaultPalette, ok := palettes["default"].(map[string]interface{}); ok {
				if blockPalette, ok := defaultPalette["block_palette"].([]interface{}); ok {
					for _, entry := range blockPalette {
						block, _ := entry.(map[string]interface{})
						states := make(map[string]interface{})
						
						// 处理可能缺失的 states 字段
						if statesData, ok := block["states"].(map[string]interface{}); ok {
							states = statesData
						}
						
						val := int32(0)
						if v, ok := block["val"].(int16); ok {
							val = int32(v)
						} else if v, ok := block["val"].(int32); ok {
							val = v
						} else if v, ok := block["val"].(int); ok {
							val = int32(v)
						}
						
						version := int32(0)
						if v, ok := block["version"].(int32); ok {
							version = v
						} else if v, ok := block["version"].(int64); ok {
							version = int32(v)
						}
						
						palette = append(palette, Block{
							Name:    block["name"].(string),
							States:  states,
							Val:     val,
							Version: version,
						})
					}
				}
			}
		}
	}

	// 解析方块实体数据
	blockEntities := make(map[string]map[string]interface{})
	if structure, ok := nbtData["structure"].(map[string]interface{}); ok {
		if palettes, ok := structure["palette"].(map[string]interface{}); ok {
			if defaultPalette, ok := palettes["default"].(map[string]interface{}); ok {
				if blockPosData, ok := defaultPalette["block_position_data"].(map[string]interface{}); ok {
					for posKey, data := range blockPosData {
						if dataMap, ok := data.(map[string]interface{}); ok {
							if entityData, ok := dataMap["block_entity_data"].(map[string]interface{}); ok {
								blockEntities[posKey] = entityData
							}
						}
					}
				}
			}
		}
	}

	return &StructureTemplate{
		Origin:        origin,
		Size:          size,
		BlockIndices:  blockIndices,
		Palette:       palette,
		BlockEntities: blockEntities,
	}, nil
}

// GetBlock 获取结构中指定相对位置的方块数据
// 返回：方块、方块实体数据、错误
func (st *StructureTemplate) GetBlock(relPos [3]int32) (Block, map[string]interface{}, error) {
	// 检查位置是否有效
	if relPos[0] < 0 || relPos[0] >= st.Size[0] ||
		relPos[1] < 0 || relPos[1] >= st.Size[1] ||
		relPos[2] < 0 || relPos[2] >= st.Size[2] {
		return Block{}, nil, fmt.Errorf("GetBlock: Location out of structural range")
	}

	// 计算一维索引
	index := relPos[0] * st.Size[1] * st.Size[2] + relPos[1] * st.Size[2] + relPos[2]
	
	// 检查索引是否有效
	if int(index) >= len(st.BlockIndices) {
		return Block{}, nil, fmt.Errorf("GetBlock: Index out of range (%d)", len(st.BlockIndices))
	}
	
	// 获取调色板索引
	paletteIndex := st.BlockIndices[index]
	if int(paletteIndex) >= len(st.Palette) {
		return Block{}, nil, fmt.Errorf("GetBlock: Invalid palette index")
	}
	
	// 获取方块状态
	block := st.Palette[paletteIndex]
	// 获取方块实体数据
	entityData, _ := st.BlockEntities[strconv.Itoa(int(index))]

	return block, entityData, nil
}