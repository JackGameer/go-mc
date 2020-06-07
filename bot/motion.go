package bot

import (
	"errors"
	"math"
	"strconv"

	"github.com/Tnze/go-mc/data"
	pk "github.com/Tnze/go-mc/net/packet"
)

// SwingArm swing player's arm.
// hand could be one of 0: main hand, 1: off hand.
// It's just animation.
func (c *Client) SwingArm(hand int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.AnimationServerbound,
		pk.VarInt(hand),
	))
}

// Respawn the player when it was dead.
func (c *Client) Respawn() error {
	return c.conn.WritePacket(pk.Marshal(
		data.ClientStatus,
		pk.VarInt(0),
	))
}

// UseItem use the item player handing.
// hand could be one of 0: main hand, 1: off hand
func (c *Client) UseItem(hand int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.UseItem,
		pk.VarInt(hand),
	))
}

// UseEntity used by player to right-clicks another entity.
// hand could be one of 0: main hand, 1: off hand.
// A Notchian server only accepts this packet if
// the entity being attacked/used is visible without obstruction
// and within a 4-unit radius of the player's position.
func (c *Client) UseEntity(entityID int32, hand int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.UseEntity,
		pk.VarInt(entityID),
		pk.VarInt(0),
		pk.VarInt(hand),
	))
}

// AttackEntity used by player to left-clicks another entity.
// The attack version of UseEntity. Has the same limit.
func (c *Client) AttackEntity(entityID int32, hand int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.UseEntity,
		pk.VarInt(entityID),
		pk.VarInt(1),
	))
}

// UseEntityAt is a variety of UseEntity with target location
func (c *Client) UseEntityAt(entityID int32, x, y, z float32, hand int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.UseEntity,
		pk.VarInt(entityID),
		pk.VarInt(2),
		pk.Float(x), pk.Float(y), pk.Float(z),
		pk.VarInt(hand),
	))
}

// Chat send chat as chat message or command at textbox.
func (c *Client) Chat(msg string) error {
	if len(msg) > 256 {
		return errors.New("message too long")
	}

	return c.conn.WritePacket(pk.Marshal(
		data.ChatMessageServerbound,
		pk.String(msg),
	))
}

// PluginMessage is used by mods and plugins to send their data.
func (c *Client) PluginMessage(channal string, msg []byte) error {
	return c.conn.WritePacket(pk.Marshal(
		data.PluginMessageServerbound,
		pk.Identifier(channal),
		pluginMessageData(msg),
	))
}

// UseBlock is used to place or use a block.
// hand is the hand from which the block is placed; 0: main hand, 1: off hand.
// face is the face on which the block is placed.
//
// Cursor position is the position of the crosshair on the block:
// cursorX, from 0 to 1 increasing from west to east;
// cursorY, from 0 to 1 increasing from bottom to top;
// cursorZ, from 0 to 1 increasing from north to south.
//
// insideBlock is true when the player's head is inside of a block's collision.
func (c *Client) UseBlock(hand, locX, locY, locZ, face int, cursorX, cursorY, cursorZ float32, insideBlock bool) error {
	return c.conn.WritePacket(pk.Marshal(
		data.PlayerBlockPlacement,
		pk.VarInt(hand),
		pk.Position{X: locX, Y: locY, Z: locZ},
		pk.VarInt(face),
		pk.Float(cursorX), pk.Float(cursorY), pk.Float(cursorZ),
		pk.Boolean(insideBlock),
	))
}

// SelectItem used to change the slot selection in hotbar.
// slot should from 0 to 8
func (c *Client) SelectItem(slot int) error {
	if slot < 0 || slot > 8 {
		return errors.New("invalid slot: " + strconv.Itoa(slot))
	}

	return c.conn.WritePacket(pk.Marshal(
		data.HeldItemChangeServerbound,
		pk.Short(slot),
	))
}

// PickItem used to swap out an empty space on the hotbar with the item in the given inventory slot.
// The Notchain client uses this for pick block functionality (middle click) to retrieve items from the inventory.
//
// The server will first search the player's hotbar for an empty slot,
// starting from the current slot and looping around to the slot before it.
// If there are no empty slots, it will start a second search from the
// current slot and find the first slot that does not contain an enchanted item.
// If there still are no slots that meet that criteria, then the server will
// use the currently selected slot. After finding the appropriate slot,
// the server swaps the items and then change player's selected slot (cause the HeldItemChange event).
func (c *Client) PickItem(slot int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.PickItem,
		pk.VarInt(slot),
	))
}

func (c *Client) playerAction(status, locX, locY, locZ, face int) error {
	return c.conn.WritePacket(pk.Marshal(
		data.PlayerDigging,
		pk.VarInt(status),
		pk.Position{X: locX, Y: locY, Z: locZ},
		pk.Byte(face),
	))
}

// Dig used to start, end or cancel a digging
// status is 0 for start digging, 1 for cancel and 2 if client think it done.
// To digging a block without cancel, use status 0 and 2 once each.
func (c *Client) Dig(status, locX, locY, locZ, face int) error {
	return c.playerAction(status, locX, locY, locZ, face)
}

// DropItemStack drop the entire selected stack
func (c *Client) DropItemStack() error {
	return c.playerAction(3, 0, 0, 0, 0)
}

// DropItem drop one item in selected stack
func (c *Client) DropItem() error {
	return c.playerAction(4, 0, 0, 0, 0)
}

// UseItemEnd used to finish UseItem, like eating food, pulling back bows.
func (c *Client) UseItemEnd() error {
	return c.playerAction(5, 0, 0, 0, 0)
}

// SwapItem used to swap the items in hands.
func (c *Client) SwapItem() error {
	return c.playerAction(6, 0, 0, 0, 0)
}

// Disconnect disconnect the server.
// Server will close the connection.
func (c *Client) Disconnect() error {
	return c.conn.Close()
}

// SendPacket send the packet to server.
func (c *Client) SendPacket(packet pk.Packet) error {
	return c.conn.WritePacket(packet)
}

// SetPosition method move your character around.
// Server will ignore this if changes too much.
func (c *Client) SetPosition(x, y, z float64, onGround bool) {
	c.Player.X, c.Player.Y, c.Player.Z = x, y, z
	c.Player.OnGround = onGround
	sendPlayerPositionPacket(c)
}

func (c *Client) SetPositionAndLook(x, y, z float64, onGround bool, yaw, pitch float32) {
	c.Player.X, c.Player.Y, c.Player.Z = x, y, z
	c.Player.OnGround = onGround
	c.Player.Yaw, c.Player.Pitch = yaw, pitch
	sendPlayerPositionAndLook(c)
}
func sendPlayerPositionAndLook(c *Client) {
	c.conn.WritePacket(pk.Marshal(
		data.PlayerPositionAndLookServerbound,
		pk.Double(c.X),
		pk.Double(c.Y),
		pk.Double(c.Z),
		pk.Float(c.Yaw),
		pk.Float(c.Pitch),
		pk.Boolean(c.OnGround),
	))
}

func sendPlayerPositionPacket(c *Client) error {
	return c.conn.WritePacket(pk.Marshal(
		data.PlayerPosition,
		pk.Double(c.Player.X),
		pk.Double(c.Player.Y),
		pk.Double(c.Player.Z),
		pk.Boolean(c.Player.OnGround),
	))
}

// LookAt method turn player's hand and make it look at a point.
func (c *Client) LookAt(x, y, z float64) {
	x0, y0, z0 := c.Player.X, c.Player.Y, c.Player.Z
	x, y, z = x-x0, y-y0, z-z0

	r := math.Sqrt(x*x + y*y + z*z)
	yaw := -math.Atan2(x, z) / math.Pi * 180
	for yaw < 0 {
		yaw = 360 + yaw
	}
	pitch := -math.Asin(y/r) / math.Pi * 180

	c.LookYawPitch(float32(yaw), float32(pitch))
}

// LookYawPitch set player's hand to the direct by yaw and pitch.
// yaw can be [0, 360) and pitch can be (-180, 180).
// if |pitch|>90 the player's hand will be very strange.
func (c *Client) LookYawPitch(yaw, pitch float32) {
	c.Player.Yaw, c.Player.Pitch = yaw, pitch
	sendPlayerLookPacket(c)
}

func sendPlayerLookPacket(c *Client) error {
	return c.conn.WritePacket(pk.Marshal(
		data.PlayerLook,
		pk.Float(c.Player.Yaw),
		pk.Float(c.Player.Pitch),
		pk.Boolean(c.Player.OnGround),
	))
}
