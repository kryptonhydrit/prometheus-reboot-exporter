package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
)

var (
	listenAddr  string
	metricsPath string
	watchedFile string
)

var rebootRequiredGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "reboot_required",
		Help: "Indicates if a reboot is required (1 = required, 0 = not required).",
	})

func startWatcher() error {
	if _, err := os.Stat(watchedFile); err == nil {
		rebootRequiredGauge.Set(1)
	} else {
		rebootRequiredGauge.Set(0)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error when creating the Watcher: %w", err)
	}

	watchDir := path.Dir(watchedFile)
	err = watcher.Add(watchDir)
	if err != nil {
		return fmt.Errorf("error when monitoring the directory: %w", err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Name == watchedFile {
					if event.Op&fsnotify.Create == fsnotify.Create {
						rebootRequiredGauge.Set(1)
					} else if event.Op&fsnotify.Remove == fsnotify.Remove {
						rebootRequiredGauge.Set(0)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Error: %s", err)
			}
		}
	}()
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:  "reboot_exporter",
		Long: "Export metrics about the reboot status from Debian-based systems",
		Run: func(cmd *cobra.Command, args []string) {
			registry := prometheus.NewRegistry()
			registry.MustRegister(rebootRequiredGauge)

			startWatcher()

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
	rootCmd.PersistentFlags().StringVar(&watchedFile, "monitored.path", "/var/run/reboot-required", "Path to file to be monitored")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
