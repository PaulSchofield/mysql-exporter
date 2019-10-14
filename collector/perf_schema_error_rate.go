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

// Scrape `perf_schema.global_status`.

package collector

import (
	"context"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const queryErrorRateQuery = `
	SELECT 
		sum(sum_errors) AS query_count
	FROM performance_schema.events_statements_summary_by_user_by_event_name WHERE event_name IN ('statement/sql/select', 'statement/sql/insert', 'statement/sql/update', 'statement/sql/delete');
	`

// Metric descriptors.
var (
	perfSchemaQueryErrorRateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "event_statement_summary_by_user_by_event_name"),
		"The total count of query_count errors.",
		nil, nil,
	)
)

// ScrapePerfSchemaQueryErrorRateSum collects from `performance_schema.event_statements_summary_by_user_by_event_name`.
type ScrapePerfSchemaQueryErrorRateSum struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfSchemaQueryErrorRateSum) Name() string {
	return "perfSchemaQueryErrorRateSum"
}

// Help describes the role of the Scraper.
func (ScrapePerfSchemaQueryErrorRateSum) Help() string {
	return "Returns error rate from performance_schema.event_statements_summary_by_user_by_event_name"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfSchemaQueryErrorRateSum) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfSchemaQueryErrorRateSum) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) error {
	// Timers here are returned in picoseconds.
	perfSchemaQueryErrorRateRows, err := db.QueryContext(ctx, queryErrorRateQuery)
	if err != nil {
		return err
	}
	defer perfSchemaQueryErrorRateRows.Close()

	var (
		total uint64
	)

	for perfSchemaQueryErrorRateRows.Next() {
		if err := perfSchemaQueryErrorRateRows.Scan(
			&total,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			perfSchemaQueryErrorRateDesc, prometheus.CounterValue, float64(total),
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfSchemaQueryErrorRateSum{}
