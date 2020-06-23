package main

import (
	"./faq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	faq.Start("./config.yml")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
