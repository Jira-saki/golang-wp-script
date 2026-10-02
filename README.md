# 🐉 golang-wp-script

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-Linux%20amd64-lightgrey)
![License](https://img.shields.io/badge/license-MIT-blue)

A lightweight Go-based security toolkit for auditing WordPress and static-hosted sites in shared hosting environments. The project compiles into small native Linux binaries for fast local scanning and remote deployment without requiring a Go runtime on the target machine.

This repository focuses on practical web-host security checks: suspicious PHP discovery, plugin inventory extraction, upload-path anomaly detection, and multi-site scanning workflows.

---

## Overview

The tools in this repo are designed for environments like Xserver-style shared hosting, where administrators need quick, low-overhead scanners that run as standalone binaries instead of resource-heavy PHP-based monitoring systems.

Core goals:

- Detect unauthorized PHP files hiding in asset directories
- Identify suspicious files in WordPress root or upload paths
- Spot non-WordPress static-site anomalies
- Inventory active WordPress plugins from plugin header metadata
- Support fleet-wide scanning across many hosted domains

---

## Prerequisites

Before building or deploying, make sure you have:

| Requirement | Version | Notes |
| --- | --- | --- |
| Go | 1.24+ | Used for local compilation |
| SSH access | — | Needed to copy binaries to the server |
| Linux target host | amd64 | Typical shared-hosting architecture |
| `scp` or `rsync` | — | Used for transfer |

> No Go runtime is required on the production target server. The compiled binaries are self-contained.

---

## Repository layout

```text
golang-wp-script/
├── bin/                 # Compiled Linux binaries
│   ├── cli-template     # Starter CLI template
│   ├── prober           # HTTP endpoint prober
│   ├── walker           # WordPress/static security scanner
│   └── xserver-walker    # Production scanner binary
├── cmd/                 # Go entry points
│   ├── cli-template/
│   ├── prober/
│   └── walker/
├── python-script/       # Supporting Python probes
├── scanwp               # Multi-site sweep helper
├── go.mod               # Go module definition
├── README.md            # Project documentation
└── .gitignore           # Repo ignores
```

---

## Build and cross-compile

Build locally on macOS or another dev machine, then deploy the native binary to Linux:

```bash
# Build the scanner
GOOS=linux GOARCH=amd64 go build -o bin/walker ./cmd/walker

# Build the prober
GOOS=linux GOARCH=amd64 go build -o bin/prober ./cmd/prober

# Optional: build the template CLI
GOOS=linux GOARCH=amd64 go build -o bin/cli-template ./cmd/cli-template
```

Deploy with SSH:

```bash
scp bin/walker user@host:~
scp bin/prober user@host:~
```

---

## Scanner behavior

The main scanner in `cmd/walker` checks for:

- PHP files inside asset-like directories such as `/favicons/`, `/ogp/`, `/images/`, `/css/`, and `/js/`
- Unexpected PHP files on static sites
- PHP files in WordPress upload directories or root-level paths that look suspicious
- Plugin metadata in `wp-content/plugins/*` for plugin name and version reporting
- Safe `index.php` guard patterns that should not trigger alerts

### Example output

```text
🔍 [Golang Scanner v2.4] Starting security sweep in: ./public_html
--------------------------
📦 [PLUGIN DETECTED] All in One SEO (v5.0.2.1)
📦 [PLUGIN DETECTED] Contact Form 7 (v6.1.7)
📦 [PLUGIN DETECTED] Wordfence Security (v9.0.2)
🚨 [CRITICAL] PHP file hiding in asset directory: public_html/favicons/favicons.php
🚨 [STATIC ANOMALY] Unexpected PHP file on static site: public_html/tsigcfnd.php
--------------------------
✅ [Golang Scanner] Scan completed.
```

---

## Multi-tenant sweep helper

The repository includes a shell helper named `scanwp` for sweeping many tenant directories in one pass.

```bash
scanwp() {
  echo "🚀 [DevSecOps] Starting Global Security Scan..."
  echo "========================================"
  cd ~ || return
  for dir in */; do
    if [ -d "${dir}public_html" ]; then
      echo "🏢 [SCANNING] Domain: ${dir%/}"
      cd "$dir" && ../xserver-walker && cd ..
      echo "========================================"
    fi
  done
}
```

This is useful when scanning a shared-host account containing many sites.

---

## Security filtration rules

The scanner intentionally reduces false positives by ignoring known-good WordPress files and metadata. Examples include:

- Gutenberg asset files like `*.asset.php`
- WordPress core assets under `wp-includes/`
- Standard plugin and theme files under `wp-content/plugins/` and `wp-content/themes/`
- Safe guard-style `index.php` files that contain minimal or empty content

This makes the tool more suitable for real production hosting environments where noise must be controlled.

---

## Notes

This project is intended for security auditing and incident-response workflows. It is not a substitute for a full web application firewall or a managed hosting defense platform, but it can be a strong lightweight companion for identifying suspicious behavior and inventorying WordPress assets quickly.

---

## License

MIT
