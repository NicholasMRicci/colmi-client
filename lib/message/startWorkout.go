package message

const TAG_WORKOUT_SEND byte = 119
const TAG_WORKOUT_RECV byte = 120

func StartWorkout() Message {
	data := make([]byte, 14)
	data[0] = 1
	data[1] = 23
	return Message{tag: TAG_WORKOUT_SEND, data: data}
}

func DecodeWorkout(msg Message) (uint8, bool) {
	if msg.tag != TAG_WORKOUT_RECV {
		return 0, false
	}
	return msg.data[4], true
}

func PauseWorkout() Message {
	data := make([]byte, 14)
	data[0] = 2
	data[1] = 23
	return Message{tag: TAG_WORKOUT_SEND, data: data}
}

func EndWorkout() Message {
	data := make([]byte, 14)
	data[0] = 4
	data[1] = 23
	return Message{tag: TAG_WORKOUT_SEND, data: data}
}
