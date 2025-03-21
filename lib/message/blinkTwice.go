package message

func BlinkTwice() Message {
	return Message{tag: 16, data: make([]byte, 14)}
}
