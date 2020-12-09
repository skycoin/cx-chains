// Code generated by github.com/Skycoin/skyencoder. DO NOT EDIT.
package coin

import "github.com/skycoin/skycoin/src/cipher/encoder"

// encodeSizeUxHead computes the size of an encoded object of type UxHead
func encodeSizeUxHead(obj *UxHead) uint64 {
	i0 := uint64(0)

	// obj.Time
	i0 += 8

	// obj.BkSeq
	i0 += 8

	return i0
}

// encodeUxHead encodes an object of type UxHead to a buffer allocated to the exact size
// required to encode the object.
func encodeUxHead(obj *UxHead) ([]byte, error) {
	n := encodeSizeUxHead(obj)
	buf := make([]byte, n)

	if err := encodeUxHeadToBuffer(buf, obj); err != nil {
		return nil, err
	}

	return buf, nil
}

// encodeUxHeadToBuffer encodes an object of type UxHead to a []byte buffer.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func encodeUxHeadToBuffer(buf []byte, obj *UxHead) error {
	if uint64(len(buf)) < encodeSizeUxHead(obj) {
		return encoder.ErrBufferUnderflow
	}

	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.Time
	e.Uint64(obj.Time)

	// obj.BkSeq
	e.Uint64(obj.BkSeq)

	return nil
}

// decodeUxHead decodes an object of type UxHead from a buffer.
// Returns the number of bytes used from the buffer to decode the object.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
func decodeUxHead(buf []byte, obj *UxHead) (uint64, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.Time
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.Time = i
	}

	{
		// obj.BkSeq
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.BkSeq = i
	}

	return uint64(len(buf) - len(d.Buffer)), nil
}

// decodeUxHeadExact decodes an object of type UxHead from a buffer.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
// If the buffer is longer than required to decode the object, returns encoder.ErrRemainingBytes.
func decodeUxHeadExact(buf []byte, obj *UxHead) error {
	if n, err := decodeUxHead(buf, obj); err != nil {
		return err
	} else if n != uint64(len(buf)) {
		return encoder.ErrRemainingBytes
	}

	return nil
}
