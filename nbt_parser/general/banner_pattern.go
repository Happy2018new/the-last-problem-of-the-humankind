package nbt_parser_general

import "github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"

// 描述旗帜的种类
const (
	BannerTypeNormal  int32 = iota // 普通旗帜
	BannerTypeOminous              // 不祥旗帜
)

// BannerBaseColorDefault 是旗帜的默认颜色 (黑色)
const BannerBaseColorDefault int32 = iota

// BannerPattern 是旗帜的单个图案
type BannerPattern struct {
	Color   int32  `mapstructure:"Color"`
	Pattern string `mapstructure:"Pattern"`
}

// Marshal ..
func (b *BannerPattern) Marshal(io protocol.IO) {
	io.Varint32(&b.Color)
	io.String(&b.Pattern)
}
