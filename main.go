package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adeniyistephen/fast-speedtest/app"
	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	log := log.New(os.Stdout, "SPEEDY : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			ReadTimeout     time.Duration `conf:"default:5m"`
			WriteTimeout    time.Duration `conf:"default:5m"`
			ShutdownTimeout time.Duration `conf:"default:5m"`
		}
	}
	cfg.Version.SVN = build

	if err := conf.Parse(os.Args[1:], "SPEEDY", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("SPEEDY", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("SPEEDY", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// App Starting

	expvar.NewString("build").Set(build)
	log.Printf("main: Started: Application initializing: version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config:\n%v\n", out)

	//===================================================
	// Start API Service

	log.Println("main: Initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      app.Routes(build, shutdown, log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Shutdown
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v: Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
