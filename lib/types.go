package lib

import (
	"tinygo.org/x/bluetooth"
)

type Ring struct {
	rx         bluetooth.DeviceCharacteristic
	tx         bluetooth.DeviceCharacteristic
	disconnect func() error
}

type Message struct {
	tag  byte
	data []byte
}

func NewMessage(tag byte, data []byte) Message {
	return Message{tag, data}
}
