// Package command provides teonet proxy server commands binary protocol.
package command

import (
	"errors"
	"fmt"
)

// Teonet proxy commands
const (
	cmdNone Command = iota
	Connect
	Disconnect
	ConnectTo
	NewAPIClient
	cmdCount
)

var (
	ErrNotEnoughData  = fmt.Errorf("not enough data")
	ErrWrongChecksum  = fmt.Errorf("wrong checksum")
	ErrUnknownCommand = fmt.Errorf("unknown command")
)

// TeonetCmd represents a command packet in the Teonet proxy protocol. It
// contains a command byte and a data slice.
type TeonetCmd struct {
	Cmd  Command
	Data []byte
	Err  error
}

type Command byte

// String returns the string representation of the Command.
//
// It returns a string that represents the value of the Command
// constant. If the value is one of the predefined constants
// (Connect, Dsconnect, ConnectTo, NewAPIClient), it returns
// the corresponding string. Otherwise, it returns "Unknown".
// String is part of the fmt.Stringer interface.
func (c Command) String() string {
	switch c & 0x7F {
	case Connect:
		return "Connect"
	case Disconnect:
		return "Dsconnect"
	case ConnectTo:
		return "ConnectTo"
	case NewAPIClient:
		return "NewAPIClient"
	default:
		return "Unknown"
	}
}

// MarshalBinary converts the TeonetCmd struct into a binary representation.
//
// It returns a byte slice containing the binary representation of the struct
// and an error if there was an issue during the conversion.
func (c TeonetCmd) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 0, len(c.Data)+2)
	

	if c.Err != nil {
		data = append(data, byte(c.Cmd|0x80))
		data = append(data, []byte(c.Err.Error())...)
	} else {
		data = append(data, byte(c.Cmd))
		data = append(data, c.Data...)
	}

	data = append(data, c.checksum(data))

	return
}

// UnmarshalBinary unmarshals binary data into a TeonetCmd object.
//
// The function takes a byte slice `data` as input and unmarshals it into the
// `TeonetCmd` object. It performs the following steps:
//   - Checks the length of the data slice. If it is less than 2, it returns an
//     error `ErrNotEnoughData`.
//   - Checks the checksum of the data. If it does not match the checksum at the
//     end of the data slice, it returns an error `ErrWrongChecksum`.
//   - Sets the command byte from the data slice at index 0.
//   - Sets the data slice from the data slice at index 1 to the second last
//     element.
//   - Returns any error encountered during the unmarshaling process.
//
// Parameters:
// - `data []byte`: The binary data to be unmarshaled.
//
// Returns:
// - `err error`: An error encountered during the unmarshaling process.
func (c *TeonetCmd) UnmarshalBinary(data []byte) (err error) {

	// Check packet length
	if len(data) < 2 {
		err = ErrNotEnoughData
		return
	}

	// Check checksum
	if c.checksum(data[:len(data)-1]) != data[len(data)-1] {
		err = ErrWrongChecksum
		return
	}

	cmd := data[0] & 0x7F

	// Check command number
	if !(cmd > byte(cmdNone) && cmd < byte(cmdCount)) {
		err = ErrUnknownCommand
		return
	}

	// Get command and data
	c.Cmd = Command(cmd)
	if data[0]&0x80 == 0 {
		c.Data = data[1 : len(data)-1]
	} else {
		c.Err = errors.New(string(data[1 : len(data)-1]))
	}

	return
}

// checksum calculates checksum for given data
func (c TeonetCmd) checksum(data []byte) byte {
	var sum byte
	for i := 0; i < len(data); i++ {
		sum += data[i]
	}
	return sum
}
