---
name: vhs-tape
description: Create a VHS .tape file for recording terminal GIF/MP4/WebM demos. Use when the user wants to create a .tape file for VHS terminal recordings.
---

# VHS .tape File Reference

VHS reads `.tape` files and renders them into terminal recordings (GIF, MP4, WebM, or PNG sequences). A tape file is a linear sequence of commands — no loops or conditionals.

## File Structure Rule

**`Set` commands MUST come before all interaction commands.** Order within the Set block does not matter.

---

## Output & Require

```
Output demo.gif          # GIF (most common)
Output demo.mp4          # MP4 video
Output demo.webm         # WebM video
Output frames/           # Directory for PNG frame sequence
Output demo.txt          # ASCII text (for testing)

Require echo             # Abort with clear error if program not in PATH
Require fzf
```

Multiple `Output` lines are allowed — VHS renders to all of them in one run.

---

## Set Commands

### Terminal Dimensions
```
Set Width  1200          # pixels (default: 1200)
Set Height 600           # pixels (default: 600)
```

### Typography
```
Set FontSize      22                   # pt (default: 22)
Set FontFamily    "RobotoMono Nerd Font"  # must be installed on system
Set LetterSpacing 1                    # multiplier (default: 1.0)
Set LineHeight    1.2                  # multiplier (default: 1.0)
```

### Shell
```
Set Shell "bash"         # bash, zsh, fish, etc.
```

### Timing & Playback
```
Set TypingSpeed   50ms       # delay between typed chars (default: 50ms)
Set Framerate     60         # capture framerate (default: 50)
Set PlaybackSpeed 2          # output speed multiplier (0.5 = slow, 2.0 = fast)
Set LoopOffset    60.4       # % into GIF where looping restarts (avoids blank-terminal loop)
Set LoopOffset    20.99%     # same, with explicit %
```

### Wait Behavior
```
Set WaitTimeout  15s         # how long Wait polls before failing (default: 15s)
Set WaitPattern  /regex/     # default pattern for bare Wait command
```

### Theme
```
Set Theme "Catppuccin Mocha"    # named theme (350+ available; run `vhs themes` to list)
Set Theme "Dracula"
Set Theme "Nord"
Set Theme "TokyoNight"
Set Theme "One Dark"
Set Theme "Gruvbox Dark"
Set Theme "Rose Pine"

# Inline JSON custom theme:
Set Theme { "name": "Custom", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "purple": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightPurple": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#29283b", "foreground": "#b3b0d6", "selectionBackground": "#3d3c58", "cursorColor": "#b3b0d6" }
```

### Visual / Layout
```
Set Padding      50          # internal padding in pixels
Set Margin       20          # external border in pixels
Set MarginFill   "#6B50FF"   # margin area color (hex)
Set BorderRadius 10          # corner radius in pixels
Set CursorBlink  false       # cursor animation (default: true)
```

### Window Bar (macOS-style chrome)
```
Set WindowBar         "Colorful"       # colored traffic-light buttons (recommended for docs)
Set WindowBar         "ColorfulRight"  # buttons on the right
Set WindowBar         "Rings"          # monochrome rings
Set WindowBar         "RingsRight"     # rings on the right
Set WindowBarSize     40               # height in pixels
```

---

## Sleep

```
Sleep 1s          # 1 second
Sleep 500ms       # 500 milliseconds
Sleep .5          # 0.5 seconds
Sleep 1m          # 1 minute
```

### Time Format Summary

| Format     | Example  | Meaning        |
|------------|----------|----------------|
| `ms` suffix | `500ms` | Milliseconds   |
| `s` suffix  | `1.5s`  | Seconds        |
| `m` suffix  | `1m`    | Minutes        |
| Decimal    | `.5`     | 0.5 seconds    |

---

## Type

```
Type "hello world"             # type at global TypingSpeed
Type 'single quotes work'
Type `backticks work`

Type@500ms "slow emphasis"     # override speed for this command only
Type@10ms  "fast burst"
```

---

## Key Press Commands

**Syntax:** `<Key>[@duration] [count]`

