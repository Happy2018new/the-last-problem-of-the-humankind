package main

import (
	"sync"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
)

func SystemTestingUUIDSafeString() {
	tA := time.Now()

	// Chat
	for range 100 {
		ud := uuid.New()
		chatContent := utils.MakeUUIDSafeString(ud)
		channel := make(chan struct{})

		equalUUID, _ := utils.FromUUIDSafeString(chatContent)
		if equalUUID != ud {
			panic("SystemTestingUUIDSafeString: UUID Safe String not equal")
		}

		doOnce := new(sync.Once)
		uniqueID := api.PacketListener().ListenPacket(
			[]uint32{packet.IDText},
			func(p packet.Packet) {
				if p.(*packet.Text).Message == chatContent {
					doOnce.Do(func() { close(channel) })
				}
			},
		)
		api.Commands().SendChat(chatContent)

		timer := time.NewTimer(time.Second * 5)
		defer timer.Stop()
		select {
		case <-timer.C:
			panic("SystemTestingUUIDSafeString: Time out")
		case <-channel:
			api.PacketListener().DestroyListener(uniqueID)
		}
	}

	pterm.Success.Printfln("SystemTestingUUIDSafeString: PASS (Time used = %v)", time.Since(tA))
}
