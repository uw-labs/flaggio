package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	ApplicationName        = "flaggio"
	ApplicationDescription = "Self hosted feature flag solution"
	ApplicationVersion     = "0.1.0"
	GitSummary             = ""
	GitBranch              = ""
	BuildStamp             = ""
)

const (
	logFormatterText = "text"
	logFormatterJSON = "json"
)

func build() string {
	return fmt.Sprintf("%s[%s] (%s)", GitBranch, GitSummary, BuildStamp)
}

func main() { // nolint:gocyclo // dependencies
	app := cli.App{
		Name:        ApplicationName,
		Description: ApplicationDescription,
		Version:     ApplicationVersion,
		Flags:       flags,
		Action: func(_ *cli.Context) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logger := logrus.New()
			logLevel, err := logrus.ParseLevel(cfg.logLevel)
			if err != nil {
				return err
			}
			logger.SetLevel(logLevel)
			switch cfg.logFormatter {
			case logFormatterText:
				logger.SetFormatter(new(logrus.TextFormatter))
			case logFormatterJSON:
				logger.SetFormatter(new(logrus.JSONFormatter))
			default:
				return fmt.Errorf("invalid formatter: %s", cfg.logFormatter)
			}

			logger.
				WithFields(logrus.Fields{"version": ApplicationVersion, "build": build()}).
				Infof("starting %s application", ApplicationName)

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			// setup tracer
			if cfg.isTracingEnabled() {
				tracer, closer, err := newTracer(cfg.jaegerAgentHost, logger.WithField("app", "tracer"))
				if err != nil {
					return err
				}
				defer closer.Close()
				opentracing.SetGlobalTracer(tracer)
			}

			errs := make(chan error, 1)
			var wg sync.WaitGroup
			if !cfg.noAPI {
				// start API server
				go func() {
					err := startAPI(ctx, &wg, logger.WithField("app", "api"))
					if err != nil {
						errs <- err
					}
				}()
			}
			if !cfg.noAdmin {
				// start Admin server
				go func() {
					err := startAdmin(ctx, &wg, logger.WithField("app", "admin"))
					if err != nil {
						errs <- err
					}
				}()
			}

			for {
				select {
				case err := <-errs:
					cancel()
					wg.Wait()
					return err
				case <-done:
					logger.Debug("got os.Interrupt, cancelling main context")
					cancel()
				case <-ctx.Done():
					logger.Trace("context done")
					wg.Wait()
					logger.Info("shutdown completed")
					return ctx.Err()
				}
			}
		},
	}

	if err := app.Run(os.Args); err != nil && !errors.Is(err, context.Canceled) {
		logrus.Fatal(err)
	}
}
