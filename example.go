package heosapi

import (
	"fmt"
	"os"
)

func example() {
	var heos Heos

	if err := heos.Connect("1.1.1.1:1255"); err != nil {
		fmt.Printf("connect: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := heos.Disconnect(); err != nil {
			fmt.Printf("disconnect: %s\n", err)
			os.Exit(1)
		}
	}()

	if err := heos.Send(Command{
		Group:   "system",
		Command: "register_for_change_events",
	}, map[string]string{
		"enable": "on",
	}); err != nil {
		fmt.Printf("send: %s\n", err)
		os.Exit(1)
	}

	ch, errCh := heos.EventStream()
	for {
		select {
		case resp := <-ch:
			fmt.Printf("%#+v\n", resp)
		case err := <-errCh:
			fmt.Printf("receive: %s\n", err)
			os.Exit(1)
		}
	}
}
