package microcontroller

import (
	"bytes"
	"errors"
	"time"

	"github.com/tarm/serial"
)

type Arduino struct {
	port string
	conn *serial.Port
}

func NewArduino(port string, baud int) (ard *Arduino, err error) {
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	ard = &Arduino{
		port: port,
		conn: s,
	}

	return ard, nil
}

// Send performs several important checks before appending a newline and sending
// the message to the microcontroller.
func (ard *Arduino) Send(msg []byte) (err error) {
	if bytes.Count(msg, []byte("\n")) > 0 {
		return errors.New("msg must not contain a new line")
	}

	if len(msg) > 50 {
		return errors.New("msg must not be longer than 50 bytes")
	}

	if len(msg) == 0 {
		return errors.New("msg is empty, sending nothing")
	}

	msg = append(msg, '\n')

	_, err = ard.conn.Write(msg)
	if err != nil {
		return err
	}

	return ard.conn.Flush()
}

// Receive uses a 1KiB buffer and fills it until a newline is read. The returned
// msg does not contain the newline.
func (ard *Arduino) Receive() (msg []byte, err error) {
	buffer := make([]byte, 1024)
	count, err := ard.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	buffer = bytes.TrimSuffix(buffer[:count], []byte("\n"))

	return buffer, nil
}

// Request is a helper function to send a message and wait for a response
func (ard *Arduino) Request(msg []byte) (resp []byte, err error) {
	err = ard.Send(msg)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 500)

	return ard.Receive()
}
