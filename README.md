# 🐉 golang-wp-script

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)
![Platform](https://img.shields.io/badge/platform-Linux%20amd64-lightgrey)

A high-performance, lightweight DevSecOps toolkit written in Go — built for automating security audits, inventorying software assets, monitoring uptimes, and hunting malware across multi-tenant WordPress shared hosting environments (e.g., Xserver).

Instead of relying on resource-heavy PHP plugins, this repo compiles into standalone native binaries (~2.5MB–3.3MB) that execute at blistering speeds with a near-zero memory footprint.

---

## 📋 Prerequisites

Before building or deploying, make sure you have the following:

| Requirement | Version | Notes |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21+ | Required for local cross-compilation |
| SSH access | — | Required to deploy binaries to the hosting server |
| Linux target server | amd64 | Shared hosting architecture (e.g., Xserver) |
| `scp` or `rsync` | — | For transferring compiled binaries |

> **No Go runtime is needed on the target server.** Binaries are fully self-contained.

---

## 📂 Repository Structure

```text
golang-wp-script/
├── bin/                 # Compiled native Linux binaries (production-ready)
│   ├── xserver-walker   # Production-hardened scanner with asset mapper & filters
│   ├── prober           # HTTP status and endpoint health prober
│   └── cli-template     # Scaffold skeleton for new Go CLI tools
├── cmd/                 # Source code entry points
│   ├── walker/          # Malware path hunter & plugin version auditor
│   ├── prober/
│   └── cli-template/
└── README.md
```

---

## 🎯 Component Breakdown

| Binary | Size | Purpose | Core Technology |
|---|---|---|---|
| `walker` / `xserver-walker` | ~2.6 MB | Recursively scans `uploads/` directories for unauthorized `.php` backdoors | `filepath.WalkDir` with custom path suppression |
| `prober` | ~3.3 MB | Bulk HTTP probing to verify site uptimes and response health across domains | Native `net/http` client with connection pooling |
| `cli-template` | ~2.5 MB | Boilerplate scaffold for rapidly building new automation utilities | Standard Go I/O flags |

---

## 🛠️ Build & Cross-Compile

All tools are cross-compiled locally on your development machine (for example, macOS) and deployed as native Linux binaries.

```bash
# Cross-compile for Linux amd64 from macOS
GOOS=linux GOARCH=amd64 go build -o bin/xserver-walker ./cmd/walker
GOOS=linux GOARCH=amd64 go build -o bin/prober ./cmd/prober

# Deploy using your SSH aliases
scp bin/xserver-walker xserver:~
```

## 💻 Operational Workflows

1. **Dual-Action Auditing (`xserver-walker`)**

   Navigate into the target tenant directory and run the standalone tool. It will extract top-level plugins and enforce security posture policies.

```bash
cd /home/user/target-domain.com
../xserver-walker
```

Expected real-time output (example):

```text
🔍 [Golang Scanner] Starting security sweep in: ./public_html
--------------------------
📦 [PLUGIN DETECTED] All in One SEO (v4.9.9)
📦 [PLUGIN DETECTED] Contact Form 7 (v6.1.6)
📦 [PLUGIN DETECTED] Otter – Page Builder Blocks & Extensions for Gutenberg (v3.2.0)
📦 [PLUGIN DETECTED] Wordfence Security (v8.2.2)
⚠️  [SUSPICIOUS] PHP file hiding in uploads directory: public_html/wp-content/uploads/backdoor.php
--------------------------
✅ [Golang Scanner] Scan completed.
```

2. **Multi-Tenant Fleet Automation (`scanwp`)**

Add this macro to your server's `~/.bashrc` to sweep multiple hosted domains in sequence:

```bash
scanwp() {
  echo "🚀 [DevSecOps] Starting Global WordPress Security Scan..."
  echo "========================================"
  cd ~ || return
  for dir in */ ; do
    if [ -d "${dir}public_html" ]; then
      echo "🏢 [SCANNING] Domain: ${dir%/}"
      cd "$dir" && ../xserver-walker && cd ..
      echo "========================================"
    fi
  done
}
```

## 🛡️ Smart Noise Filtering & Engine Features (v2.2)

`xserver-walker` uses targeted heuristic filtering to eliminate noise in production multi-tenant environments:

- **Two-Tier Threat Classification:**
  - `⚠️ [SUSPICIOUS]`: Unflagged `.php` files placed inside `/wp-content/uploads/`.
  - `🚨 [HIGH RISK]`: Obfuscated `index.php` backdoors in upload paths containing control-flow jumps (`goto`) or execution payloads.
- **`isSafeIndexPHP()` Inspection:** Inspects `index.php` files inside uploads. Legitimate empty files or standard `Silence is golden.` headers pass cleanly, while encoded payloads trigger high-severity alerts.
- **Directory Path Precision:** Matches strict `/wp-content/uploads/` path strings rather than loose keyword matching, preventing false positives on legitimate plugin classes.
- **Known-Good Exclusions:** Automatically skips administrative firewall rulesets (for example, AIOS `uploads/aios/firewall-rules/`).

Example scanner output:

```text
🔍 [Golang Scanner] Starting security sweep in: ./public_html
--------------------------
📦 [PLUGIN DETECTED] All in One SEO (v5.0.0.1)
📦 [PLUGIN DETECTED] Wordfence Security (v9.0.0)
🚨 [HIGH RISK] Malicious index.php found in uploads: public_html/wp-content/uploads/2026/08/index.php
⚠️  [SUSPICIOUS] PHP file hiding in uploads directory: public_html/wp-content/uploads/test.php
--------------------------
✅ [Golang Scanner] Scan completed.
```

## 🚨 Incident Response Case Study

This toolkit was used during a live production remediation. The scanner identified a persistent Base64 payload injected into system configuration files. The attacker attempted to disguise tampering with a fake comment header and a spoofed process signature.

**Lead Engineer:** Jirasak

**Focus:** Infrastructure Automation, Fleet Protection, Systems Security

**License:** MIT
