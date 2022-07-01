package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/itkq/sentry-exporter/sentry"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	logger         log.Logger
	client         *sentry.Client
	projectSlugMap map[int]string

	errors       *prometheus.Desc
	transactions *prometheus.Desc
}

const (
	prometheusNamespace = "sentry"
)

func NewCollector(logger log.Logger, client *sentry.Client, projectSlugMap map[int]string) *Collector {
	return &Collector{
		logger:         logger,
		client:         client,
		projectSlugMap: projectSlugMap,

		errors: prometheus.NewDesc(
			prometheus.BuildFQName(prometheusNamespace, "", "errors_total"),
			"Total errors",
			[]string{"project", "project_slug"},
			nil,
		),
		transactions: prometheus.NewDesc(
			prometheus.BuildFQName(prometheusNamespace, "", "transactions_total"),
			"Total transactions",
			[]string{"project", "project_slug"},
			nil,
		),
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.TODO()

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := c.collectErrorsByProject(ctx, ch); err != nil {
			level.Error(c.logger).Log(err)
			fmt.Println(err)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := c.collectTransactionsByProject(ctx, ch); err != nil {
			level.Error(c.logger).Log(err)
			fmt.Println(err)
		}
	}()

	waitGroup.Wait()
}

func (c *Collector) collectErrorsByProject(ctx context.Context, ch chan<- prometheus.Metric) error {
	req := &sentry.RetrieveEventCountsForOrgV2Request{
		Field:       "sum(quantity)",
		StatsPeriod: "1h",
		Interval:    "1h", // default and minimum
		GroupBy:     []string{"project"},
		Category:    "error",
		Project:     []string{"-1"},
	}

	resp, err := c.client.RetrieveEventCountsForOrgV2(ctx, req)
	if err != nil {
		return err
	}

	for _, g := range resp.Groups {
		slug, ok := c.projectSlugMap[g.By.Project]
		if !ok {
			slug = ""
		}
		ch <- prometheus.MustNewConstMetric(
			c.errors,
			prometheus.CounterValue,
			float64(g.Totals["sum(quantity)"]),
			fmt.Sprint(g.By.Project),
			slug,
		)
	}

	return nil
}

func (c *Collector) collectTransactionsByProject(ctx context.Context, ch chan<- prometheus.Metric) error {
	req := &sentry.RetrieveEventCountsForOrgV2Request{
		Field:       "sum(quantity)",
		StatsPeriod: "1h",
		Interval:    "1h", // default and minimum
		GroupBy:     []string{"project"},
		Category:    "transaction",
		Project:     []string{"-1"},
	}

	resp, err := c.client.RetrieveEventCountsForOrgV2(ctx, req)
	if err != nil {
		return err
	}

	for _, g := range resp.Groups {
		slug, ok := c.projectSlugMap[g.By.Project]
		if !ok {
			slug = ""
		}
		ch <- prometheus.MustNewConstMetric(
			c.transactions,
			prometheus.CounterValue,
			float64(g.Totals["sum(quantity)"]),
			fmt.Sprint(g.By.Project),
			slug,
		)
	}

	return nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {}
