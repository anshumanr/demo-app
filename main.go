package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"demo-app/v1/api"
	log "demo-app/v1/logger"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	logger := log.CreateLogger("../../logs/proxy.log", 20, 1, 28, false)

	r := api.GetServerInstance(9876, logger)
	r.StartServer()

	<-done

}
