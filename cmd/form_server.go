package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var formServerCmd = &cobra.Command{
	Use:   "form-server",
	Short: "Start the contact form server",
	Long: `Start the Sinatra server for handling contact form 
submissions with Resend email integration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Starting form server on port %d\n", formPort)
		// TODO: Implement Sinatra server startup
		return fmt.Errorf("form-server command not yet implemented")
	},
}

var (
	formPort int
)

func init() {
	formServerCmd.Flags().IntVarP(&formPort, "port", "p", 4567, "Port for form server")
	rootCmd.AddCommand(formServerCmd)
}