package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ErrorInvalidID      = "Invalid message ID"
	ErrorInvalidPayload = "Invalid message payload"
)

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type messageID uint8

type Message struct {
	ID      messageID
	Payload []byte
}

func (msg *Message) Serialize() []byte {
	buf := bytes.Buffer{}
	length := uint32(1 + len(msg.Payload))

	binary.Write(&buf, binary.BigEndian, length)

	buf.WriteByte(byte(msg.ID))
	buf.Write(msg.Payload)
	return buf.Bytes()
}

func BuildHandshake(infoHash, peerID [20]byte) []byte {
	pstr := "BitTorrent protocol"
	buf := bytes.Buffer{}
	buf.WriteByte(byte(len(pstr))) // Total length of the message
	buf.WriteString(pstr)          // Protocol identifier
	buf.Write(make([]byte, 8))     // Reserved 8 bytes
	buf.Write(infoHash[:])         // InfoHash
	buf.Write(peerID[:])           //Peer ID
	return buf.Bytes()
}

func BuildHaveMessage(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}

func BuildRequest(index, begin, length int) *Message {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, uint32(index))
	binary.Write(&buf, binary.BigEndian, uint32(begin))
	binary.Write(&buf, binary.BigEndian, uint32(length))
	return &Message{ID: MsgRequest, Payload: buf.Bytes()}
}

func ParseHandshake(reader io.Reader) (infoHash, peerID [20]byte, err error) {
	lengthBuf := make([]byte, 1)
	_, err = io.ReadFull(reader, lengthBuf)
	if err != nil {
		return
	}
	length := int(lengthBuf[0])
	buf := make([]byte, 48+length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return
	}
	copy(infoHash[:], buf[length+8:length+28])
	copy(peerID[:], buf[length+28:])
	return
}

func (msg *Message) ParsePiece(index int, buf []byte) (blockLen int, err error) {
	if msg.ID != MsgPiece {
		err = fmt.Errorf(ErrorInvalidID)
		return
	}
	if len(msg.Payload) < 8 {
		err = fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
		return
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		err = fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
		return
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		err = fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
		return
	}
	block := msg.Payload[8:]
	if begin+len(block) > len(buf) {
		err = fmt.Errorf("Data too long [%d] for offset %d with length %d", len(block), begin, len(buf))
		return
	}
	blockLen = copy(buf[begin:], block)
	return
}

func (msg *Message) ParseHave() (int, error) {
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf(ErrorInvalidPayload)
	}
	if msg.ID != MsgHave {
		return 0, fmt.Errorf(ErrorInvalidID)
	}
	return int(binary.BigEndian.Uint32(msg.Payload)), nil
}

func Read(reader io.Reader) (*Message, error) {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(reader, lenBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	if length == 0 {
		return nil, nil
	}
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(reader, messageBuf)
	if err != nil {
		return nil, err
	}
	return &Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}, nil
}
