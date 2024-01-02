package command

import (
	"bytes"
	"errors"
	"testing"
)

func TestUnmarshalBinary(t *testing.T) {
	// Create a new TeonetCmd instance
	cmd := &TeonetCmd{}

	// Test case 1: valid binary data
	data := []byte{0, 0, 0, 0, 1, 2, 3}
	err := cmd.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case 2: valid binary data
	data = []byte{0, 0, 0, 0, 1, 1}
	err = cmd.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case 3: invalid packet length
	data = []byte{1}
	err = cmd.UnmarshalBinary(data)
	if err != ErrNotEnoughData {
		t.Errorf("Expected ErrNotEnoughData, got: %v", err)
	}

	// Test case 4: invalid checksum
	data = []byte{0, 0, 0, 0, 1, 2, 4}
	err = cmd.UnmarshalBinary(data)
	if err != ErrWrongChecksum {
		t.Errorf("Expected ErrWrongChecksum, got: %v", err)
	}

	// Test case 4: unknown command
	data = []byte{0, 0, 0, 0, 10, 2, 12}
	err = cmd.UnmarshalBinary(data)
	if err != ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got: %v", err)
	}

	// Test case 5: check message id
	data = []byte{1, 0, 0, 0, 1, 2, 4}
	err = cmd.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if cmd.Id != 1 {
		t.Errorf("Expected id: 1, got: %v", cmd.Id)
	}
}

func TestMarshalBinary(t *testing.T) {

	// Test case 1: Create a new TeonetCmd instance
	cmd := &TeonetCmd{Cmd: Connect}
	data, err := cmd.MarshalBinary()
	if err != nil {
		t.Errorf("error converting to binary, data: %v, error: %v", data, err)
	}

	// Test case 2: Create a new TeonetCmd instance
	cmd = &TeonetCmd{Cmd: Connect, Data: []byte("hello")}
	data, err = cmd.MarshalBinary()
	if err != nil {
		t.Errorf("error converting to binary, data: %v, error: %v", data, err)
	}

	// Unmarshal the binary data and compare it with the original TeonetCmd
	unmarshaledCmd := &TeonetCmd{}
	err = unmarshaledCmd.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("error unmarshaling binary data: %v", err)
	}

	if unmarshaledCmd.Cmd != cmd.Cmd {
		t.Errorf("expected cmd: %v, got: %v", cmd.Cmd, unmarshaledCmd.Cmd)
	}

	if !bytes.Equal(unmarshaledCmd.Data, cmd.Data) {
		t.Errorf("expected data: %v, got: %v", cmd.Data, unmarshaledCmd.Data)
	}

	// Test case 3: Create a new TeonetCmd instance with error
	cmd = &TeonetCmd{Cmd: Connect, Err: errors.New("test error")}
	data, err = cmd.MarshalBinary()
	if err != nil {
		t.Errorf("error converting to binary, data: %v, error: %v", data, err)
	}

	// Unmarshal the binary data and compare it with the original TeonetCmd
	unmarshaledCmd = &TeonetCmd{}
	err = unmarshaledCmd.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("error unmarshaling binary data: %v", err)
	}

	if unmarshaledCmd.Cmd != cmd.Cmd {
		t.Errorf("expected cmd: %v, got: %v", cmd.Cmd, unmarshaledCmd.Cmd)
	}

	if unmarshaledCmd.Err.Error() != cmd.Err.Error() {
		t.Errorf("expected error: %v, got: %v", cmd.Data, unmarshaledCmd.Data)
	}
}
