# Nushell Notes for prompt-pulse

As of 2026-04-23, `prompt-pulse` does not generate a Nushell integration script.

The supported generated shell targets are:

- `bash`
- `zsh`
- `fish`
- `ksh`

That means `prompt-pulse shell nushell` is not a valid current workflow.

## Current Recommendation

If you use Nushell, wire the interactive path manually to the separate
`prompt-pulse-tui` binary instead of expecting generated prompt-pulse shell
helpers.

Example keybinding entry for `$env.config.keybindings`:

```nu
{
    name: prompt_pulse_tui
    modifier: control
    keycode: char_p
    mode: [emacs vi_normal vi_insert]
    event: {
        send: executehostcommand
        cmd: "prompt-pulse-tui"
    }
}
```

Optional banner helper:

```nu
def pp-banner [] {
    let session_id = $"($nu.pid)-(date now | format date '%s')"
    prompt-pulse --banner --session-id $session_id
}
```

## What Not To Copy Forward

The following older patterns are stale and should not be reused:

- `prompt-pulse --tui`
- `prompt-pulse shell nushell`
- generated `pp-tui`, `pp-daemon-start`, or `pp-daemon-stop` Nushell helpers

If Nushell generation is implemented in the future, it should be documented as
a new supported shell target instead of inferred from old notes like this one.
