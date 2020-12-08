// Code generated by github.com/Skycoin/skyencoder. DO NOT EDIT.
package daemon

import (
	"errors"
	"math"

	"github.com/Skycoin/cx-chains/src/cipher/encoder"
)

// encodeSizeDisconnectMessage computes the size of an encoded object of type DisconnectMessage
func encodeSizeDisconnectMessage(obj *DisconnectMessage) uint64 {
	i0 := uint64(0)

	// obj.ReasonCode
	i0 += 2

	// obj.Reserved
	i0 += 4 + uint64(len(obj.Reserved))

	return i0
}

// encodeDisconnectMessage encodes an object of type DisconnectMessage to a buffer allocated to the exact size
// required to encode the object.
func encodeDisconnectMessage(obj *DisconnectMessage) ([]byte, error) {
	n := encodeSizeDisconnectMessage(obj)
	buf := make([]byte, n)

	if err := encodeDisconnectMessageToBuffer(buf, obj); err != nil {
		return nil, err
	}

	return buf, nil
}

// encodeDisconnectMessageToBuffer encodes an object of type DisconnectMessage to a []byte buffer.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func encodeDisconnectMessageToBuffer(buf []byte, obj *DisconnectMessage) error {
	if uint64(len(buf)) < encodeSizeDisconnectMessage(obj) {
		return encoder.ErrBufferUnderflow
	}

	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.ReasonCode
	e.Uint16(obj.ReasonCode)

	// obj.Reserved length check
	if uint64(len(obj.Reserved)) > math.MaxUint32 {
		return errors.New("obj.Reserved length exceeds math.MaxUint32")
	}

	// obj.Reserved length
	e.Uint32(uint32(len(obj.Reserved)))

	// obj.Reserved copy
	e.CopyBytes(obj.Reserved)

	return nil
}

// decodeDisconnectMessage decodes an object of type DisconnectMessage from a buffer.
// Returns the number of bytes used from the buffer to decode the object.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
func decodeDisconnectMessage(buf []byte, obj *DisconnectMessage) (uint64, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.ReasonCode
		i, err := d.Uint16()
		if err != nil {
			return 0, err
		}
		obj.ReasonCode = i
	}

	{
		// obj.Reserved

		ul, err := d.Uint32()
		if err != nil {
			return 0, err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return 0, encoder.ErrBufferUnderflow
		}

		if length != 0 {
			obj.Reserved = make([]byte, length)

			copy(obj.Reserved[:], d.Buffer[:length])
			d.Buffer = d.Buffer[length:]
		}
	}

	return uint64(len(buf) - len(d.Buffer)), nil
}

// decodeDisconnectMessageExact decodes an object of type DisconnectMessage from a buffer.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
// If the buffer is longer than required to decode the object, returns encoder.ErrRemainingBytes.
func decodeDisconnectMessageExact(buf []byte, obj *DisconnectMessage) error {
	if n, err := decodeDisconnectMessage(buf, obj); err != nil {
		return err
	} else if n != uint64(len(buf)) {
		return encoder.ErrRemainingBytes
	}

	return nil
}
