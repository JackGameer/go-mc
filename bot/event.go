package bot

import (
	"github.com/Tnze/go-mc/bot/world/entity"
	"github.com/Tnze/go-mc/chat"

	pk "github.com/Tnze/go-mc/net/packet"
)

type eventBroker struct {
	GameStart          func() error
	ChatMsg            func(msg chat.Message, pos byte) error
	Disconnect         func(reason chat.Message) error
	HealthChange       func() error
	ExperienceChange   func() error
	Die                func() error
	SoundPlay          func(name string, category int, x, y, z float64, volume, pitch float32) error
	PluginMessage      func(channel string, data []byte) error
	HeldItemChange     func(slot int) error
	WindowsItem        func(id byte, slots []entity.Slot) error
	WindowsItemChange  func(id byte, slotID int, slot entity.Slot) error
	SpawnObject        func(entityID int, UUID [16]byte, mobType int, x, y, z float64, pitch, yaw float32, data int, velocityX, velocitY, velocityZ int16) error
	SpawnEntity        func(entityID int, UUID pk.UUID, mobType int, x, y, z float64, yaw, pitch, headPitch int8, velocityX, velocitY, velocityZ int16) error
	DestroyEntities    func(entityIDs []int) error
	EntityRelativeMove func(EntityID, DeltaX, DeltaY, DeltaZ int, onGround bool) error
	// ReceivePacket will be called when new packet arrive.
	// Default handler will run only if pass == false.
	ReceivePacket func(p pk.Packet) (pass bool, err error)
}
