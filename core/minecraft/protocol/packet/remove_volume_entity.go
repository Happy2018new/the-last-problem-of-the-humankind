package packet

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
)

// RemoveVolumeEntity indicates a volume entity to be removed from server to client.
type RemoveVolumeEntity struct {
	// EntityRuntimeID ..
	EntityRuntimeID uint64
	// Dimension ...
	Dimension int32
}

// ID ...
func (*RemoveVolumeEntity) ID() uint32 {
	return IDRemoveVolumeEntity
}

func (pk *RemoveVolumeEntity) Marshal(io protocol.IO) {
	io.Uint64(&pk.EntityRuntimeID)
	io.Varint32(&pk.Dimension)
}
