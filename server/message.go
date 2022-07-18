package server

import "encoding/binary"

// ref: https://www.postgresql.org/docs/14/protocol-message-formats.html

type TransactionStatus = byte

const (
	TransactionIdle   TransactionStatus = 0x49 // I
	Transaction       TransactionStatus = 0x54 // T
	TransactionFailed TransactionStatus = 0x45 // E

	nullEnd = 0x00
)

func makeMsg(tag byte, body []byte) []byte {
	length := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(length, uint32(len(body)+payloadBytesLength))

	msg := make([]byte, 0)
	msg = append(msg, tag)
	msg = append(msg, length...)
	msg = append(msg, body...)

	return msg
}

func makeReadyForQueryMsg(status TransactionStatus) []byte {
	return makeMsg('Z', []byte{status})
}

func makeCommandCompleteMsg(s string) []byte {
	body := make([]byte, 0)
	body = append(body, []byte(s)...)
	body = append(body, nullEnd)

	return makeMsg('C', body)
}

func makeParameterStatusMsg(param, value string) []byte {
	body := make([]byte, 0)
	body = append(body, []byte(param)...)
	body = append(body, nullEnd)
	body = append(body, []byte(value)...)
	body = append(body, nullEnd)

	return makeMsg('S', body)
}

func makeErrorMsg(err error) []byte {
	const errMsgEnd = 0x00
	errMsg := err.Error()

	body := make([]byte, 0)

	body = append(body, 'S') // Severity
	body = append(body, []byte("ERROR")...)
	body = append(body, nullEnd)

	body = append(body, 'V') // Severity
	body = append(body, []byte("ERROR")...)
	body = append(body, nullEnd)

	body = append(body, 'M') // message
	body = append(body, []byte(errMsg)...)
	body = append(body, nullEnd)

	body = append(body, errMsgEnd)

	return makeMsg('E', body)
}
