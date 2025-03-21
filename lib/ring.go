package lib

import (
	"errors"
	"log"

	"github.com/NicholasMRicci/colmi-client/lib/message"
	"tinygo.org/x/bluetooth"
)

var serviceUUID, _ = bluetooth.ParseUUID("6E40FFF0-B5A3-F393-E0A9-E50E24DCCA9E")
var txUUID, _ = bluetooth.ParseUUID("6E400002-B5A3-F393-E0A9-E50E24DCCA9E")
var rxUUID, _ = bluetooth.ParseUUID("6E400003-B5A3-F393-E0A9-E50E24DCCA9E")

func AquireRing(BLE *bluetooth.Adapter, name string) (Ring, error) {
	var found bluetooth.ScanResult

	log.Println("Scanning for ring")
	err := BLE.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == "COLMI R10_AB02" {
			defer BLE.StopScan()
			found = device
			log.Println("Found Ring")
		}
	})
	if err != nil {
		return Ring{}, err
	}

	conn := Must1(BLE.Connect(found.Address, bluetooth.ConnectionParams{}))

	services := Must1(conn.DiscoverServices([]bluetooth.UUID{serviceUUID}))
	if len(services) != 1 {
		log.Fatalln("Ring wrong services")
	}
	characteristics := Must1(services[0].DiscoverCharacteristics([]bluetooth.UUID{txUUID, rxUUID}))
	if len(characteristics) != 2 {
		log.Fatalln("Ring is not advertizing characteristics")
	}
	log.Println("Ring Good to go")

	return Ring{tx: characteristics[0], rx: characteristics[1], disconnect: conn.Disconnect}, nil
}

func (r Ring) Send(msg message.Message) error {
	log.Println("Sending Message")
	bytes, err := msg.GetBytes()
	if err != nil {
		return err
	}
	r.tx.WriteWithoutResponse(bytes)

	return nil
}

func (r Ring) BeginReads(data chan message.Message) error {
	if r.messageChan != nil {
		return errors.New("Handler already registered")
	}
	return r.rx.EnableNotifications(func(buf []byte) {
		msg, err := message.FromBytes(buf)
		if err != nil && len(buf) != 0 {
			log.Printf("Bad message idk what to do %v", msg)
		}
		if len(buf) != 16 {
			log.Printf("Bad message idk what to do %v", msg)
		}
		data <- msg
	})
}

func (r Ring) StopReads() error {
	if r.messageChan == nil {
		return errors.New("No reads right now")
	}
	return r.rx.EnableNotifications(nil)
}

func (r Ring) Disconnect() error {
	log.Println("Disconnecting from ring")
	return r.disconnect()
}
