package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/islamghany/service/foundation/logger"
)

func main() {
	// ----------------------------------------------------------
	// intialize the logger
	loggerEvents := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			fmt.Printf("ERROR: %s\n", r.Message)
		},
		Info: func(ctx context.Context, r logger.Record) {
			fmt.Printf("INFO: %s\n", r.Message)
		},
	}
	tracerIDFn := func(ctx context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	}
	log := logger.NewWithEvents(os.Stdout, logger.LevelInfo, "sales-api", tracerIDFn, loggerEvents)
	// -----------------------------------------------------------
	ctx := context.Background()
	log.Debug(ctx, "main", "startup", "main is started")
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -----------------------------------------------------------
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -----------------------------------------------------------
	// Shutdown
	shutdown := make(chan os.Signal, 1)
	// Notify the shutdown channel when it is time to shutdown.
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown

	log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
	defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}
