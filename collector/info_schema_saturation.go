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

// Scrape `info_schema.global_status for saturation`.

package collector

import (
	"context"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const globalStatusSaturationQuery = `
	SELECT
		sum(variable_value)
		FROM information_schema.global_status
		WHERE Variable_name = 'THREADS_RUNNING' ;
	`

// Metric descriptors.
var (
	informationSchemaGlobalStatusSaturationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "global_status_saturation"),
		"Saturation rate.",
		nil, nil,
	)
)

// ScrapeGlobalStatusSaturationSum collects from `information_schema.global_status`.
type ScrapeGlobalStatusSaturationSum struct{}

// Name of the Scraper. Should be unique.
func (ScrapeGlobalStatusSaturationSum) Name() string {
	return "globalStatusSaturationSum"
}

// Help describes the role of the Scraper.
func (ScrapeGlobalStatusSaturationSum) Help() string {
	return "Returns saturation rate from information_schema.global_status"
}

// Version of MySQL from which scraper is available.
func (ScrapeGlobalStatusSaturationSum) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeGlobalStatusSaturationSum) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) error {
	// Timers here are returned in picoseconds.
	globalStatusSaturationRows, err := db.QueryContext(ctx, globalStatusSaturationQuery)
	if err != nil {
		return err
	}
	defer globalStatusSaturationRows.Close()

	var (
		total  uint64
	)

	for globalStatusSaturationRows.Next() {
		if err := globalStatusSaturationRows.Scan(
			&total,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			informationSchemaGlobalStatusSaturationDesc, prometheus.CounterValue, float64(total),
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapeGlobalStatusSaturationSum{}
