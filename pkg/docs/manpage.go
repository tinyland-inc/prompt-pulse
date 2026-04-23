package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ManPage holds the content for a Unix man page.
type ManPage struct {
	// Name is the command or file name.
	Name string

	// Section is the man page section number (e.g., "1", "5").
	Section string

	// ShortDesc is the one-line description for the NAME section.
	ShortDesc string

	// Synopsis is the command usage line.
	Synopsis string

	// Description is the full description body.
	Description string

	// Options lists command-line flags and their descriptions.
	Options string

	// Examples shows usage examples.
	Examples string

	// SeeAlso lists related man pages.
	SeeAlso string
}

// dcGenerateManPage builds a man page for the given command name and section.
func dcGenerateManPage(name, section string) *ManPage {
	switch name + "." + section {
	case "prompt-pulse.1":
		return dcManPromptPulse()
	case "prompt-pulse-daemon.1":
		return dcManDaemon()
	case "prompt-pulse-banner.1":
		return dcManBanner()
	case "prompt-pulse-tui.1":
		return dcManTUI()
	case "prompt-pulse.toml.5":
		return dcManConfig()
	default:
		return &ManPage{
			Name:      name,
			Section:   section,
			ShortDesc: "unknown command",
		}
	}
}

// dcAllManPages returns man pages for all documented commands.
func dcAllManPages() []*ManPage {
	return []*ManPage{
		dcManPromptPulse(),
		dcManDaemon(),
		dcManBanner(),
		dcManTUI(),
		dcManConfig(),
	}
}

// WriteManPages writes all man pages as roff files to the given base directory,
// creating man1/ and man5/ subdirectories as needed. Returns the number of pages written.
func WriteManPages(baseDir string) (int, error) {
	pages := dcAllManPages()
	for _, mp := range pages {
		subdir := filepath.Join(baseDir, "man"+mp.Section)
		if err := os.MkdirAll(subdir, 0o755); err != nil {
			return 0, fmt.Errorf("creating %s: %w", subdir, err)
		}
		path := filepath.Join(subdir, mp.Name+"."+mp.Section)
		content := dcRenderManRoff(mp)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return 0, fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return len(pages), nil
}

// dcRenderManRoff renders a ManPage in roff/troff format suitable for man(1).
func dcRenderManRoff(mp *ManPage) string {
	var b strings.Builder
	date := time.Now().Format("January 2006")

	// Header
	b.WriteString(fmt.Sprintf(".TH %s %s \"%s\" \"prompt-pulse v2.0.0\" \"prompt-pulse Manual\"\n",
		strings.ToUpper(mp.Name), mp.Section, date))

	// NAME
	b.WriteString(".SH NAME\n")
	b.WriteString(fmt.Sprintf("%s \\- %s\n", mp.Name, mp.ShortDesc))

	// SYNOPSIS
	if mp.Synopsis != "" {
		b.WriteString(".SH SYNOPSIS\n")
		b.WriteString(fmt.Sprintf(".B %s\n", mp.Synopsis))
	}

	// DESCRIPTION
	if mp.Description != "" {
		b.WriteString(".SH DESCRIPTION\n")
		b.WriteString(mp.Description + "\n")
	}

	// OPTIONS
	if mp.Options != "" {
		b.WriteString(".SH OPTIONS\n")
		b.WriteString(mp.Options + "\n")
	}

	// EXAMPLES
	if mp.Examples != "" {
		b.WriteString(".SH EXAMPLES\n")
		b.WriteString(mp.Examples + "\n")
	}

	// SEE ALSO
	if mp.SeeAlso != "" {
		b.WriteString(".SH SEE ALSO\n")
		b.WriteString(mp.SeeAlso + "\n")
	}

	return b.String()
}

// dcRenderManMarkdown renders a ManPage as a Markdown document for web docs.
func dcRenderManMarkdown(mp *ManPage) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s(%s)\n\n", mp.Name, mp.Section))

	b.WriteString("## NAME\n\n")
	b.WriteString(fmt.Sprintf("%s - %s\n\n", mp.Name, mp.ShortDesc))

	if mp.Synopsis != "" {
		b.WriteString("## SYNOPSIS\n\n")
		b.WriteString(fmt.Sprintf("```\n%s\n```\n\n", mp.Synopsis))
	}

	if mp.Description != "" {
		b.WriteString("## DESCRIPTION\n\n")
		b.WriteString(mp.Description + "\n\n")
	}

	if mp.Options != "" {
		b.WriteString("## OPTIONS\n\n")
		b.WriteString(mp.Options + "\n\n")
	}

	if mp.Examples != "" {
		b.WriteString("## EXAMPLES\n\n")
		b.WriteString(mp.Examples + "\n\n")
	}

	if mp.SeeAlso != "" {
		b.WriteString("## SEE ALSO\n\n")
		b.WriteString(mp.SeeAlso + "\n\n")
	}

	return b.String()
}

