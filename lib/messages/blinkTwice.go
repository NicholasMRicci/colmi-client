package messages

import "github.com/NicholasMRicci/colmi-client/lib"

func BlinkTwice() lib.Message {
	return lib.NewMessage(16, make([]byte, 14))
}
