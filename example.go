package heosapi

import (
	"fmt"
	"os"
)

func example() {
	heos := NewHeos("1.1.1.1:1255")

	if err := heos.Connect(); err != nil {
		fmt.Printf("connect: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := heos.Disconnect(); err != nil {
			fmt.Printf("disconnect: %s\n", err)
			os.Exit(1)
		}
	}()

	resp, err := heos.Send(Command{
		Group:   "system",
		Command: "heart_beat",
	}, map[string]string{})
	if err != nil {
		fmt.Printf("send: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", resp)
}
