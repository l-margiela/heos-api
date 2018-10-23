package heosapi

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Response holds device's response
type Response struct {
	Heos    HeosResp    `json:"heos"`
	Payload PayloadResp `json:"payload"`
}

// HeosResp holds the first part of the device's response
type HeosResp struct {
	Command Command           `json:"command"`
	Result  string            `json:"result"`
	Message map[string]string `json:"message"`
}

type heosResp struct {
	Command string `json:"command"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

// UnmarshalJSON implements json.Unmarshaler interface
func (heos *HeosResp) UnmarshalJSON(raw []byte) error {
	var resp heosResp
	if err := json.Unmarshal(raw, &resp); err != nil {
		return err
	}

	cmd, err := parseHeosCommand(resp.Command)
	if err != nil {
		return err
	}
	heos.Command = cmd
	heos.Result = resp.Result

	msg, err := parseHeosMessage(resp.Message)
	if err != nil {
		return err
	}
	heos.Message = msg

	return nil
}

func parseHeosCommand(cmd string) (Command, error) {
	split := strings.Split(cmd, "/")

	return Command{
		Group:   split[0],
		Command: split[1],
	}, nil
}

func parseHeosMessage(msg string) (map[string]string, error) {
	split := strings.Split(msg, "&")

	parsed := map[string]string{}
	for _, pair := range split {
		pairS := strings.SplitN(pair, "=", 2)

		k := pairS[0]
		v := pairS[1]
		parsed[k] = v
	}

	return parsed, nil
}

// PayloadResp hold the second part of the device's response
type PayloadResp interface{}

func parseResponse(respRaw []byte) (Response, error) {
	event := bytes.TrimRight(respRaw, "\x00")

	var resp Response
	if err := json.Unmarshal(event, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
