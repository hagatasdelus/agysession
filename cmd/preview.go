package cmd

import (
	"errors"

	"github.com/hagatasdelus/agysession/internal/preview"
	"github.com/spf13/cobra"
)

var (
	previewQuery string
	previewRegex bool
)

var previewCmd = &cobra.Command{
	Use:   "preview [flags] <sessionId>",
	Short: "render the preview pane for a session id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("session id required")
		}
		return preview.Run(args[0], preview.Options{
			Query: previewQuery,
			Regex: previewRegex,
		})
	},
}

func init() {
	previewCmd.Flags().StringVar(&previewQuery, "query", "", "highlight matches of <query> in the preview (fixed-string)")
	previewCmd.Flags().BoolVar(&previewRegex, "regex", false, "treat --query as a regular expression")
	rootCmd.AddCommand(previewCmd)
}
