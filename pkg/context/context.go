package context

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ApplicationContext returns a cancellable background context that handles os signals (SIGKILL & SIGTEM)
func ApplicationContext() context.Context {
	ctx, ctxCancel := context.WithCancel(context.Background())

	signal.Ignore()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGKILL|syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalChan

		fmt.Println("signal captured correctly to shutdown !!")
		ctxCancel()

		time.Sleep(time.Second)
		os.Exit(1)
	}()

	return ctx
}
