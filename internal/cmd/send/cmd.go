package send

import (
	"github.com/abh1sheke/hermes-mailer/internal/cmd/send/multi"
	"github.com/abh1sheke/hermes-mailer/internal/cmd/send/single"
	"github.com/spf13/cobra"
)

// Cmd is the command defenition for the "single" command.
// "single" allows users to send email messages from a single sender.
var Cmd = &cobra.Command{
	Use:   "send",
	Short: "Send emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(single.Cmd)
	Cmd.AddCommand(multi.Cmd)
}
