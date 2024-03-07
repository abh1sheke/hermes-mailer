package multi

import "github.com/spf13/cobra"

var senders, receivers, subject, host, readReceipts string
var textContent, htmlContent string
var workers uint8

// Cmd is the command defenition for the "multi" command.
// "multi" allows users to send email messages from multiple senders.
var Cmd = &cobra.Command{
	Use:   "multi",
	Short: "Send email messages from multiple senders",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&senders, "senders", "s", "", "Path to file containing senders")
	Cmd.Flags().StringVarP(&receivers, "receivers", "r", "", "Path to file containing receivers")
	Cmd.Flags().StringVarP(&subject, "subject", "S", "", "Sets the subject for the email messages")
	Cmd.Flags().StringVar(&host, "host", "", "Sets the SMTP host server for the senders")
	Cmd.Flags().StringVarP(&readReceipts, "read-receipts", "R", "", "Sets the email to which read-receipts are sent")
	Cmd.Flags().StringVarP(&textContent, "text", "t", "", "Path to the file containing plaintext email content")
	Cmd.Flags().StringVarP(&htmlContent, "html", "", "", "Path to the file containig html email content")

	Cmd.Flags().Uint8VarP(&workers, "workers", "", 2, "Sets the number of simultaneous send operations")

	Cmd.MarkFlagsRequiredTogether("senders", "receivers", "subject", "host", "text", "html")
}
