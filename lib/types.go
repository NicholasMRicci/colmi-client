package lib

import (
	"github.com/NicholasMRicci/colmi-client/lib/message"
	"tinygo.org/x/bluetooth"
)

type Ring struct {
	rx          bluetooth.DeviceCharacteristic
	tx          bluetooth.DeviceCharacteristic
	disconnect  func() error
	messageChan chan message.Message
}
