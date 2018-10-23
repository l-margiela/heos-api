package heosapi

import (
	"fmt"
	"net"
)

// Heos is low-level API for communicating with Denon HEOS speaker family
type Heos struct {
	conn net.Conn
}

// Connect initialises connection to the speaker
func (heos *Heos) Connect(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	heos.conn = conn

	return nil
}

// Disconnect disconnects from the speaker
func (heos *Heos) Disconnect() error {
	return heos.conn.Close()
}

// Send sends given command with parameters to the speaker
func (heos *Heos) Send(cmd Command, params map[string]string) error {
	_, err := fmt.Fprintf(heos.conn, "heos://%s/%s?%s\r\n", cmd.Group, cmd.Command, paramsToStr(params))

	return err
}

// Recv provides channel with responses from the speaker.
// Also provides error channel
func (heos *Heos) Recv() (<-chan Response, <-chan error) {
	ch := make(chan Response)
	errCh := make(chan error)
	go heos.recv(ch, errCh)

	return ch, errCh
}

func (heos *Heos) recv(ch chan<- Response, errCh chan<- error) {
	for {
		event := make([]byte, 10000)
		if _, err := heos.conn.Read(event); err != nil {
			errCh <- err
			break
		}

		resp, err := parseResponse(event)
		if err != nil {
			errCh <- err
			break
		}

		ch <- resp
	}
	close(ch)
	close(errCh)
}
