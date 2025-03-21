package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NicholasMRicci/colmi-client/lib"
	"tinygo.org/x/bluetooth"
)

func main() {
	lib.Must(bluetooth.DefaultAdapter.Enable())

	server := lib.NewServer()
	server.Start()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	server.Stop()
}
