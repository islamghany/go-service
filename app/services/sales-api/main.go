package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/islamghany/service/foundation/logger"
)

var build = "develop"

func main() {
	// ----------------------------------------------------------
	// intialize the logger
	loggerEvents := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			fmt.Printf("ERROR: %s\n", r.Message)
		},
	}
	tracerIDFn := func(ctx context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	}
	log := logger.NewWithEvents(os.Stdout, logger.LevelInfo, "sales-api", tracerIDFn, loggerEvents)
	// -----------------------------------------------------------
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -----------------------------------------------------------
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	// -----------------------------------------------------------
	// Configuration
	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s,mask"`
			APIHost            string        `conf:"default:0.0.0.0:8000"`
			DebugHost          string        `conf:"default:0.0.0.0:4000"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Service Project",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(build)

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
