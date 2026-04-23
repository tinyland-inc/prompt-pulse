# Prompt-Pulse Secrets Index

Current secret and credential surfaces for `prompt-pulse` live in these docs:

- [SECRETS_SUMMARY.md](SECRETS_SUMMARY.md): current repo-truth boundary and supported inputs
- [SECRETS_QUICK_REFERENCE.md](SECRETS_QUICK_REFERENCE.md): daily commands and environment patterns
- [SECRETS_SETUP.md](SECRETS_SETUP.md): current v2 setup flow for app-level config and env injection

## Boundary

This repo documents the app-level contract:

- `~/.config/prompt-pulse/config.toml`
- supported environment variables and `*_FILE` secret paths
- repo-native validation commands

`tinyland-inc/lab` owns the fleet-specific integration layer:

- Home Manager generation of `config.toml`
- launchd/systemd wiring
- sops-nix / KeePass-backed secret injection
- shell-start and workstation policy

## Current Supported Inputs

- `ANTHROPIC_ADMIN_KEY` or `ANTHROPIC_ADMIN_KEY_FILE`
- `ANTHROPIC_ADMIN_KEYS_FILE` for multi-account `name:key` lines
- `CIVO_TOKEN` or `CIVO_API_KEY_FILE`
- `CIVO_REGION`
- `DIGITALOCEAN_TOKEN` or `DIGITALOCEAN_TOKEN_FILE`
- `PPULSE_PROTOCOL`
- `PPULSE_THEME`
- `PPULSE_LAYOUT`

## Current Commands

```bash
prompt-pulse --diagnose
prompt-pulse --health
prompt-pulse --banner
prompt-pulse-tui
```

## Archived Historical Notes

The earlier February 5, 2026 yoga-specific operator notebooks are preserved as
historical context under [archive/](archive/):

- [archive/SECRETS_SUMMARY_2026-02-05.md](archive/SECRETS_SUMMARY_2026-02-05.md)
- [archive/SECRETS_QUICK_REFERENCE_2026-02-05.md](archive/SECRETS_QUICK_REFERENCE_2026-02-05.md)
- [archive/SECRETS_SETUP_2026-02-05.md](archive/SECRETS_SETUP_2026-02-05.md)
- [archive/SECRETS_INDEX_2026-02-05.md](archive/SECRETS_INDEX_2026-02-05.md)
