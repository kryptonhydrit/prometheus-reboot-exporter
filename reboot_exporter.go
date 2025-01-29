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
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	listenAddr  = flag.String("web.telemetry-port", "9001", "Port to listen for telemetry.")
	metricsPath = flag.String("web.telemetry-path", "/metrics", "Path to expose metrics")
)

// Global variable for the custom metric
var rebootRequiredGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "reboot_required",
		Help: "Indicates if a reboot is required (1 = required, 0 = not required).",
	})

// Function to check if the reboot required
func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func checkRebootRequired() {
	go func() {
		for {
			if checkFileExists("/var/run/reboot-required") {
				rebootRequiredGauge.Set(1)
			} else {
				rebootRequiredGauge.Set(0)
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	flag.Parse()

	prometheus.MustRegister(rebootRequiredGauge)

	checkRebootRequired()

	log.Println("Starting reboot_exporter", version.Info())
	log.Printf("Server running at http://0.0.0.0:%s%s", *listenAddr, *metricsPath)
	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":"+*listenAddr, nil))
}