func dcManPromptPulse() *ManPage {
	return &ManPage{
		Name:      "prompt-pulse",
		Section:   "1",
		ShortDesc: "terminal dashboard with live data, waifu rendering, and shell surfaces",
		Synopsis:  "prompt-pulse [flags]",
		Description: `prompt-pulse is a terminal dashboard that displays system metrics, cloud billing,
Tailscale network status, Kubernetes cluster health, and Claude API usage.
It supports multiple display modes: banner (inline), starship output, shell integration,
health checks, and daemon-backed cache updates. The interactive dashboard is provided by
the separate prompt-pulse-tui binary.

prompt-pulse renders waifu images using the best available terminal protocol
(Kitty Unicode placeholders, iTerm2 inline images, Sixel, or half-blocks)
and adapts its layout to terminal width.`,
		Options: `.TP
.B \-\-banner
Display an inline banner with current status data.
.TP
.B \-\-daemon
Run the background data collection daemon.
.TP
.B \-\-starship <segment>
Render a starship segment: claude, billing, infra, k8s, system, or all.
.TP
.B \-\-shell <shell>
Print shell integration for bash, zsh, fish, or ksh.
.TP
.B \-\-health
Show daemon health and collected data summary.
.TP
.B \-\-diagnose
Print runtime diagnostics and config search paths.
.TP
.B \-\-migrate
Migrate a legacy v1 config.yaml file to v2 config.toml.
.TP
.B \-\-config <path>
Override the config file path.
.TP
.B \-\-theme <name>
Override the color theme (default, gruvbox, nord, catppuccin, dracula, tokyo-night).
.TP
.B \-\-version
Print version and exit.`,
		Examples: `.nf
# Show banner
prompt-pulse --banner

# Start daemon in background
prompt-pulse --daemon &

# Check daemon health
prompt-pulse --health

# Initialize shell integration
eval "$(prompt-pulse --shell bash)"
.fi`,
		SeeAlso: `.BR prompt-pulse-daemon (1),
.BR prompt-pulse-banner (1),
.BR prompt-pulse-tui (1),
.BR prompt-pulse.toml (5)`,
	}
}

func dcManDaemon() *ManPage {
	return &ManPage{
		Name:      "prompt-pulse-daemon",
		Section:   "1",
		ShortDesc: "prompt-pulse background data collection mode",
		Synopsis:  "prompt-pulse --daemon",
		Description: `The prompt-pulse daemon runs in the background, collecting data from configured
sources at regular intervals. The CLI health surface is exposed through
prompt-pulse --health.

The daemon caches collected data so that banner and TUI modes can display
information instantly without waiting for API calls.`,
		Options: `.TP
.B \-\-daemon
Start the daemon in the foreground.
.TP
.B \-\-health
Print daemon health: PID, uptime, collector status, and data freshness.
.TP
.B \-\-config <path>
Override the config file path used by the daemon.`,
		Examples: `.nf
# Start daemon
prompt-pulse --daemon

# Check health
prompt-pulse --health

# Start with explicit config
prompt-pulse --daemon --config /path/to/config.toml
.fi`,
		SeeAlso: `.BR prompt-pulse (1),
.BR prompt-pulse.toml (5)`,
	}
}

