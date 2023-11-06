package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

type cronCollector struct {
	log        log.Logger
	subsystem  string
	instance   string
	jobsStatus *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &cronCollector{
		subsystem: "cron",
	})
}

func (c *cronCollector) Name() string {
	return c.subsystem
}

func (c *cronCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel
	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.jobsStatus = buildPrometheusDesc(c.subsystem, "job_status",
		"Cron job status by name and description (1 = enabled, 0 = disabled)",
		[]string{"schedule", "description", "command", "origin"},
	)
}

func (c *cronCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.jobsStatus
}

func (c *cronCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	crons, err := client.FetchCronTable()
	if err != nil {
		return err
	}
	for _, cron := range crons.Cron {
		ch <- prometheus.MustNewConstMetric(
			c.jobsStatus,
			prometheus.GaugeValue,
			float64(cron.Status),
			cron.Schedule,
			cron.Description,
			cron.Command,
			cron.Origin,
			c.instance,
		)
	}
	return nil
}