# Prompt-Pulse Secrets Summary

## Current Repo Truth

`prompt-pulse` v2 is TOML-first and environment-driven.

- Config file: `~/.config/prompt-pulse/config.toml`
- Separate interactive UI: `prompt-pulse-tui`
- Repo-native diagnostics: `prompt-pulse --diagnose` and `prompt-pulse --health`
- Canonical source repo: `tinyland-inc/prompt-pulse`
- Fleet integration repo: `tinyland-inc/lab`

## Supported Secret Surfaces

### Claude

Current Claude usage collection is Admin API based.

- `ANTHROPIC_ADMIN_KEY`
- `ANTHROPIC_ADMIN_KEY_FILE`
- `ANTHROPIC_ADMIN_KEYS_FILE`

`ANTHROPIC_ADMIN_KEYS_FILE` should contain one `name:key` entry per line.

### Billing

Current billing support is:

- Civo via `CIVO_TOKEN` or `CIVO_API_KEY_FILE`
- DigitalOcean via `DIGITALOCEAN_TOKEN` or `DIGITALOCEAN_TOKEN_FILE`

### Other Environment Overrides

- `CIVO_REGION`
- `PPULSE_PROTOCOL`
- `PPULSE_THEME`
- `PPULSE_LAYOUT`

## Repo vs Integration Boundary

This repo owns:

- Go app source
- config schema
- environment override contract
- CLI / banner / starship / shell behavior

`tinyland-inc/lab` owns:

- generated `config.toml`
- Home Manager integration
- launchd/systemd service wiring
- fleet secret materialization to `*_FILE` paths

## Quick Verification

```bash
prompt-pulse --diagnose
prompt-pulse --health
env | grep -E '^(ANTHROPIC_ADMIN|CIVO|DIGITALOCEAN|PPULSE_)'
```

## Historical Archive

The earlier yoga-specific operator audit was moved to
[archive/SECRETS_SUMMARY_2026-02-05.md](archive/SECRETS_SUMMARY_2026-02-05.md)
so it no longer presents itself as the live repo contract.
