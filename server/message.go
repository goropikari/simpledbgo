package server

import "encoding/binary"

// ref: https://www.postgresql.org/docs/14/protocol-message-formats.html

func makeReadyForQueryMsg() []byte {
	return []byte{0x5a, 0x00, 0x00, 0x00, 0x05, 0x49}
}

func makeCommandCompleteMsg(s string) []byte {
	body := make([]byte, 0)
	body = append(body, []byte(s)...)
	body = append(body, 0x00)
	l := len(body)
	lb := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lb, uint32(l+payloadBytesLength))
	payload := make([]byte, 0)
	payload = append(payload, 0x43) // 0x43 -> C: CommandComplete
	payload = append(payload, lb...)
	payload = append(payload, body...)

	return payload
}

func makeParameterStatusMsg(param, value string) []byte {
	msg := make([]byte, 0)
	length := make([]byte, payloadBytesLength)
	body := make([]byte, 0)
	body = append(body, []byte(param)...)
	body = append(body, 0x00)
	body = append(body, []byte(value)...)
	body = append(body, 0x00)
	binary.BigEndian.PutUint32(length, uint32(len(body)+payloadBytesLength))

	msg = append(msg, 'S')
	msg = append(msg, length...)
	msg = append(msg, body...)

	return msg
}
