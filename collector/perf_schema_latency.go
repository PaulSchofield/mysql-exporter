// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Scrape `perf_schema.events_statements_summary_global_by_event_name for latency`.

package collector

import (
	"context"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const queryLatencyRateQuery = `
	SELECT 
		(avg_timer_wait)/1e9 AS avg_latency_ms FROM performance_schema.events_statements_summary_global_by_event_name
	WHERE event_name = 'statement/sql/select';
	TRUNCATE TABLE performance_schema.events_statements_summary_global_by_event_name ;
	`

// Metric descriptors.
var (
	perfSchemaQueryLatencyRateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "event_statement_summary_by_user_by_event_name"),
		"Show Latency",
		nil, nil,
	)
)

// ScrapePerfSchemaQueryLatencyRateSum collects from `performance_schema.event_statements_summary_by_user_by_event_name`.
type ScrapePerfSchemaQueryLatencyRateSum struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfSchemaQueryLatencyRateSum) Name() string {
	return "perfSchemaQueryLatencyRateSum"
}

// Help describes the role of the Scraper.
func (ScrapePerfSchemaQueryLatencyRateSum) Help() string {
	return "Returns Latency rate from performance_schema.event_statements_summary_by_user_by_event_name"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfSchemaQueryLatencyRateSum) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfSchemaQueryLatencyRateSum) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) Latency {
	// Timers here are returned in picoseconds.
	perfSchemaQueryLatencyRateRows, err := db.QueryContext(ctx, queryLatencyRateQuery)
	if err != nil {
		return err
	}
	defer perfSchemaQueryLatencyRateRows.Close()

	var (
		total uint64
	)

	for perfSchemaQueryLatencyRateRows.Next() {
		if err := perfSchemaQueryLatencyRateRows.Scan(
			&total,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			perfSchemaQueryLatencyRateDesc, prometheus.CounterValue, float64(total),
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfSchemaQueryLatencyRateSum{}
