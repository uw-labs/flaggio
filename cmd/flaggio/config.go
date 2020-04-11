package main

import (
	"github.com/urfave/cli/v2"
)

type config struct {
	databaseURI, redisURI                  string
	apiAddr, adminAddr, uiBuildPath        string
	logFormatter, logLevel                 string
	corsAllowedOrigins, corsAllowedHeaders cli.StringSlice
	corsDebug, noAPI, noAdmin, noAdminUI   bool
	playgroundEnabled                      bool
	jaegerAgentHost                        string
}

func (c *config) isCachingEnabled() bool {
	return c.redisURI != ""
}

func (c *config) isTracingEnabled() bool {
	return c.jaegerAgentHost != ""
}

var cfg = config{}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "database-uri",
		Usage:       "Database URI",
		EnvVars:     []string{"DATABASE_URI"},
		Destination: &cfg.databaseURI,
		Required:    true,
	},
	&cli.StringFlag{
		Name:        "redis-uri",
		Usage:       "Redis URI",
		EnvVars:     []string{"REDIS_URI"},
		Destination: &cfg.redisURI,
	},
	&cli.StringFlag{
		Name:        "build-path",
		Usage:       "UI build absolute path",
		EnvVars:     []string{"BUILD_PATH"},
		Destination: &cfg.uiBuildPath,
	},
	&cli.StringSliceFlag{
		Name:        "cors-allowed-origins",
		Usage:       "CORS allowed origins separated by comma",
		EnvVars:     []string{"CORS_ALLOWED_ORIGINS"},
		Destination: &cfg.corsAllowedOrigins,
	},
	&cli.StringSliceFlag{
		Name:        "cors-allowed-headers",
		Usage:       "CORS allowed headers",
		EnvVars:     []string{"CORS_ALLOWED_HEADERS"},
		Destination: &cfg.corsAllowedHeaders,
	},
	&cli.BoolFlag{
		Name:        "cors-debug",
		Usage:       "CORS debug logging",
		EnvVars:     []string{"CORS_DEBUG"},
		Hidden:      true,
		Destination: &cfg.corsDebug,
	},
	&cli.BoolFlag{
		Name:        "no-api",
		Usage:       "Don't start the API server",
		EnvVars:     []string{"NO_API"},
		Destination: &cfg.noAPI,
	},
	&cli.BoolFlag{
		Name:        "no-admin",
		Usage:       "Don't start the admin server",
		EnvVars:     []string{"NO_ADMIN"},
		Destination: &cfg.noAdmin,
	},
	&cli.BoolFlag{
		Name:        "no-admin-ui",
		Usage:       "Don't start the admin UI",
		EnvVars:     []string{"NO_ADMIN_UI"},
		Destination: &cfg.noAdminUI,
	},
	&cli.BoolFlag{
		Name:        "playground",
		Usage:       "Enable graphql playground",
		EnvVars:     []string{"PLAYGROUND"},
		Destination: &cfg.playgroundEnabled,
	},
	&cli.StringFlag{
		Name:        "api-addr",
		Usage:       "Sets the bind address for the API",
		EnvVars:     []string{"API_ADDR"},
		Value:       ":8080",
		Destination: &cfg.apiAddr,
	},
	&cli.StringFlag{
		Name:        "admin-addr",
		Usage:       "Sets the bind address for the admin",
		EnvVars:     []string{"ADMIN_ADDR"},
		Value:       ":8081",
		Destination: &cfg.adminAddr,
	},
	&cli.StringFlag{
		Name:        "log-formatter",
		Usage:       "Sets the log formatter for the application. Valid values are: text, json",
		EnvVars:     []string{"LOG_FORMATTER"},
		Value:       logFormatterJSON,
		Destination: &cfg.logFormatter,
	},
	&cli.StringFlag{
		Name:        "log-level",
		Usage:       "Sets the log level for the application",
		EnvVars:     []string{"LOG_LEVEL"},
		Value:       "info",
		Destination: &cfg.logLevel,
	},
	&cli.StringFlag{
		Name:        "jaeger-agent-host",
		Usage:       "The address of the jaeger agent (host:port)",
		EnvVars:     []string{"JAEGER_AGENT_HOST"},
		Destination: &cfg.jaegerAgentHost,
	},
}
