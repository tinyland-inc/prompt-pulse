# Prompt-Pulse User Guide

## Overview

prompt-pulse is a shell-integrated infrastructure monitoring dashboard written in
Go. It consolidates status information from multiple sources into a single
glanceable view, available as Starship prompt segments, a system login banner, or
a full interactive TUI via the separate `prompt-pulse-tui` binary.

```
+------------------------------------------------------------------+
| prompt-pulse                                                      |
|                                                                   |
|  +-----------+-----------+-----------+                            |
|  | [Claude]  |  Billing  |  Infra    |  <-- tab bar               |
|  +-----------+-----------+-----------+                            |
|                                                                   |
|  Claude AI Usage                                                  |
|                                                                   |
|  personal  [SUB] pro  OK                                          |
|  5h usage  [========------] 62%     Resets in 43m                 |
|  7d usage  [====----------] 31%     Resets at 18:00               |
|  ------                                                           |
|  work      [API] tier_2  OK                                       |
|  Requests  [====---------] 34%     800 / 2000 used               |
|  Tokens    [==-----------] 18%     9k / 50k used                  |
|                                                                   |
|  q: quit | tab: switch | 1-3: jump    Updated: 14:32:07          |
+------------------------------------------------------------------+
```

### What It Collects

- **Claude Admin Usage** -- Usage across one or more Anthropic orgs via admin
  keys, including current and previous month token and cost summaries.
- **Cloud Billing** -- Current-month spend and resource summaries from Civo and
  DigitalOcean.
- **Infrastructure** -- Tailscale mesh node connectivity and Kubernetes cluster
  health (node readiness, CPU, memory, pod counts).
- **System Status** -- Automatic health evaluation with configurable thresholds
  that drive waifu banner mood selection.

### Display Modes

| Mode | Flag | Description |
|------|------|-------------|
| Single pass | *(none)* | Collect once, write to cache, exit |
| Daemon | `--daemon` | Background polling loop with PID file |
| Interactive TUI | `prompt-pulse-tui` | Separate Rust TUI for dashboard navigation |
| Banner | `--banner` | Terminal login banner with optional waifu image |
| Starship | `--starship <module>` | One-line output for Starship prompt segment |

---

## Quick Start

### 1. Install

```bash
# Using Nix (recommended)
nix profile install .#prompt-pulse

# Or build from source
go build -o prompt-pulse .
```

### 2. Create a minimal configuration

```bash
mkdir -p ~/.config/prompt-pulse

cat > ~/.config/prompt-pulse/config.toml << 'EOF'
[collectors.claude]
enabled = true
interval = "5m"

[[collectors.claude.account]]
name = "personal"
EOF
```

### 3. Run a single collection pass

```bash
prompt-pulse
```

This fetches data from all configured providers, writes results to the cache at
`~/.cache/prompt-pulse/`, and exits.

### 4. Check status from the command line

```bash
prompt-pulse --starship claude
# Output: personal:pro 62%

prompt-pulse --starship billing
# Output: $142 ($180 forecast)

prompt-pulse --starship infra
# Output: ts:5/7 k8s:civo:healthy
```

### 5. Launch the TUI

```bash
prompt-pulse-tui
```

### 6. Add to your shell

