package golib

import (
	"context"
	"gitlab.com/golibs-starter/golib/log"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Graceful(app *fx.App) {
	signChan := make(chan os.Signal, 1)
	go func() {
		app.Run()
	}()
	signal.Notify(signChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer func() {
		log.Infof("Close another connection")
		cancel()
	}()
	log.Infof("Stopping App")
	if err := app.Stop(ctx); err == context.DeadlineExceeded {
		log.Infof("Halted active connections")
	}
	close(signChan)
	log.Infof("App Stopped")
}
