# Prompt-Pulse Secrets Quick Reference

## Current Commands

```bash
# runtime diagnostics
prompt-pulse --diagnose

# daemon/cache health
prompt-pulse --health

# render the current banner from cache
prompt-pulse --banner

# open the separate interactive dashboard
prompt-pulse-tui
```

## Current Config Path

```text
~/.config/prompt-pulse/config.toml
```

Use `pkg/config/testdata/full.toml` for a full example.

## Current Environment Patterns

```bash
export ANTHROPIC_ADMIN_KEY="sk-ant-admin01-..."
export CIVO_TOKEN="..."
export DIGITALOCEAN_TOKEN="..."
```

File-backed secret patterns are also supported:

```bash
export ANTHROPIC_ADMIN_KEY_FILE="$HOME/.config/sops-nix/secrets/api/anthropic_admin"
export CIVO_API_KEY_FILE="$HOME/.config/sops-nix/secrets/infrastructure/civo_api_key"
export DIGITALOCEAN_TOKEN_FILE="$HOME/.config/sops-nix/secrets/infrastructure/digitalocean_token"
```

Multi-account Claude input:

```bash
export ANTHROPIC_ADMIN_KEYS_FILE="$HOME/.config/sops-nix/secrets/api/anthropic_admin_keys"
```

Each line should be `name:key`.

## Current Supported Keys

- `ANTHROPIC_ADMIN_KEY`
- `ANTHROPIC_ADMIN_KEY_FILE`
- `ANTHROPIC_ADMIN_KEYS_FILE`
- `CIVO_TOKEN`
- `CIVO_API_KEY_FILE`
- `CIVO_REGION`
- `DIGITALOCEAN_TOKEN`
- `DIGITALOCEAN_TOKEN_FILE`
- `PPULSE_PROTOCOL`
- `PPULSE_THEME`
- `PPULSE_LAYOUT`

## Troubleshooting

```bash
# inspect effective environment
env | grep -E '^(ANTHROPIC_ADMIN|CIVO|DIGITALOCEAN|PPULSE_)'

# validate TOML shape quickly
python3 -c "import pathlib, tomllib; tomllib.loads(pathlib.Path('$HOME/.config/prompt-pulse/config.toml').read_text())"

# check current repo build surface
go test ./...
```

If you are using Tinyland fleet automation, do not edit `lab` paths from memory.
Use the current `tinyland-inc/lab` operator docs for sops/Home Manager wiring.
