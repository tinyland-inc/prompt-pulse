# prompt-pulse

Canonical writable source for the Go `prompt-pulse` application.

This repo owns the Go application source. `tinyland-inc/lab` owns pinned
consumption, Home Manager integration, shell-start behavior, and fleet/runtime
validation.

## Authority

- Go app source authority: this repo
- Rust TUI authority: `Jesssullivan/prompt-pulse-tui`
- Fleet/runtime integration authority: `tinyland-inc/lab`
- Waifu backend authority: `tinyland-inc/waifu-mirror`

## Current Module-Path Note

The canonical repo home is GitHub:

- repo: `https://github.com/tinyland-inc/prompt-pulse`

The Go module path is still the legacy path:

- module: `gitlab.com/tinyland/lab/prompt-pulse`

That module path is intentionally unchanged in this tranche. Until a dedicated
module-path migration lands, `go install` and internal imports continue to use
the legacy module path even though the canonical writable repo lives on GitHub.

## Freeze Boundary

Until a deliberate module-path migration is approved, keep these boundaries
explicit:

- repo home and source authority: GitHub (`tinyland-inc/prompt-pulse`)
- Go module path and internal imports: `gitlab.com/tinyland/lab/prompt-pulse`
- Bazel/Gazelle prefix: must match the legacy module path, not the GitHub repo
  home
- historical sync helpers: may still exist for compatibility, but they do not
  make GitLab the canonical writable repo

## Contribution Flow

1. Edit source here.
2. Validate source here.
3. Test integration from `lab`.
4. Repin `lab` to the landed upstream revision.

Do not treat `lab/cmd/prompt-pulse` as the normal day-to-day writable source.

## Local Remote Policy

- `origin` = canonical GitHub repo
- explicit extra remotes only, such as `mirror`, `personal`, or `archive`

## More Docs

- usage and configuration: [docs/README.md](docs/README.md)
- full guide: [docs/PROMPT_PULSE_GUIDE.md](docs/PROMPT_PULSE_GUIDE.md)
