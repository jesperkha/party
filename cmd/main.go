package main

import (
	"log"
	"os"
	"syscall"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/pipoker/config"
	"github.com/jesperkha/pipoker/server"
)

func main() {
	notif := notifier.New()

	config := config.Load()
	s := server.New(config)

	go s.ListenAndServe(notif)

	notif.NotifyOnSignal(os.Interrupt, syscall.SIGTERM)
	log.Println("shutdown")
}
