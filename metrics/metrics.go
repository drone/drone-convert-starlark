// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ConversionCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "drone_starlark_conversion_request_count",
		Help: "Total number of Starlark pipelines submitted for conversion",
	})
	ConversionSuccessCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "drone_starlark_conversion_success_count",
		Help: "Total number of Starlark pipelines successfully converted",
	})
	ConversionFailCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "drone_starlark_conversion_fail_count",
		Help: "Total number of Starlark pipelines unsuccessfully converted",
	})
	ConversionLatency = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "drone_starlark_conversion_latencies",
		Help: "The time required to service each conversion request",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)
