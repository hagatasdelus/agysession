# agysession

> An fzf-powered session picker for `agy --conversation`.

`agysession` lists every Antigravity CLI session under `~/.gemini/antigravity-cli/brain`, lets you fuzzy-find across all of your projects with a live preview pane, and resumes the one you pick in its original working directory.

This tool is inspired by [`ccsession`](https://github.com/sorafujitani/ccsession) but adapted for the Antigravity CLI (`agy`).

---

## Features

- **Cross-project listing** — every session from every workspace in one view, sorted by last activity.
- **Ultra-fast scanning** — reads the global index `history.jsonl` to instantly list sessions without traversing the entire filesystem.
- **Three search modes** — fuzzy (default), directory-only, and full-text grep over JSONL transcripts.
- **Live preview** — last 30 messages of the highlighted session, with timestamps and roles. In grep mode the matched query is highlighted in the preview so you can spot the hit at a glance.
- **Faithful resume** — `chdir`s back to the session's original `cwd` (workspace) before exec'ing `agy --conversation`, so paths and tooling Just Work.
- **Single static binary** — written in Go utilizing `spf13/cobra` CLI framework.

## Requirements

| Tool | Required for |
| --- | --- |
| [`fzf`](https://github.com/junegunn/fzf) | interactive picker |
| `agy` ([Antigravity CLI](https://antigravity.google/docs/cli-overview)) | resuming sessions |

## Install

### Go

```sh
go install github.com/hagatasdelus/agysession@latest
```

Version metadata is recovered from `runtime/debug.ReadBuildInfo`, so `agysession --version` works for `go install` builds as well.

## Usage

```sh
agysession                            # list -> fzf -> resume
agysession list  [--grep Q] [--regex] # emit TSV rows to stdout
agysession preview [--query Q] [--regex] <id> # render the preview pane (Q highlighted)
agysession resume  <id>               # chdir to the session's CWD, exec `agy --conversation`
agysession --version
agysession --help
```

### Keys inside fzf

| Key      | Mode |
| -------- | --- |
| `Ctrl-G` | grep — refilters by user/model content on every keystroke; matches are highlighted in the preview |
| `Ctrl-O` | dir — fuzzy match restricted to the directory column |
| `Ctrl-F` | fuzzy — default; matches across time / dir / label |
| `Enter`  | resume the selected session |
| `Esc`    | cancel |

## How it works

1. `agysession list` reads `~/.gemini/antigravity-cli/history.jsonl` to collect active session metadata, dedupes them by `conversationId`, and prints one TSV row per session (`id`, `epoch`, relative time, cwd basename, label).
2. `fzf` consumes the TSV. The three key bindings swap fzf's matcher between fuzzy mode, directory-only mode, and grep mode (which reloads via `agysession list --grep <query>` on every keystroke). The current query is also forwarded to the preview as `agysession preview --query <query> <id>`, which highlights its matches in the rendered messages.
3. On `Enter`, `agysession resume <id>` resolves the session's original workspace (CWD), `chdir`s into it, and `execve`s `agy --conversation <id>` so the resumed process fully replaces the picker.

## Inspired by

- [`ccsession`](https://github.com/sorafujitani/ccsession)

## License

[MIT](./LICENSE)
