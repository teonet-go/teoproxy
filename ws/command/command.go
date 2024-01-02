// Package command provides teonet proxy server commands binary protocol.
package command

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Teonet proxy commands
const (
	cmdNone      Command = iota
	Connect              // Connect to Teonet
	Disconnect           // Disconnect from Teonet
	ConnectTo            // Connect to peer
	NewAPIClient         // New APIClient
	SendTo               // Send api command to peer
	cmdCount             // Number of commands
)

var (
	ErrNotEnoughData  = fmt.Errorf("not enough data")
	ErrWrongChecksum  = fmt.Errorf("wrong checksum")
	ErrUnknownCommand = fmt.Errorf("unknown command")
)

// TeonetCmd represents a command packet in the Teonet proxy protocol. It
// contains a command byte and a data slice.
type TeonetCmd struct {
	Id   uint32  // Packet ID
	Cmd  Command // Command
	Data []byte  // Data
	Err  error   // Error
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
	case SendTo:
		return "SendTo"
	default:
		return "Unknown"
	}
}

// New creates and returns a new TeonetCmd instance.
//
// It takes a `cmd` parameter of type `Command`, which represents the command
// for the TeonetCmd instance. The `data` parameter is a byte slice that
// contains the data for the TeonetCmd instance.
//
// The function returns a pointer to the newly created TeonetCmd instance.
func New(cmd Command, data []byte) *TeonetCmd {
	return &TeonetCmd{Cmd: cmd, Data: data}
}

// NewEmpty returns a new instance of TeonetCmd with no initial values.
//
// Returns a pointer to TeonetCmd.
func NewEmpty() *TeonetCmd {
	return &TeonetCmd{}
}

// MarshalBinary converts the TeonetCmd struct into a binary representation.
//
// It returns a byte slice containing the binary representation of the struct
// and an error if there was an issue during the conversion.
func (c TeonetCmd) MarshalBinary() (data []byte, err error) {

	// Add packet id
	idBinarySlice := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBinarySlice, c.Id)
	data = append(data, idBinarySlice...)

	// Add command and data
	if c.Err != nil {
		data = append(data, byte(c.Cmd|0x80))
		data = append(data, []byte(c.Err.Error())...)
	} else {
		data = append(data, byte(c.Cmd))
		data = append(data, c.Data...)
	}

	// Add checksum
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

	const (
		idLen   = 4               // Length of packet id
		cmdLen  = 1               // Length of command byte
		cmdIdx  = idLen           // Index of command byte
		dataIdx = cmdIdx + cmdLen // Index of data
	)

	// Check packet length
	if len(data) < idLen+cmdLen+1 {
		err = ErrNotEnoughData
		return
	}

	// Check checksum
	if c.checksum(data[:len(data)-1]) != data[len(data)-1] {
		err = ErrWrongChecksum
		return
	}

	cmd := data[cmdIdx] & 0x7F      // Command
	isErr := data[cmdIdx]&0x80 != 0 // The data contains error message

	// Check command number
	if !(cmd > byte(cmdNone) && cmd < byte(cmdCount)) {
		err = ErrUnknownCommand
		return
	}

	// Get packet id
	c.Id = binary.LittleEndian.Uint32(data[:idLen])

	// Get command and data or error message
	c.Cmd = Command(cmd)
	if !isErr {
		c.Data = data[dataIdx : len(data)-1]
	} else {
		c.Err = errors.New(string(data[dataIdx : len(data)-1]))
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