func dcManBanner() *ManPage {
	return &ManPage{
		Name:      "prompt-pulse-banner",
		Section:   "1",
		ShortDesc: "prompt-pulse inline terminal banner",
		Synopsis:  "prompt-pulse --banner [options]",
		Description: `The banner mode displays a compact inline summary of system status in the terminal.
It adapts to terminal width with four layout modes: compact (<80 cols),
standard (120+ cols), wide (160+ cols), and ultra-wide (200+ cols).

The banner reads data from the daemon's cache for instant display.
If the daemon is not running, it will attempt a quick data fetch with a timeout.`,
		Options: `.TP
.B \-\-no-waifu
Disable waifu image in the banner.
.TP
.B \-\-width <cols>
Override terminal width detection.
.TP
.B \-\-instant
Use pre-rendered cached banner (default: true).
.TP
.B \-\-timeout <duration>
Maximum time to wait for daemon data.`,
		Examples: `.nf
# Show banner
prompt-pulse --banner

# Banner without waifu
prompt-pulse --banner --no-waifu

# Force compact mode
prompt-pulse --banner --width 79
.fi`,
		SeeAlso: `.BR prompt-pulse (1),
.BR prompt-pulse-tui (1),
.BR prompt-pulse.toml (5)`,
	}
}

func dcManTUI() *ManPage {
	return &ManPage{
		Name:      "prompt-pulse-tui",
		Section:   "1",
		ShortDesc: "prompt-pulse full-screen terminal dashboard",
		Synopsis:  "prompt-pulse-tui [options]",
		Description: `The TUI mode launches a full-screen interactive dashboard using Bubbletea v2.
It displays widgets in a configurable grid layout with real-time data updates.

Navigation supports vim-style keys (h/j/k/l), mouse clicks, and tab cycling.
The layout is defined by presets or custom row configurations in the config file.`,
		Options: `.TP
.B \-\-preset <name>
Override layout preset (dashboard, minimal, ops, billing).
.TP
.B \-\-refresh <duration>
Override widget refresh interval.
.TP
.B \-\-no-mouse
Disable mouse support.
.TP
.B \-\-no-vim
Disable vim-style keybindings.`,
		Examples: `.nf
# Launch TUI
prompt-pulse-tui

# Use minimal layout
prompt-pulse-tui --preset minimal

# Ops layout without mouse
prompt-pulse-tui --preset ops --no-mouse
.fi`,
		SeeAlso: `.BR prompt-pulse (1),
.BR prompt-pulse-banner (1),
.BR prompt-pulse.toml (5)`,
	}
}

func dcManConfig() *ManPage {
	return &ManPage{
		Name:      "prompt-pulse.toml",
		Section:   "5",
		ShortDesc: "prompt-pulse configuration file format",
		Synopsis:  "$XDG_CONFIG_HOME/prompt-pulse/config.toml",
		Description: `prompt-pulse uses a TOML configuration file with nested tables for each subsystem.
The file is searched for in $XDG_CONFIG_HOME/prompt-pulse/config.toml, falling back
to ~/.config/prompt-pulse/config.toml.

If no configuration file is found, built-in defaults are used. Environment variables
can override specific settings (see ENVIRONMENT section below).

The configuration is organized into these top-level tables: general, layout,
collectors, image, theme, shell, and banner.`,
		Options: `Environment variable overrides:

.TP
.B ANTHROPIC_ADMIN_KEY
Overrides collectors.claude.admin_key.
.TP
.B CIVO_TOKEN
Overrides collectors.billing.civo.api_key.
.TP
.B DIGITALOCEAN_TOKEN
Overrides collectors.billing.digitalocean.api_key.
.TP
.B PPULSE_PROTOCOL
Overrides image.protocol.
.TP
.B PPULSE_THEME
Overrides theme.name.
.TP
.B PPULSE_LAYOUT
Overrides layout.preset.`,
		Examples: `.nf
[general]
log_level = "info"
daemon_poll_interval = "15m"

[layout]
preset = "dashboard"

[collectors.sysmetrics]
enabled = true
interval = "1s"

[collectors.tailscale]
enabled = true
interval = "30s"

[image]
protocol = "auto"
waifu_enabled = true

[theme]
name = "catppuccin"

[shell]
tui_keybinding = "\\C-p"
show_banner_on_startup = true
.fi`,
		SeeAlso: `.BR prompt-pulse (1),
.BR prompt-pulse-daemon (1)`,
	}
}
