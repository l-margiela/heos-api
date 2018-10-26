package heosapi

import (
	"bufio"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// Heos is low-level API for communicating with Denon HEOS speaker family
type Heos struct {
	conn net.Conn
	host string
}

// NewHeos returns Heos instancec
func NewHeos(host string) Heos {
	var heos Heos
	heos.host = host

	return heos
}

// Connect initialises connection to the speaker
func (heos *Heos) Connect() error {
	conn, err := newConn(heos.host)
	if err != nil {
		return err
	}
	heos.conn = conn

	return nil
}

func newConn(host string) (net.Conn, error) {
	return net.Dial("tcp", host)
}

// Disconnect disconnects from the speaker
func (heos *Heos) Disconnect() error {
	return heos.conn.Close()
}

// Send sends given command with parameters to the speaker
func (heos *Heos) Send(cmd Command, params map[string]string) (Response, error) {
	_, err := fmt.Fprintf(heos.conn, "heos://%s/%s?%s\r\n", cmd.Group, cmd.Command, paramsToStr(params))
	if err != nil {
		return Response{}, err
	}

	return heos.readNextResp()
}

func responseSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if len(data) < i+1 {
			return 0, nil, nil
		}
		if data[i] == '\r' && data[i+1] == '\n' {
			return i + 2, data[:i], nil
		}
	}

	return 0, data, bufio.ErrFinalToken
}

func (heos *Heos) readNextResp() (Response, error) {
	scanner := bufio.NewScanner(heos.conn)
	scanner.Split(responseSplit)

	if !scanner.Scan() {
		return Response{}, fmt.Errorf("no response")
	}
	if err := scanner.Err(); err != nil {
		return Response{}, errors.Wrap(err, "reading input")
	}
	return parseResponse([]byte(scanner.Text()))
}

// EventStream provides channel with responses from the speaker.
// Also provides error channel
func (heos *Heos) EventStream() (<-chan Response, <-chan error) {
	errCh := make(chan error, 1)
	ch := make(chan Response)
	success := true

	var resp Response
	var err error

	newHeos := NewHeos(heos.host)
	if connErr := newHeos.Connect(); err != nil {
		errCh <- connErr
		goto done
	}

	resp, err = newHeos.Send(Command{
		Group:   "system",
		Command: "register_for_change_events",
	}, map[string]string{
		"enable": "on",
	})
	if err != nil {
		errCh <- err
		goto done
	}
	if resp.Heos.Result == "fail" {
		respErr, ok := resp.Heos.Message["text"]
		if !ok {
			errCh <- fmt.Errorf("register_for_change_events failed")
			goto done
		}
		errCh <- fmt.Errorf("register_for_change_events failed: %s", respErr)
		goto done
	}

	go eventStream(newHeos, ch, errCh)
	success = true

done:
	if !success {
		close(errCh)
	}
	return ch, errCh
}

func eventStream(heos Heos, ch chan<- Response, errCh chan<- error) {
	scanner := bufio.NewScanner(heos.conn)
	scanner.Split(responseSplit)

	for scanner.Scan() {
		resp, err := parseResponse([]byte(scanner.Text()))
		if err != nil {
			errCh <- err
			goto done
		}
		ch <- resp
	}
	if err := scanner.Err(); err != nil {
		errCh <- err
		goto done
	}

done:
	close(ch)
	close(errCh)
}
