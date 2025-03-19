package lib

import (
	"errors"
	"log"

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
		panic(err)
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

func (r Ring) Send(msg Message) error {
	log.Println("Sending Message")
	crc := msg.tag
	if len(msg.data) != 14 {
		return errors.New("wrong data size")
	}
	for _, piece := range msg.data {
		crc += piece
	}
	bytes := make([]byte, 0, 16)
	bytes = append(bytes, msg.tag)
	bytes = append(bytes, msg.data...)
	bytes = append(bytes, crc)
	r.tx.Write(bytes)

	return nil
}

func (r Ring) Read() Message {
	recv := make([]byte, 16)
	n, err := r.rx.Read(recv)
	if n != 16 || err != nil {
		panic(err)
	}
	calcedCrc := byte(0)
	for _, val := range recv {
		calcedCrc += val
	}
	if calcedCrc != recv[15] {
		panic("CRC Mismatch")
	}
	return Message{tag: recv[0], data: recv[1:14]}
}

func (r Ring) Disconnect() error {
	log.Println("Disconnecting from ring")
	return r.disconnect()
}
