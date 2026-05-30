package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/hagatasdelus/agysession/version"
	"github.com/spf13/cobra"
)

var (
	commit     = ""
	date       = ""
	excludeDir string
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= 7 && commit == "" {
				commit = s.Value[:7]
			}
		case "vcs.time":
			if date == "" {
				date = s.Value
			}
		}
	}
	if commit != "" {
		version.Revision = commit
	}
}

var rootCmd = &cobra.Command{
	Use:     "agysession [flags]",
	Short:   "agysession - fzf frontend for \"agy --conversation\"",
	Version: fmt.Sprintf("%s (rev: %s)", version.Version, version.Revision),
	RunE: func(cmd *cobra.Command, args []string) error {
		if excludeDir != "" {
			os.Setenv("AGYSESSION_EXCLUDE_DIR", excludeDir)
		}
		return runDefault()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&excludeDir, "exclude-dir", "", "hide sessions whose cwd contains <s> (case-insensitive)")
}

func runDefault() error {
	if _, err := exec.LookPath("fzf"); err != nil {
		return fmt.Errorf("fzf is required but not found in PATH")
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command("bash", "-c", defaultScript)
	cmd.Env = append(os.Environ(), "AGYSESSION_BIN="+self)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

const defaultScript = `set -u
exclude_args=()
exclude_arg=""
if [ -n "${AGYSESSION_EXCLUDE_DIR:-}" ]; then
  exclude_args=(--exclude-dir "$AGYSESSION_EXCLUDE_DIR")
  exclude_arg=$(printf -- '--exclude-dir %q' "$AGYSESSION_EXCLUDE_DIR")
fi
id=$("$AGYSESSION_BIN" list --color=always "${exclude_args[@]+"${exclude_args[@]}"}" | fzf \
  --ansi \
  --delimiter=$'\t' \
  --with-nth=3,4,5 \
  --nth=1,2,3 \
  --no-sort \
  --no-hscroll \
  --color='hl:-1:reverse,hl+:-1:reverse' \
  --preview "$AGYSESSION_BIN preview --query {q} {1}" \
  --preview-window=right,60%,wrap \
  --header='[fuzzy] ctrl-g: grep / ctrl-o: dir / ctrl-f: fuzzy / enter: resume' \
  --bind 'start:unbind(change)' \
  --bind "change:reload(sleep 0.05; $AGYSESSION_BIN list --color=always $exclude_arg --grep {q})" \
  --bind "ctrl-g:transform:echo \"change-prompt(grep> )+disable-search+reload(sleep 0.05; $AGYSESSION_BIN list --color=always $exclude_arg --grep {q})+rebind(change)\"" \
  --bind "ctrl-o:transform:echo \"change-prompt(dir> )+enable-search+change-nth(2)+reload($AGYSESSION_BIN list --color=always $exclude_arg)+unbind(change)\"" \
  --bind "ctrl-f:transform:echo \"change-prompt(> )+enable-search+change-nth(1,2,3)+reload($AGYSESSION_BIN list --color=always $exclude_arg)+unbind(change)\"" \
  | awk -F'\t' '{print $1}') || true
if [ -n "$id" ]; then
  exec "$AGYSESSION_BIN" resume "$id"
fi
`
