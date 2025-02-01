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
	listenAddr   string
	metricsPath  string
	sentinalPath string
)

var rebootRequiredGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "reboot_required",
		Help: "Indicates if a reboot is required (1 = required, 0 = not required).",
	})

func main() {
	var rootCmd = &cobra.Command{
		Use:  "reboot_exporter",
		Long: "Export metrics about the reboot status from Debian-based systems",
		Run: func(cmd *cobra.Command, args []string) {
			registry := prometheus.NewRegistry()
			registry.MustRegister(rebootRequiredGauge)

			go func() {
				for {
					if _, err := os.Stat(sentinalPath); err == nil {
						rebootRequiredGauge.Set(1)
					} else {
						rebootRequiredGauge.Set(0)
					}
					time.Sleep(1 * time.Minute)
				}
			}()

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

	rootCmd.PersistentFlags().StringVar(&listenAddr, "web.metrics.port", "11011", "Port on which to expose metrics")
	rootCmd.PersistentFlags().StringVar(&metricsPath, "web.metrics.path", "/metrics", "Path under which to expose metrics")
	rootCmd.PersistentFlags().StringVar(&sentinalPath, "sentinal.path", "/var/run/reboot-required", "Path to file to be monitored")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
