# Prometheus - Reboot exporter

- [Prometheus - Reboot exporter](#prometheus---reboot-exporter)
  - [About](#about)
  - [Getting started](#getting-started)
    - [Build](#build)
  - [Flags](#flags)
  - [Metrics](#metrics)
  - [Sources](#sources)

## About
This exporter provides metrics from Debian-based systems by monitoring a specific file to determine if a reboot is required.

I created the project to gain some experience with go and to build a dashboard to easily manage pending reboots.

## Getting started

> [!CAUTION]
> The exporter works only on Debian-based systems

### Build
To build a binary, clone the repository and run:
```bash
make
```

Run the binary:
```bash
./bin/./reboot_exporter
```

## Flags
| Flag | Type | Description |
| --- | --- | --- |
| --web.telemetry-path | string | Path under which to expose metrics. (default "/metrics") |
| --web.telemetry-port | string | Port on which to expose metrics (default 11011) |

## Metrics
| Metric | Description |
| --- | --- |
| reboot_required | 1 = reboot required, 0 = reboot not required |

## Sources
Github: https://github.com/kryptonhydrit/prometheus-reboot-exporter