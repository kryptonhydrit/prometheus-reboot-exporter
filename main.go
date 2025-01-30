// Copyright 2025 kryptonhydrit
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
)

var (
	listenAddr  string
	metricsPath string
)

// Global variable for the custom metric
var rebootRequiredGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "reboot_required",
		Help: "Indicates if a reboot is required (1 = required, 0 = not required).",
	})

func checkRebootRequired() {
	go func() {
		for {
			if _, err := os.Stat("/var/run/reboot-required"); err == nil {
				rebootRequiredGauge.Set(1)
			} else {
				rebootRequiredGauge.Set(0)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
}

func main() {
	var rootCmd = &cobra.Command{
		Use:  "reboot_exporter",
		Long: "Export metrics about the reboot status from Debian-based systems",
		Run: func(cmd *cobra.Command, args []string) {
			registry := prometheus.NewRegistry()
			registry.MustRegister(rebootRequiredGauge)

			checkRebootRequired()

			log.Println("Starting reboot_exporter", version.Info())
			log.Printf("Server running at http://0.0.0.0:%s%s", listenAddr, metricsPath)
			http.Handle(metricsPath, promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				},
			))
			log.Fatalln(http.ListenAndServe(":"+listenAddr, nil))
		},
	}

	rootCmd.PersistentFlags().StringVar(&listenAddr, "web.telemetry-port", "11011", "Port on which to expose metrics")
	rootCmd.PersistentFlags().StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
