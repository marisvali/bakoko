package world

import (
	"bytes"
	"net"
	. "playful-patterns.com/bakoko/ints"
)

func ReadData(conn net.Conn) ([]byte, error) {
	// Create a buffer to read data into.
	buffer := make([]byte, 1024)
	// Create the slice where we'll keep our full message.
	var fullMessage []byte
	// Expected length == 0 means we don't know the expected length yet.
	expectedLen := I(0)

	for {
		// Read (some) data from the client.
		nBytesAcuallyReceived, err := conn.Read(buffer)
		if err != nil {
			return fullMessage, err
		}

		// Append whatever data we got to our full message.
		fullMessage = append(fullMessage, buffer[:nBytesAcuallyReceived]...)

		// Get the expected length (if possible).
		// The first int64 indicates the length of the rest of the message.
		// Check if we got at least an int64.
		if expectedLen.Eq(I(0)) && len(fullMessage) >= 8 {
			// Yes, we have it, now we can extract the expected length.
			Deserialize(bytes.NewBuffer(fullMessage), &expectedLen)
		}

		// Check if we got the expected number of bytes or not.
		if I(len(fullMessage)).Geq(expectedLen) {
			return fullMessage[8:], nil // Skip the length from the gui-data.
		}
	}
}

func WriteData(conn net.Conn, data []byte) error {
	// Docs for the io.Writer interface says that if the entire data hasn't
	// been written then err will be non-nil.

	// Write the total length of the data first.
	buf := new(bytes.Buffer)
	dataLen := int64(len(data) + 8) // include the dataLen field in the total length
	Serialize(buf, dataLen)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// Write the data itself.
	_, err = conn.Write(data)
	return err
}
