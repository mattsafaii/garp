package cmd

import (
	"fmt"
	"garp/internal/deploy"

	"github.com/spf13/cobra"
)

var deployHistoryCmd = &cobra.Command{
	Use:   "deploy-history",
	Short: "Show deployment history",
	Long:  `Display recent deployment history with details about each deployment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeployHistory()
	},
}

var (
	historyLimit int
)

func runDeployHistory() error {
	history, err := deploy.NewDeploymentHistory()
	if err != nil {
		return fmt.Errorf("failed to load deployment history: %v", err)
	}

	recent := history.GetRecentDeployments(historyLimit)

	if len(recent) == 0 {
		fmt.Println("No deployments found.")
		return nil
	}

	fmt.Printf("Recent deployments (showing %d):\n\n", len(recent))

	for _, record := range recent {
		status := "✅ SUCCESS"
		if !record.Success {
			status = "❌ FAILED"
		}

		fmt.Printf("ID: %s\n", record.ID)
		fmt.Printf("Time: %s\n", record.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("Strategy: %s\n", record.Strategy)
		if record.Target != "" {
			fmt.Printf("Target: %s\n", record.Target)
		}
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Duration: %v\n", record.Duration)

		if record.URL != "" {
			fmt.Printf("URL: %s\n", record.URL)
		}

		if record.GitBranch != "" {
			fmt.Printf("Git Branch: %s\n", record.GitBranch)
		}

		if record.GitCommit != "" {
			fmt.Printf("Git Commit: %s\n", record.GitCommit)
		}

		if len(record.Messages) > 0 {
			fmt.Println("Messages:")
			for _, msg := range record.Messages {
				fmt.Printf("  • %s\n", msg)
			}
		}

		if len(record.Errors) > 0 {
			fmt.Println("Errors:")
			for _, err := range record.Errors {
				fmt.Printf("  ⚠️  %s\n", err)
			}
		}

		fmt.Println()
	}

	return nil
}

func init() {
	deployHistoryCmd.Flags().IntVar(&historyLimit, "limit", 10, "Number of recent deployments to show")
	rootCmd.AddCommand(deployHistoryCmd)
}
