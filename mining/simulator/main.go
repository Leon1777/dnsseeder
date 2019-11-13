package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/daglabs/btcd/signal"
	"github.com/daglabs/btcd/util/panics"
)

func main() {
	defer panics.HandlePanic(log, backendLog, nil)
	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command-line arguments: %s", err)
		os.Exit(1)
	}

	if cfg.Verbose {
		enableRPCLogging()
	}

	addressList, err := getAddressList(cfg)
	if err != nil {
		panic(errors.Errorf("Couldn't load address list: %s", err))
	}

	clients, err := connectToServers(cfg, addressList)
	if err != nil {
		panic(errors.Errorf("Error connecting to servers: %s", err))
	}
	defer disconnect(clients)

	spawn(func() {
		err = mineLoop(clients)
		if err != nil {
			panic(errors.Errorf("Error in main loop: %s", err))
		}
	})

	interrupt := signal.InterruptListener()
	<-interrupt
}

func disconnect(clients []*simulatorClient) {
	for _, client := range clients {
		client.Disconnect()
	}
}
