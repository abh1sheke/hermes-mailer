package cmd

import "github.com/spf13/cobra"
import "github.com/abh1sheke/hermes-mailer/internal/cmd/send"

var rootCmd = &cobra.Command{
	Use:   "hermes",
	Short: "A command-line tool for various email operations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(send.Cmd)

	rootCmd.PersistentFlags().Uint8P("log-level", "l", 4, "Sets the log level")
}

// Execute runs the 'Root' cobra-cli command.
func Execute() error {
	return rootCmd.Execute()
}
