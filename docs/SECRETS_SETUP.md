# Prompt-Pulse Secrets Setup

## Scope

This guide covers the current app-level contract for `prompt-pulse` v2.

It does not try to document one specific workstation’s Nix, sops, or Home
Manager state. If you are on the Tinyland managed fleet, use
`tinyland-inc/lab` for that integration layer.

## 1. Create or Review Config

The live config path is:

```text
~/.config/prompt-pulse/config.toml
```

Minimal example:

```toml
[collectors.claude]
enabled = true
interval = "5m"

[[collectors.claude.account]]
name = "personal"
organization_id = ""

[collectors.billing]
enabled = true

[collectors.billing.civo]
enabled = true
region = "nyc1"

[collectors.billing.digitalocean]
enabled = true

[collectors.tailscale]
enabled = true

[collectors.kubernetes]
enabled = false
```

For a fuller example, start from `pkg/config/testdata/full.toml`.

## 2. Provide Secrets

### Claude Admin API

Use either a direct env var or a file-backed env var:

```bash
export ANTHROPIC_ADMIN_KEY="sk-ant-admin01-..."
# or
export ANTHROPIC_ADMIN_KEY_FILE="$HOME/.config/sops-nix/secrets/api/anthropic_admin"
```

For multiple accounts:

```bash
export ANTHROPIC_ADMIN_KEYS_FILE="$HOME/.config/sops-nix/secrets/api/anthropic_admin_keys"
```

Each line in that file should be:

```text
personal:sk-ant-admin01-...
work:sk-ant-admin01-...
```

### Civo

```bash
export CIVO_TOKEN="..."
# or
export CIVO_API_KEY_FILE="$HOME/.config/sops-nix/secrets/infrastructure/civo_api_key"
export CIVO_REGION="nyc1"
```

### DigitalOcean

```bash
export DIGITALOCEAN_TOKEN="..."
# or
export DIGITALOCEAN_TOKEN_FILE="$HOME/.config/sops-nix/secrets/infrastructure/digitalocean_token"
```

## 3. Optional Presentation Overrides

```bash
export PPULSE_PROTOCOL="kitty"
export PPULSE_THEME="catppuccin"
export PPULSE_LAYOUT="dashboard"
```

## 4. Validate

```bash
prompt-pulse --diagnose
prompt-pulse --health
prompt-pulse --banner
prompt-pulse-tui
```

If the daemon is not running yet, `prompt-pulse --health` will fail with a
not-running status. That is expected until your integration layer starts it.

## 5. Integration Boundary

If you are using Tinyland fleet automation:

- edit the generated/runtime integration in `tinyland-inc/lab`
- keep source/schema changes in this repo
- repin `lab` after upstream changes land here

## Historical Archive

The older workstation-specific operator notebook is preserved at
[archive/SECRETS_SETUP_2026-02-05.md](archive/SECRETS_SETUP_2026-02-05.md).
