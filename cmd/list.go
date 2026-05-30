package cmd

import (
	"os"

	"github.com/hagatasdelus/agysession/internal/list"
	"github.com/spf13/cobra"
)

var (
	listGrep    string
	listRegex   bool
	listColor   string
	listNoColor bool
)

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "emit TSV rows for fzf",
	RunE: func(cmd *cobra.Command, args []string) error {
		effExcludeDir := excludeDir
		if effExcludeDir == "" {
			effExcludeDir = os.Getenv("AGYSESSION_EXCLUDE_DIR")
		}
		return list.Run(list.Options{
			Grep:       listGrep,
			ExcludeDir: effExcludeDir,
			Regex:      listRegex,
			Color:      listColor,
			NoColor:    listNoColor,
			Out:        os.Stdout,
		})
	},
}

func init() {
	listCmd.Flags().StringVar(&listGrep, "grep", "", "filter sessions by user/model content (fixed-string)")
	listCmd.Flags().BoolVar(&listRegex, "regex", false, "treat --grep query as a regular expression")
	listCmd.Flags().StringVar(&listColor, "color", "auto", "color output: auto (default) | always | never")
	listCmd.Flags().BoolVar(&listNoColor, "no-color", false, "shorthand for --color=never")
	rootCmd.AddCommand(listCmd)
}
