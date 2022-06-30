package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log/level"
	"github.com/itkq/sentry-exporter/sentry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress          = kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests.").Default(":9115").Envar("WEB_LISTEN_ADDRESS").String()
	sentryApiKey           = kingpin.Flag("sentry.api-key", "The API key of Sentry").Required().Envar("SENTRY_API_KEY").String()
	sentryApiEndpoint      = kingpin.Flag("sentry.api-endpoint", "The API endpoint of Sentry").Default("https://sentry.io/api/0/").Envar("SENTRY_API_ENDPOINT").String()
	sentryOrganizationSlug = kingpin.Flag("sentry.organization-slug", "An organization slug of Sentry").Required().Envar("SENTRY_ORGANIZATION_SLUG").String()
)

func main() {
	os.Exit(run())
}
func run() int {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("sentry_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	client := sentry.NewClient(*sentryApiKey, *sentryApiEndpoint, *sentryOrganizationSlug)
	collector := NewCollector(logger, client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    *listenAddress,
		Handler: mux,
	}

	srvc := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			return 0
		case <-srvc:
			return 1
		}
	}
}
