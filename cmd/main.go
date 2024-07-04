package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/ahsansandiah/dpo-test/packages/server"
)

func run() error {
	mgr, err := manager.NewInit()
	if err != nil {
		return err
	}

	// app config
	tzLocation, err := time.LoadLocation(mgr.GetConfig().AppTz)
	if err != nil {
		return err
	}
	time.Local = tzLocation

	// server config
	server := server.NewServer(mgr.GetConfig())

	// start routes
	// end routes

	server.RegisterRouter(server.Router)

	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
