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

// Scrape `info_schema.global_status`.

package collector

import (
	"context"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const globalStatusRequestRatesQuery = `
	SELECT
		SUM(variable_value) AS TOTAL_REQUEST_RATE
		FROM information_schema.global_status
		WHERE variable_name IN ('com_select', 'com_update', 'com_delete', 'com_insert');
	`

// Metric descriptors.
var (
	informationSchemaGlobalStatusRequestSumTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "global_status_request_rate"),
		"The total count of requests.",
		nil, nil,
	)
)

// ScrapeGlobalStatusRequestRatesSum collects from `information_schema.global_status`.
type ScrapeGlobalStatusRequestRatesSum struct{}

// Name of the Scraper. Should be unique.
func (ScrapeGlobalStatusRequestRatesSum) Name() string {
	return "globalStatusRequestRatesSum"
}

// Help describes the role of the Scraper.
func (ScrapeGlobalStatusRequestRatesSum) Help() string {
	return "Returns request rate from information_schema.global_status"
}

// Version of MySQL from which scraper is available.
func (ScrapeGlobalStatusRequestRatesSum) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeGlobalStatusRequestRatesSum) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) error {
	// Timers here are returned in picoseconds.
	globalStatusRequestRatesRows, err := db.QueryContext(ctx, globalStatusRequestRatesQuery)
	if err != nil {
		return err
	}
	defer globalStatusRequestRatesRows.Close()

	var (
		total  uint64
	)

	for globalStatusRequestRatesRows.Next() {
		if err := globalStatusRequestRatesRows.Scan(
			&total,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			informationSchemaGlobalStatusRequestSumTotalDesc, prometheus.CounterValue, float64(total),
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapeGlobalStatusRequestRatesSum{}
