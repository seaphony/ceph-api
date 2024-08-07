package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	stdlog "github.com/rs/zerolog/log"
	"github.com/seaphony/ceph-api/pkg/app"
	"github.com/seaphony/ceph-api/pkg/config"
)

// this information will be collected when built, by -ldflags="-X 'main.version=$(tag)' -X 'main.commit=$(commit)'".
var (
	version            = "development"
	commit             = "not set"
	configPath         = flag.String("config", "", "set path to config directory")
	configOverridePath = flag.String("config-override", "", "set path to config override directory")
)

func main() {
	flag.Parse()
	var configs []config.Src
	if configPath != nil && *configPath != "" {
		configs = append(configs, config.Path(*configPath))
	}
	if configOverridePath != nil && *configOverridePath != "" {
		configs = append(configs, config.Path(*configOverridePath))
	}

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() {
		<-signals
		zerolog.Ctx(ctx).Info().Msg("received shutdown signal.")
		cancel()
	}()
	var conf config.Config
	err := config.Get(&conf, configs...)
	if err != nil {
		stdlog.Fatal().Err(err).Msg("critical error. Unable to read app config")
	}

	err = app.Start(ctx, conf, config.Build{
		Version: version,
		Commit:  commit,
	})
	if err != nil {
		stdlog.Err(err).Msg("critical error. Shutdown application")
		os.Exit(1)
	}
}
