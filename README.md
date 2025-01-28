# Prometheus - Reboot exporter

- [Prometheus - Reboot exporter](#prometheus---reboot-exporter)
  - [About](#about)
  - [Getting started](#getting-started)
    - [Build](#build)
    - [Running exporter](#running-exporter)
  - [Flags](#flags)
  - [Metrics](#metrics)

## About
This exporter provides metrics from Debian-based systems by monitoring a specific file to determine if a reboot is required.

I created the project to gain some experience with go and to build a dashboard to easily manage pending reboots.

## Getting started

### Build
```bash
make
```

### Running exporter
```bash
bin/./reboot_exporter
```

## Flags
```bash
Usage of bin/./reboot_exporter:
  -web.telemetry-path string
        Path to expose metrics (default "/metrics")
  -web.telemetry-port string
        Port to listen for telemetry. (default "9001")
```

## Metrics
| Metric | Description |
| --- | --- |
| reboot_required | 1 = reboot required, 0 = reboot not required |