```
Enter                    # press Enter once
Enter 3                  # press Enter 3 times
Enter@100ms 3            # press 3 times with 100ms between each

Backspace
Backspace 2
Backspace@200ms 5

Delete
Delete 2

Insert
Tab
Tab 2

Space
Space 3

Escape

Up
Down
Left
Right
Up 2
Down@50ms 3

PageUp
PageDown
PageUp 2
PageDown@100ms 3

ScrollUp
ScrollDown
ScrollUp 2
ScrollDown@100ms 2
```

---

## Modifier Key Commands

```
Ctrl+C
Ctrl+L
Ctrl+R
Ctrl+D
Ctrl+Z
Ctrl+A              # any letter
Ctrl+Left           # ctrl + arrow
Ctrl+Right
Ctrl+Up
Ctrl+Down

Alt+.
Alt+L
Alt+b
Alt+f

Shift+Enter
Shift+Tab

# Multiple modifiers:
Ctrl+Shift+Alt+C
Ctrl+Alt+Down
```

---

## Wait

Blocks until a regex matches or the timeout expires. More reliable than guessing `Sleep` durations for commands with variable output time.

```
Wait                          # wait using WaitPattern setting
Wait /regex/                  # wait until regex matches current line
Wait+Line /regex/             # explicit: match current line
Wait+Screen /regex/           # match anywhere on screen buffer

Wait@5s /pattern/             # 5-second timeout override for this Wait
Wait+Screen@10s /done/
```

Polls every 10ms. Uses Go regex syntax.

---

## Hide / Show

Hide setup boilerplate (prompt customization, `cd`, env exports) from the recording:

```
Hide
Type "export PS1='$ '"
Enter
Sleep 300ms
Type "cd /tmp/demo"
Enter
Sleep 300ms
Show
# Recording content starts here
```

---

## Screenshot

Capture a single PNG frame at the current moment:

```
Screenshot examples/output.png
```

---

## Clipboard

```
Copy "text to copy"    # copy to clipboard
Paste                  # paste clipboard contents into terminal
```

---

## Environment Variables

```
Env HOME "/tmp/demo"
Env FOO  "bar"
```

---

## Source

Include another tape file's commands inline:

```
Source setup.tape
Source common/header.tape
```

---

## Comments

```
# This is a comment
```

---

## Best Practices

1. **Put all `Set` commands first** — they must precede any interaction commands.
2. **Use `Hide`/`Show` for setup** — hide shell prompt changes, `cd` commands, env exports that clutter the demo.
3. **Add `Sleep` after commands** — give the terminal time to render output before the next keystroke or the recording ends.
4. **Use `Wait` instead of long `Sleep`** for commands with variable duration — `Wait /\$\s/` waits for the prompt to return rather than guessing.
5. **Use `Require`** to fail fast with a clear error if a dependency is missing.
6. **Multiple `Output` lines** — produce GIF + MP4 in one run.
7. **`Set LoopOffset`** — set to the percentage where the action starts so the loop feels natural, not jumping to a blank terminal.
8. **`Set WindowBar "Colorful"`** — adds macOS window chrome; makes demos look polished in documentation.
9. **`Type@speed` modifiers** — use slow typing for emphasis, fast for showing speed.
10. **`Source` for shared setup** — put common `Set` blocks and `Hide`/`Show` setup in a shared tape for reuse.

---

## Complete Example

```
# tfx demo recording
Output demo/tfx-demo.gif
Output demo/tfx-demo.mp4

Require tfx

Set Shell        "zsh"
Set Width        1200
Set Height       600
Set FontSize     14
Set FontFamily   "RobotoMono Nerd Font"
Set Theme        "Catppuccin Mocha"
Set Padding      20
Set WindowBar    "Colorful"
Set BorderRadius 8
Set TypingSpeed  100ms
Set Framerate    60
Set WaitTimeout  30s
Sleep 1s
Type "./tfx tui"
Enter
Sleep 2s

# ============================================================
# RECORDING STARTS HERE
# ============================================================

```

---

## CLI Reference

```bash
vhs demo.tape              # execute a tape file
vhs new demo.tape          # create a starter tape file
vhs record > session.tape  # record a live terminal session to tape
vhs themes                 # list all 350+ available theme names
vhs manual                 # view the full manual
```

## Task

Create a `.tape` file for: $ARGUMENTS
