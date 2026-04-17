package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agenkit",
	Short: "Manage AI agent configurations from one source of truth",
	Long: "Agenkit is a CLI for defining MCP servers, hooks, and instructions once\n" +
		"and applying them across tools like Copilot, Claude, Cursor, and Codex.\n\n" +
		"It keeps agent configuration declarative, reviewable, and portable\n" +
		"through a single manifest instead of scattered provider-specific files.",

	SilenceUsage:  true,
	SilenceErrors: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose responses")

}

func Execute() error {
	return rootCmd.Execute()
}
