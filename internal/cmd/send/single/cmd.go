package single

import "github.com/spf13/cobra"

var sender, password, receivers, receiversFile, subject, host string
var textContent, htmlContent string

// Cmd is the command defenition for the "single" command.
// "single" allows users to send email messages from a single sender.
var Cmd = &cobra.Command{
	Use:   "single",
	Short: "Send email messages from a single sender",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&sender, "sender", "s", "", "Sets the sender email address")
	Cmd.Flags().StringVarP(&password, "password", "p", "", "Sets the password for the sender email address")
	Cmd.Flags().StringVarP(&receivers, "receivers", "r", "", "A comma seperated string of receiver email addresses")
	Cmd.Flags().StringVarP(&receiversFile, "receivers-file", "R", "", "Path to file containing receivers")
	Cmd.Flags().StringVarP(&subject, "subject", "S", "", "Sets the subject for the email messages")
	Cmd.Flags().StringVarP(&host, "host", "H", "", "Sets the SMTP host server for the sender")
	Cmd.Flags().StringVarP(&textContent, "text", "t", "", "Path to the file containing the plaintext email content")
	Cmd.Flags().StringVar(&htmlContent, "html", "", "Path to the file containing the HTML email content")

	Cmd.MarkFlagsRequiredTogether("sender", "password", "subject", "host", "text", "html")

	Cmd.MarkFlagsOneRequired("receivers", "receivers-file")
	Cmd.MarkFlagsMutuallyExclusive("receivers", "receivers-file")
}