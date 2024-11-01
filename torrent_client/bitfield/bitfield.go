package bitfield

import (
	"net"
)

type BitField []byte

func RecvBitfield(conn net.Conn, infoHash [20]byte) error {

	return nil
}

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8
	bitIndex := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}
	return bf[byteIndex]&(1<<(7-bitIndex)) != 0
}

func (bf BitField) SetPiece(index int) {
	byteIndex := index / 8
	bitIndex := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7-bitIndex)
}
