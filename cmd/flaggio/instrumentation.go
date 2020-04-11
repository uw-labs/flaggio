package main

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func newTracer(jaegerHost string, logger *logrus.Entry) (opentracing.Tracer, io.Closer, error) {
	sender, err := jaeger.NewUDPTransport(jaegerHost, 0)
	if err != nil {
		return nil, nil, err
	}

	metricsFactory := prometheus.New().Namespace(metrics.NSOptions{Name: ApplicationName, Tags: nil})

	jlogger := &jaegerLogger{logger: logger}
	c := jaegercfg.Configuration{
		ServiceName: "flaggio",
	}
	return c.NewTracer(
		jaegercfg.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(5*time.Second),
			jaeger.ReporterOptions.Logger(jlogger),
		)),
		jaegercfg.Logger(jlogger),
		jaegercfg.Metrics(metricsFactory),
		jaegercfg.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
		jaegercfg.Sampler(jaeger.NewRateLimitingSampler(1)),
	)
}

func tracingMiddleware(operationName string, logger *logrus.Entry) func(next http.Handler) http.Handler {
	tracer := opentracing.GlobalTracer()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
				logger.WithError(err).Debug("failed to extract opentracing headers")
			}
			serverSpan := tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx), opentracing.Tags{
				string(ext.Component):  "http",
				string(ext.HTTPMethod): r.Method,
				string(ext.HTTPUrl):    r.URL.String(),
				"request.id":           middleware.GetReqID(r.Context()),
			})
			defer serverSpan.Finish()
			newReq := r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))

			next.ServeHTTP(w, newReq)
		})
	}
}

type jaegerLogger struct {
	logger *logrus.Entry
}

func (l *jaegerLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}
