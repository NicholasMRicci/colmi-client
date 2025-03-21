package message

import (
	"errors"
	"log"
)

type Message struct {
	tag  byte
	data []byte
}

func (msg Message) GetBytes() ([]byte, error) {
	crc := msg.tag
	if len(msg.data) != 14 {
		return nil, errors.New("wrong data size")
	}
	for _, piece := range msg.data {
		crc += piece
	}
	bytes := make([]byte, 0, 16)
	bytes = append(bytes, msg.tag)
	bytes = append(bytes, msg.data...)
	bytes = append(bytes, crc)
	return bytes, nil
}

func FromBytes(recv []byte) (Message, error) {
	calcedCrc := byte(0)
	for _, val := range recv[0:15] {
		calcedCrc += val
	}
	if calcedCrc != recv[15] {
		log.Printf("%v", recv)
		return Message{}, errors.New("bad crc")
	}
	return Message{tag: recv[0], data: recv[1:14]}, nil
}
