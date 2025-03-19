package main

import (
	"github.com/NicholasMRicci/colmi-client/lib"
	"github.com/NicholasMRicci/colmi-client/lib/messages"
	"tinygo.org/x/bluetooth"
)

func main() {
	lib.Must(bluetooth.DefaultAdapter.Enable())

	// Start scanning.
	ring, err := lib.AquireRing(bluetooth.DefaultAdapter, "")
	lib.Must(err)
	defer func() { lib.Must(ring.Disconnect()) }()

	lib.Must(ring.Send(messages.BlinkTwice()))

}