Source the shell integration in your RC file to get daemon/status/banner
helpers plus the `Ctrl+P` keybinding that launches `prompt-pulse-tui` when it
is installed. See the [Shell Integration](#shell-integration) section for
details.

---

## Installation

### Nix (Recommended)

prompt-pulse is packaged as a Nix flake output. If your system uses the
crush-dots flake:

```bash
# Install to user profile
nix profile install .#prompt-pulse

# Or add to your home-manager packages
# (see nix/home-manager/prompt-pulse.nix)
```

### Go Toolchain Install Note

Requires Go 1.23 or later:

```bash
# Public `go install ...@latest` is not currently published for this repo.
# Use a GitHub checkout and build locally until the module-path migration lands.
```

### From Source

```bash
git clone https://github.com/tinyland-inc/prompt-pulse.git
cd prompt-pulse
go mod download
go build -o prompt-pulse .
```

Build with version metadata:

```bash
go build -ldflags "\
  -X main.version=0.1.0 \
  -X main.commit=$(git rev-parse --short HEAD) \
  -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o prompt-pulse .
```

---

## Configuration

### Config File Location

The default configuration file is:

```
~/.config/prompt-pulse/config.toml
```

Override with the `--config` flag:

```bash
prompt-pulse --config /path/to/custom-config.toml
```

If the config file does not exist, prompt-pulse uses built-in defaults and
continues without error.

### Full Configuration Reference

prompt-pulse v2 uses TOML, not YAML. A full repo-backed example lives at
`pkg/config/testdata/full.toml`.

Minimal structure:

```toml
[general]
daemon_poll_interval = "15m"

[collectors.claude]
enabled = true
interval = "5m"

[[collectors.claude.account]]
name = "personal"

[shell]
tui_keybinding = "\\C-p"
show_banner_on_startup = true
```

The runtime search path is:

1. `$XDG_CONFIG_HOME/prompt-pulse/config.toml`
2. `~/.config/prompt-pulse/config.toml`

### Example Configurations

#### Minimal (Claude only)

```toml
[collectors.claude]
enabled = true
interval = "5m"

[[collectors.claude.account]]
name = "personal"
organization_id = ""
```

Set `ANTHROPIC_ADMIN_KEY` or `ANTHROPIC_ADMIN_KEY_FILE` in the environment.

#### Full (current repo surface)

```toml
[general]
daemon_poll_interval = "10m"

[collectors.claude]
enabled = true
interval = "10m"

[[collectors.claude.account]]
name = "personal"
organization_id = ""

[[collectors.claude.account]]
name = "work"
organization_id = ""

[collectors.billing]
enabled = true
interval = "20m"

[collectors.billing.civo]
enabled = true
region = "nyc1"

[collectors.billing.digitalocean]
enabled = true

[collectors.tailscale]
enabled = true
interval = "30s"

[collectors.kubernetes]
enabled = true
interval = "90s"
contexts = ["civo-cluster", "homelab-rke2"]
namespaces = ["default", "fuzzy-dev"]

[image]
waifu_enabled = true
waifu_category = "waifu"

[theme]
name = "catppuccin"

[shell]
tui_keybinding = "\\C-p"
show_banner_on_startup = true
instant_banner = true
```

#### Custom Thresholds (status evaluation only)

Status evaluation thresholds are currently set at compile time via
`status.DefaultEvaluatorConfig()`. The defaults are:

| Threshold | Default | Description |
|-----------|---------|-------------|
| Claude warning | 80% | 5-hour utilization triggers warning |
| Claude critical | 95% | 5-hour utilization triggers critical |
| Billing budget warning | 80% | % of budget triggers warning |
| Billing budget critical | 100% | % of budget triggers critical |
| Tailscale warning | 50% | % of nodes online triggers warning |
| K8s node ready minimum | 80% | % of nodes ready triggers warning |

---

## CLI Usage

### Flags

```
prompt-pulse [flags]

Flags:
  -banner           Display system status banner
  -daemon           Run background polling daemon
  -starship string  Output one-line Starship format (claude|billing|infra)
  -config string    Path to configuration file
                    (default: ~/.config/prompt-pulse/config.toml)
  -shell string     Output shell integration script (bash|zsh|fish|ksh)
  -health           Check daemon health status
  -diagnose         Claude diagnostics
  -migrate          Run v1-to-v2 config migration
  -verbose          Enable verbose (debug-level) logging
  -version          Print version and exit
```

### Default Mode (no flags)

When invoked with no flags, prompt-pulse runs a **single collection pass**:
collects data from all configured providers, writes results to the cache
directory, and exits. This is useful for cron jobs or manual cache refreshes.

```bash
prompt-pulse                         # collect once with default config
prompt-pulse --config /etc/pp.yaml   # collect once with custom config
prompt-pulse --verbose               # collect with debug logging
```

### Version

```bash
prompt-pulse --version
# Output: prompt-pulse 0.1.0 (abc1234) built 2026-01-15T12:00:00Z
```

### Shell Functions (installed via shell integration)

When the shell integration is sourced, prompt-pulse generates daemon/status and
banner helpers plus keybindings that launch `prompt-pulse-tui` when it is
installed.

- Bash uses underscored helper names: `pp_start`, `pp_stop`, `pp_status`,
  `pp_banner`
- Zsh, Fish, and Ksh use dashed helper names: `pp-start`, `pp-stop`,
  `pp-status`, `pp-banner`
- `Ctrl+P` launches `prompt-pulse-tui`
- when the shell generator enables the waifu keybinding, `Ctrl+W` launches
  `prompt-pulse-tui --expand waifu`

---

## Daemon Mode

The daemon runs a continuous background polling loop that keeps the cache
fresh for Starship segments and banner generation.

### Starting the Daemon

```bash
# Directly
prompt-pulse --daemon

# Or via the generated shell helper
# Bash: pp_start
# Zsh/Fish/Ksh: pp-start
```

### PID File Management

The daemon writes a PID file to `{cache_dir}/prompt-pulse.pid` (default:
`~/.cache/prompt-pulse/prompt-pulse.pid`).

- On startup, it checks if another instance is already running by reading the
  PID file and sending signal 0 to the recorded process.
- If the PID file contains a stale PID (process no longer exists), the file is
  cleaned up automatically.
- On shutdown, the PID file is removed.
- If another daemon is already running, startup fails with an error message
  including the existing PID.

### Polling Intervals

The daemon uses two levels of interval control:

1. **Global poll interval** (`daemon.poll_interval`): The ticker that drives
   collection passes. Default: 15 minutes.

2. **Per-collector intervals**: Each collector declares its own recommended
   interval via the `Interval()` method. If a collector ran more recently than
   its declared interval, it is skipped during that pass. This prevents
   rate-limited APIs from being hit too frequently even with a short global
   interval.

### Concurrent Collection

During each pass, all eligible collectors run concurrently in separate
goroutines. Results are written to the cache independently, so a slow provider
does not block others.

### Graceful Shutdown

The daemon handles `SIGINT` and `SIGTERM` signals:

1. Receives signal
2. Logs "received shutdown signal"
3. Cancels the context (stops in-flight collectors)
4. Logs final cache state (key names and ages)
5. Removes the PID file
6. Exits cleanly

```bash
# Stop via signal
kill $(cat ~/.cache/prompt-pulse/prompt-pulse.pid)

# Or via the generated shell helper
# Bash: pp_stop
# Zsh/Fish/Ksh: pp-stop
```

---

## TUI Dashboard

The interactive dashboard no longer lives behind a `prompt-pulse --tui` flag.
It is provided by the separate Rust binary `prompt-pulse-tui`.

```bash
prompt-pulse-tui
```

The default shell integration binds `Ctrl+P` to launch it, and the waifu
shortcut path uses `prompt-pulse-tui --expand waifu`.

Source authority for the TUI is the separate repo:

- `https://github.com/Jesssullivan/prompt-pulse-tui`

---

## Starship Integration

prompt-pulse integrates with [Starship](https://starship.rs/) via custom
modules that call the binary and display one-line status in your prompt.

### Setup

**Step 1**: Generate the Starship configuration snippet:

The generated TOML defines three custom modules. Add them to
`~/.config/starship.toml`:

```toml
# prompt-pulse Starship custom modules

[custom.pp_claude]
command = "prompt-pulse --starship claude"
when = "command -v prompt-pulse"
format = "[$symbol($output)]($style) "
symbol = ""
style = "purple"
shell = ["bash", "--noprofile", "--norc"]

[custom.pp_billing]
command = "prompt-pulse --starship billing"
when = "command -v prompt-pulse"
format = "[$symbol($output)]($style) "
symbol = "$"
style = "green"
shell = ["bash", "--noprofile", "--norc"]

[custom.pp_infra]
command = "prompt-pulse --starship infra"
when = "command -v prompt-pulse"
format = "[$symbol($output)]($style) "
symbol = ""
style = "cyan"
shell = ["bash", "--noprofile", "--norc"]
```

**Step 2**: Add the module references to your Starship `format` string:

```toml
format = """
$directory\
$git_branch\
$git_status\
${custom.pp_claude}\
${custom.pp_billing}\
${custom.pp_infra}\
$character"""
```

### Available Modules

#### pp_claude

Displays per-account usage in the format:

```
account_name:tier NN% | account_name:tier NN%
```

Examples:
```
personal:pro 62% | work:tier_2 34%
personal:pro 62%
work-api:ERR
```

For current builds, the Claude segment summarizes admin-key-backed account
usage. Accounts with errors show `:ERR`.

#### pp_billing

Displays aggregate spending:

```
$142 ($180 forecast)
$85
$200 ($250 forecast) OVER BUDGET
```

#### pp_infra

Displays infrastructure summary:

```
ts:5/7 k8s:civo:healthy
ts:3/3
k8s:homelab:degraded
```

### Staleness Indicator

When cached data is older than the configured `poll_interval` (TTL), Starship
output is suffixed with ` ?` to indicate staleness:

```
personal:pro 62% ?
```

This tells you the daemon may not be running or the last collection failed.

---

## Shell Integration

prompt-pulse generates shell integration scripts that provide keybindings and
convenience functions. Four shells are supported.

### Bash

Add to `~/.bashrc`:

```bash
# prompt-pulse integration
eval "$(prompt-pulse shell bash)"
# Or source a pre-generated file:
# source ~/.config/prompt-pulse/bash-integration.sh
```

The integration provides:
- `Ctrl+P` keybinding to launch `prompt-pulse-tui` (via `bind -x`)
- `pp_start`, `pp_stop`, `pp_status`, `pp_banner`

### Zsh

Add to `~/.zshrc`:

```zsh
# prompt-pulse integration
eval "$(prompt-pulse shell zsh)"
```

The integration provides:
- `Ctrl+P` keybinding to launch `prompt-pulse-tui` (via `zle` widget + `bindkey`)
- `pp-start`, `pp-stop`, `pp-status`, `pp-banner`
- Tab completion for `prompt-pulse` flags

### Fish

Add to `~/.config/fish/config.fish`:

```fish
# prompt-pulse integration
prompt-pulse shell fish | source
```

The integration provides:
- `Ctrl+P` keybinding to launch `prompt-pulse-tui` (via `bind`)
- `pp-start`, `pp-stop`, `pp-status`, `pp-banner`
- Tab completions for all flags including `--starship` subcommands

Nushell generation is not currently implemented by `prompt-pulse shell`.

---

## Waifu Banner

The banner mode displays a system status summary suitable for terminal login
(e.g., in `.bashrc` or `.zshrc`). Optionally, it includes an anime-style image
from the configured waifu mirror endpoint, rendered inline using terminal image
protocols.

```bash
prompt-pulse --banner
```

### Enabling

Enable the waifu collector and image display in your TOML config:

```toml
[collectors.waifu]
enabled = true
interval = "1h"
endpoint = "https://waifu.example.internal"
category = "waifu"
cache_dir = "/tmp/prompt-pulse/waifu"
max_images = 64

[image]
waifu_enabled = true
waifu_category = "waifu"
```

### Banner Layout

When a waifu image is available, the banner renders side-by-side:

```
[image]  | hostname :: healthy
[image]  | ────────────────────────
[image]  |
[image]  | Claude
[image]  |   personal: 62% (5h) | 31% (7d)
[image]  |   work: 800/2000 req
[image]  |
[image]  | Billing
[image]  |   $142 this month ($180 forecast)
[image]  |
[image]  | Infrastructure
[image]  |   ts: 5/7 online
[image]  |   k8s: civo (healthy)
[image]  |
[image]  |   uptime: 14d 6h 32m
```

Without an image, the info panel renders full-width.

### Category Selection

`collectors.waifu.category` selects the category requested from the configured
mirror endpoint. `image.waifu_category` controls the display-side default used
by the banner and TUI. The shipped defaults use `"waifu"` for both values.

### Cache Management

Images are cached in the configured waifu cache directory. If
`collectors.waifu.cache_dir` is unset, the daemon stores images under
`{general.cache_dir}/waifu/`. Cache behavior:

- **Collector retention**: `collectors.waifu.max_images` bounds the number of
  downloaded images retained in the collector cache.
- **Display cache**: `[image]` settings control image-protocol session caching
  and on-disk image cache size for rendered banner output.
- **Prefetching**: The daemon prefetches images in the background so banner and
  TUI image rendering can stay local and fast.
- **Atomic writes**: Images are written via temp file + rename to prevent
  corruption from concurrent access.

---

## Status Evaluation

prompt-pulse evaluates system health by analyzing cached collector data against
configurable thresholds. The result drives waifu category selection and banner
status display.

### Health Levels

| Level | Severity | Description |
|-------|----------|-------------|
| `healthy` | 0 (best) | All systems normal |
| `unknown` | 1 | Insufficient data to evaluate |
| `warning` | 2 | Something needs attention |
| `critical` | 3 (worst) | Immediate attention required |

The **overall** status is the worst (highest severity) across all components.

### Thresholds

Default thresholds used for evaluation:

| Component | Condition | Level |
|-----------|-----------|-------|
| Claude | 5h utilization > 80% | Warning |
| Claude | 5h utilization > 95% | Critical |
| Claude | Account status not "ok" | Warning |
| Billing | Spend > 80% of budget | Warning |
| Billing | Spend > 100% of budget | Critical |
| Billing | Forecast exceeds budget | Warning |
| Billing | Provider status "error" | Warning |
| Tailscale | < 50% nodes online | Warning |
| Tailscale | 0 nodes online | Critical |
| Kubernetes | Cluster status "degraded" | Warning |
| Kubernetes | Cluster status "offline" | Critical |
| Kubernetes | < 80% nodes ready | Warning |

### Component Rules

Each component (claude, billing, infra) is evaluated independently. The worst
result from any sub-check becomes that component's level. The overall system
level is the worst across all three components.

---

## Circuit Breaker

Collectors are optionally wrapped with a circuit breaker that handles persistent
API failures gracefully.

### How It Works

The circuit breaker has three states:

```
            success
  +------+  ------>  +--------+
  | Open |           | Closed |  <-- normal operation
  +------+  <------  +--------+
      |     failure      |
      |   (>= max)       | failure
      |                   | (< max)
      v                   v
  +----------+       +--------+
  | wait...  |       | Closed |
  +----------+       +--------+
      |
      | timeout elapsed
      v
  +-----------+   success    +--------+
  | Half-Open | ----------> | Closed |
  +-----------+              +--------+
      |
      | failure (re-open with backoff)
      v
  +------+
  | Open |
  +------+
```

1. **Closed** (normal): Requests pass through. Failures are counted.
2. **Open** (tripped): After `max_failures` consecutive failures, the circuit
   opens. Requests are blocked and return a synthetic warning. The collector is
   skipped entirely, saving API calls and reducing log noise.
3. **Half-Open** (probe): After `reset_timeout` elapses, one request is allowed
   through as a probe. If it succeeds, the circuit closes. If it fails, the
   circuit re-opens with an exponentially increased timeout.

### Configuration

Default circuit breaker settings:

| Setting | Default | Description |
|---------|---------|-------------|
| `max_failures` | 3 | Consecutive failures before opening |
| `reset_timeout` | 1 minute | Initial wait before half-open probe |
| `max_reset_timeout` | 30 minutes | Cap on exponential backoff |
| `backoff_multiplier` | 2.0 | Timeout multiplier on each re-open |

### Monitoring

Circuit breaker statistics are available programmatically via `Stats()`:

- `State`: Current state (closed/open/half_open)
- `ConsecutiveFails`: Current failure streak
- `TotalFailures`: Lifetime failure count
- `TotalSuccesses`: Lifetime success count
- `LastFailure` / `LastSuccess`: Timestamps
- `CurrentTimeout`: Current reset timeout (includes backoff)
- `ConsecutiveSkips`: Number of collection passes skipped while open

---

## Provider Setup

### Claude Accounts

prompt-pulse currently tracks Anthropic org usage through Admin API keys.

```toml
[collectors.claude]
enabled = true
interval = "5m"

[[collectors.claude.account]]
name = "personal"
organization_id = ""
```

Supported environment inputs:

```bash
export ANTHROPIC_ADMIN_KEY="sk-ant-admin01-..."
# or
export ANTHROPIC_ADMIN_KEY_FILE="$HOME/.config/sops-nix/secrets/api/anthropic_admin"
```

For multiple accounts, use `ANTHROPIC_ADMIN_KEYS_FILE` with one `name:key`
entry per line.

### Civo

```toml
[collectors.billing]
enabled = true

[collectors.billing.civo]
enabled = true
region = "nyc1"
```

Set the environment variable:

```bash
export CIVO_TOKEN="your-civo-api-key"
```

Get your API key from the [Civo Dashboard](https://dashboard.civo.com/security).

### DigitalOcean

```toml
[collectors.billing]
enabled = true

[collectors.billing.digitalocean]
enabled = true
```

```bash
export DIGITALOCEAN_TOKEN="your-do-token"
```

Create a token at [DigitalOcean API Settings](https://cloud.digitalocean.com/account/api/tokens).

### Tailscale

```toml
[collectors.tailscale]
enabled = true
interval = "30s"
```

The current collector uses the local Tailscale LocalAPI socket. There is no
repo-native `TAILSCALE_API_KEY` config field in v2.

### Kubernetes

```toml
[collectors.kubernetes]
enabled = true
interval = "90s"
contexts = ["civo-cluster"]
namespaces = ["fuzzy-dev"]
```

prompt-pulse queries the Kubernetes API for the configured local kube contexts.
Current v2 config uses context names and namespaces only.

---

## Troubleshooting

### Common Issues

#### "failed to load config" on startup

The config file likely has a TOML syntax error. Validate with Python 3.11+:

```bash
python3 -c "import pathlib, tomllib; tomllib.loads(pathlib.Path('$HOME/.config/prompt-pulse/config.toml').read_text())"
```

#### Starship modules show nothing

1. Verify the daemon is running or run a manual collection:
   ```bash
   prompt-pulse --verbose
   ```
2. Check that cache files exist:
   ```bash
   ls -la ~/.cache/prompt-pulse/
   ```
3. Test individual modules:
   ```bash
   prompt-pulse --starship claude
   prompt-pulse --starship billing
   prompt-pulse --starship infra
   ```

Empty output means no cached data. Check the log file for errors.

#### "daemon already running (PID XXXX)"

Another daemon instance is active. Either stop it first or check if the PID
file is stale:

```bash
# Check if the process actually exists
kill -0 $(cat ~/.cache/prompt-pulse/prompt-pulse.pid) 2>/dev/null
echo $?  # 0 = running, 1 = stale

# Remove stale PID file
rm ~/.cache/prompt-pulse/prompt-pulse.pid
```

#### Banner shows "(no data)" for all sections

The banner only reads from cache, never from the network. Run a collection
first:

```bash
prompt-pulse                 # single collection pass
prompt-pulse --banner        # now shows data
```

Or start the daemon for continuous updates.

#### Waifu image not showing in banner

1. Verify waifu is enabled in both `[collectors.waifu]` and `[image]`.
2. The daemon must prefetch images. Run the daemon for at least one cycle.
3. Check the image cache:
   ```bash
   ls -la ~/.cache/prompt-pulse/waifu/
   ```
4. Verify your terminal supports image protocols (kitty, iTerm2, sixel).

### Debug Logging

Enable verbose logging to see detailed collection information:

```bash
prompt-pulse --verbose --daemon
```

Or for a single pass:

```bash
prompt-pulse --verbose
```

Logs are written to both stderr and the log file
(`~/.local/log/prompt-pulse.log` by default).

### Cache Inspection

Cache files are human-readable JSON stored at `~/.cache/prompt-pulse/`:

```bash
# List cached entries
ls ~/.cache/prompt-pulse/*.json

# View Claude data
cat ~/.cache/prompt-pulse/claude.json | python3 -m json.tool

# Check cache freshness (file modification time)
stat ~/.cache/prompt-pulse/claude.json

# View waifu image cache stats
ls -lhS ~/.cache/prompt-pulse/waifu/

# Clear all cached data
rm ~/.cache/prompt-pulse/*.json
```

---

## Architecture

### Package Structure

```
.
+-- main.go                        Entry point, flag parsing, mode dispatch
+-- pkg/
|   +-- config/                    TOML config parsing, defaults, env overrides
|   +-- daemon/                    Daemon loop, health file, cache wiring
|   +-- banner/                    Inline banner rendering
|   +-- starship/                  Starship segment rendering
|   +-- shell/                     Shell helper generation and keybindings
|   +-- collectors/
|   |   +-- claude/                Anthropic Admin API usage collector
|   |   +-- billing/               Civo and DigitalOcean billing collector
|   |   +-- tailscale/             Local Tailscale LocalAPI collector
|   |   +-- k8s/                   Kubernetes context collector
|   |   +-- sysmetrics/            Local system metrics collector
|   |   +-- waifu/                 Waifu endpoint/cache collector
|   +-- docs/                      Generated docs and manpages
|   +-- theme/                     Theme registry
|   +-- image/                     Waifu/banner image rendering
|   |   +-- infra_tab.go           Infra tab renderer
|   +-- widgets/
|       +-- gauge.go               Gauge bar rendering
|       +-- sparkline.go           Sparkline chart rendering
|       +-- table.go               Table rendering
|       +-- status.go              Status indicator rendering
|
+-- shell/
|   +-- common.go                  Shell type enum, IntegrationConfig, dispatch
|   +-- bash.go                    Bash integration generator
|   +-- zsh.go                     Zsh integration generator
|   +-- fish.go                    Fish integration generator
|
+-- collectors/waifu/
|   +-- client.go                  Mirror API client
|   +-- waifu.go                   Collector wiring and fetch loop
|
+-- waifu/
    +-- cache.go                   Local image cache management
    +-- picker.go                  Cached image selection helpers
    +-- prefetch.go                Background image prefetch coordination
    +-- session.go                 Banner/TUI session image tracking
```

### Data Flow

```
                                  +----------------------+
                                  | Configured Waifu      |
                                  | Mirror Endpoint       |
                                  +----------+-----------+
                                             |
                                    collect  |
                                             v
+------------------+              +--------+---------+
| Anthropic Admin  |              | Image Cache      |
| Civo API         | --collect--> | ~/.cache/prompt- |
| DigitalOcean API |              | pulse/waifu/     |
| Tailscale Local  |              +------------------+
| Kubernetes API   | --collect--> | JSON Cache       |
| Local sysmetrics |              | ~/.cache/prompt- |
+------------------+              | ~/.cache/prompt- |
                                  | pulse/*.json     |
        daemon loop               +--------+---------+
        (every 15m)                        |
                                  read     |  read
                               +-----------+-----------+
                               |           |           |
                               v           v           v
                          +--------+  +---------+  +--------+
                          |Starship|  | Banner  |  |  TUI   |
                          |Modules |  |Generator|  |        |
                          +--------+  +---------+  +--------+
                               |           |           |
                               v           v           v
                          [prompt]   [terminal]   [alt screen]
```

The key design principle is **cache-mediated decoupling**: the daemon writes to
cache files, and all display modes (Starship, banner, TUI) read from cache. This
means display operations are fast (no network) and the daemon can run
independently.
