package cmd

import (
	"errors"

	"github.com/hagatasdelus/agysession/internal/resume"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume <sessionId>",
	Short: "chdir to original cwd, exec agy --conversation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("session id required")
		}
		return resume.Run(args[0])
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
