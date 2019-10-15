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

// Scrape `info_schema.global_status for utilization`.

package collector

import (
	"context"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const globalStatusUtilizationQuery = `
	SELECT
		variable_value FROM information_schema.global_status
		WHERE variable_name = 'INNODB_ROWS_READ' ;
	`

// Metric descriptors.
var (
	informationSchemaGlobalStatusUtilizationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "global_status_utilization"),
		"Utilization rate.",
		nil, nil,
	)
)

// ScrapeGlobalStatusUtilizationSum collects from `information_schema.global_status`.
type ScrapeGlobalStatusUtilizationSum struct{}

// Name of the Scraper. Should be unique.
func (ScrapeGlobalStatusUtilizationSum) Name() string {
	return "globalStatusUtilizationSum"
}

// Help describes the role of the Scraper.
func (ScrapeGlobalStatusUtilizationSum) Help() string {
	return "Returns utilization rate from information_schema.global_status"
}

// Version of MySQL from which scraper is available.
func (ScrapeGlobalStatusUtilizationSum) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeGlobalStatusUtilizationSum) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) error {
	// Timers here are returned in picoseconds.
	globalStatusUtilizationRows, err := db.QueryContext(ctx, globalStatusUtilizationQuery)
	if err != nil {
		return err
	}
	defer globalStatusUtilizationRows.Close()

	var (
		total  uint64
	)

	for globalStatusUtilizationRows.Next() {
		if err := globalStatusUtilizationRows.Scan(
			&total,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			informationSchemaGlobalStatusUtilizationDesc, prometheus.CounterValue, float64(total),
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapeGlobalStatusUtilizationSum{}
