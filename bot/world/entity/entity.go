package entity

import (
	"strings"

	"github.com/Tnze/go-mc/data"
	en_us "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/Tnze/go-mc/nbt"
	pk "github.com/Tnze/go-mc/net/packet"
)

//Entity is the entity of minecraft
type Entity struct {
	EntityID int //实体ID
	Type     int
	X, Y, Z  float64
}

// The Slot data structure is how Minecraft represents an item and its associated data in the Minecraft Protocol
type Slot struct {
	Present bool
	ItemID  int32
	Count   int8
	NBT     interface{}
}

//Decode implement packet.FieldDecoder interface
func (s *Slot) Decode(r pk.DecodeReader) error {
	if err := (*pk.Boolean)(&s.Present).Decode(r); err != nil {
		return err
	}
	if s.Present {
		if err := (*pk.VarInt)(&s.ItemID).Decode(r); err != nil {
			return err
		}
		if err := (*pk.Byte)(&s.Count).Decode(r); err != nil {
			return err
		}
		if err := nbt.NewDecoder(r).Decode(&s.NBT); err != nil {
			return err
		}
	}
	return nil
}

func (s Slot) String() string {
	return data.ItemNameByID[s.ItemID]
}

//Name returns translated name
func (s Slot) Name() string {
	key := strings.Split(s.String(), ":")
	id := key[0] + "." + key[1]
	translation, ok := en_us.Map["item."+id]
	if !ok {
		translation, ok = en_us.Map["block."+id]
		if !ok {
			return id
		}
	}
	return translation
}

func (e Entity) String() string {
	return data.EntityNameByID[e.Type]
}
