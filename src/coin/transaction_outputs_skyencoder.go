// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.
package coin

import (
	"errors"
	"math"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

// encodeSizeTransactionOutputs computes the size of an encoded object of type transactionOutputs
func encodeSizeTransactionOutputs(obj *transactionOutputs) uint64 {
	i0 := uint64(0)

	// obj.Out
	i0 += 4
	{
		i1 := uint64(0)

		// x.Address.Version
		i1++

		// x.Address.Key
		i1 += 20

		// x.Coins
		i1 += 8

		// x.Hours
		i1 += 8

		// x.ProgramState
		// WARNING: obj.Out[0].ProgramState manually changed from x.ProgramState
		// WARNING: This is not considering program states in different `Out`s with different lengths
		i1 += 4 + uint64(len(obj.Out[0].ProgramState))

		i0 += uint64(len(obj.Out)) * i1
	}

	return i0
}

// encodeTransactionOutputs encodes an object of type transactionOutputs to a buffer allocated to the exact size
// required to encode the object.
func encodeTransactionOutputs(obj *transactionOutputs) ([]byte, error) {
	n := encodeSizeTransactionOutputs(obj)
	buf := make([]byte, n)

	if err := encodeTransactionOutputsToBuffer(buf, obj); err != nil {
		return nil, err
	}

	return buf, nil
}

// encodeTransactionOutputsToBuffer encodes an object of type transactionOutputs to a []byte buffer.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func encodeTransactionOutputsToBuffer(buf []byte, obj *transactionOutputs) error {
	if uint64(len(buf)) < encodeSizeTransactionOutputs(obj) {
		return encoder.ErrBufferUnderflow
	}

	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.Out maxlen check
	if len(obj.Out) > 65535 {
		return encoder.ErrMaxLenExceeded
	}

	// obj.Out length check
	if uint64(len(obj.Out)) > math.MaxUint32 {
		return errors.New("obj.Out length exceeds math.MaxUint32")
	}

	// obj.Out length
	e.Uint32(uint32(len(obj.Out)))

	// obj.Out
	for _, x := range obj.Out {

		// x.Address.Version
		e.Uint8(x.Address.Version)

		// x.Address.Key
		e.CopyBytes(x.Address.Key[:])

		// x.Coins
		e.Uint64(x.Coins)

		// x.Hours
		e.Uint64(x.Hours)

		// x.ProgramState length check
		if uint64(len(x.ProgramState)) > math.MaxUint32 {
			return errors.New("x.ProgramState length exceeds math.MaxUint32")
		}

		// x.ProgramState length
		e.Uint32(uint32(len(x.ProgramState)))

		// x.ProgramState copy
		e.CopyBytes(x.ProgramState)

	}

	return nil
}

// decodeTransactionOutputs decodes an object of type transactionOutputs from a buffer.
// Returns the number of bytes used from the buffer to decode the object.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
func decodeTransactionOutputs(buf []byte, obj *transactionOutputs) (uint64, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.Out

		ul, err := d.Uint32()
		if err != nil {
			return 0, err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return 0, encoder.ErrBufferUnderflow
		}

		if length > 65535 {
			return 0, encoder.ErrMaxLenExceeded
		}

		if length != 0 {
			obj.Out = make([]TransactionOutput, length)

			for z1 := range obj.Out {
				{
					// obj.Out[z1].Address.Version
					i, err := d.Uint8()
					if err != nil {
						return 0, err
					}
					obj.Out[z1].Address.Version = i
				}

				{
					// obj.Out[z1].Address.Key
					if len(d.Buffer) < len(obj.Out[z1].Address.Key) {
						return 0, encoder.ErrBufferUnderflow
					}
					copy(obj.Out[z1].Address.Key[:], d.Buffer[:len(obj.Out[z1].Address.Key)])
					d.Buffer = d.Buffer[len(obj.Out[z1].Address.Key):]
				}

				{
					// obj.Out[z1].Coins
					i, err := d.Uint64()
					if err != nil {
						return 0, err
					}
					obj.Out[z1].Coins = i
				}

				{
					// obj.Out[z1].Hours
					i, err := d.Uint64()
					if err != nil {
						return 0, err
					}
					obj.Out[z1].Hours = i
				}

				{
					// obj.Out[z1].ProgramState

					ul, err := d.Uint32()
					if err != nil {
						return 0, err
					}

					length := int(ul)
					if length < 0 || length > len(d.Buffer) {
						return 0, encoder.ErrBufferUnderflow
					}

					if length != 0 {
						obj.Out[z1].ProgramState = make([]byte, length)

						copy(obj.Out[z1].ProgramState[:], d.Buffer[:length])
						d.Buffer = d.Buffer[length:]
					}
				}
			}
		}
	}

	return uint64(len(buf) - len(d.Buffer)), nil
}

// decodeTransactionOutputsExact decodes an object of type transactionOutputs from a buffer.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
// If the buffer is longer than required to decode the object, returns encoder.ErrRemainingBytes.
func decodeTransactionOutputsExact(buf []byte, obj *transactionOutputs) error {
	if n, err := decodeTransactionOutputs(buf, obj); err != nil {
		return err
	} else if n != uint64(len(buf)) {
		return encoder.ErrRemainingBytes
	}

	return nil
}
