package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	AppName        = "flaggio"
	AppDescription = ""
	AppVersion     = ""
)

func main() {
	app := cli.App{
		Name:        AppName,
		Description: AppDescription,
		Version:     AppVersion,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "database-uri",
				Usage:  "Database URI",
				EnvVar: "DATABASE_URI",
				Value:  "mongodb://localhost:27017/flaggio",
			},
			cli.StringFlag{
				Name:   "build-path",
				Usage:  "UI build absolute path",
				EnvVar: "BUILD_PATH",
			},
			cli.StringSliceFlag{
				Name:   "cors-allowed-origins",
				Usage:  "CORS allowed origins separated by comma",
				EnvVar: "CORS_ALLOWED_ORIGINS",
			},
			cli.BoolFlag{
				Name:   "cors-debug",
				Usage:  "CORS debug logging",
				EnvVar: "CORS_DEBUG",
				Hidden: true,
			},
			cli.StringFlag{
				Name:   "api-port",
				Usage:  "Port the API server will listen to",
				EnvVar: "API_PORT",
				Value:  "25880",
			},
			cli.StringFlag{
				Name:   "admin-port",
				Usage:  "Port the admin server will listen to",
				EnvVar: "ADMIN_PORT",
				Value:  "25881",
			},
			cli.BoolFlag{
				Name:   "no-api",
				Usage:  "Don't start the API server",
				EnvVar: "NO_API",
			},
			cli.BoolFlag{
				Name:   "no-admin",
				Usage:  "Don't start the admin server",
				EnvVar: "NO_ADMIN",
			},
			cli.BoolFlag{
				Name:   "no-admin-ui",
				Usage:  "Don't start the admin UI",
				EnvVar: "NO_ADMIN_UI",
			},
			cli.StringFlag{
				Name:   "app-env",
				Usage:  "Environment this application is running on. Valid values are: dev, production",
				EnvVar: "APP_ENV",
				Value:  "production",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			var adminSrv, apiSrv *http.Server
			if !c.Bool("no-api") {
				srv, err := startAPI(ctx, c)
				if err != nil {
					return err
				}
				apiSrv = srv
			}
			if !c.Bool("no-admin") {
				srv, err := startAdmin(ctx, c)
				if err != nil {
					return err
				}
				adminSrv = srv
			}
			<-done
			cancel()

			logrus.Info("shutting down server ...")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if adminSrv != nil {
				if err := adminSrv.Shutdown(shutdownCtx); err != nil {
					logrus.Fatalf("admin server shutdown failed: %+v", err)
				}
			}
			if apiSrv != nil {
				if err := apiSrv.Shutdown(shutdownCtx); err != nil {
					logrus.Fatalf("api server shutdown failed: %+v", err)
				}
			}
			logrus.Info("shutdown complete")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
