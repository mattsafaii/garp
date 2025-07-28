package cmd

import (
	"fmt"
	"garp-cli/internal/deploy"

	"github.com/spf13/cobra"
)

var deployConfigCmd = &cobra.Command{
	Use:   "deploy-config",
	Short: "Manage deployment configurations",
	Long:  `Manage deployment environment configurations for different targets.`,
}

var setConfigCmd = &cobra.Command{
	Use:   "set [environment-name]",
	Short: "Set deployment configuration for an environment",
	Long:  `Set deployment configuration for a specific environment (e.g., staging, production).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSetConfig,
}

var getConfigCmd = &cobra.Command{
	Use:   "get [environment-name]",
	Short: "Get deployment configuration for an environment",
	Long:  `Display deployment configuration for a specific environment.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGetConfig,
}

var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployment configurations",
	Long:  `List all configured deployment environments.`,
	RunE:  runListConfig,
}

var removeConfigCmd = &cobra.Command{
	Use:   "remove [environment-name]",
	Short: "Remove deployment configuration",
	Long:  `Remove deployment configuration for a specific environment.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRemoveConfig,
}

var (
	configStrategy string
	configValues   []string
)

func runSetConfig(cmd *cobra.Command, args []string) error {
	envName := args[0]
	
	configManager, err := deploy.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to initialize config manager: %v", err)
	}
	
	// Parse config values into map
	configMap := make(map[string]string)
	for _, val := range configValues {
		// Parse key=value format
		if len(val) > 0 {
			// Simple parsing - in production you might want more robust parsing
			fmt.Printf("Adding config: %s\n", val)
			configMap["raw"] = val // Simplified for demo
		}
	}
	
	config := deploy.EnvironmentConfig{
		Strategy: configStrategy,
		Config:   configMap,
	}
	
	err = configManager.SetEnvironment(envName, config)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}
	
	fmt.Printf("✅ Configuration saved for environment '%s'\n", envName)
	return nil
}

func runGetConfig(cmd *cobra.Command, args []string) error {
	envName := args[0]
	
	configManager, err := deploy.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to initialize config manager: %v", err)
	}
	
	config, err := configManager.GetEnvironment(envName)
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}
	
	fmt.Printf("Environment: %s\n", config.Name)
	fmt.Printf("Strategy: %s\n", config.Strategy)
	fmt.Println("Configuration:")
	for key, value := range config.Config {
		fmt.Printf("  %s: %s\n", key, value)
	}
	
	return nil
}

func runListConfig(cmd *cobra.Command, args []string) error {
	configManager, err := deploy.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to initialize config manager: %v", err)
	}
	
	environments := configManager.ListEnvironments()
	
	if len(environments) == 0 {
		fmt.Println("No deployment configurations found.")
		return nil
	}
	
	fmt.Println("Configured environments:")
	for _, env := range environments {
		config, err := configManager.GetEnvironment(env)
		if err != nil {
			fmt.Printf("  %s (error loading config)\n", env)
			continue
		}
		fmt.Printf("  %s (%s)\n", env, config.Strategy)
	}
	
	return nil
}

func runRemoveConfig(cmd *cobra.Command, args []string) error {
	envName := args[0]
	
	configManager, err := deploy.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to initialize config manager: %v", err)
	}
	
	err = configManager.RemoveEnvironment(envName)
	if err != nil {
		return fmt.Errorf("failed to remove config: %v", err)
	}
	
	fmt.Printf("✅ Configuration removed for environment '%s'\n", envName)
	return nil
}

func init() {
	// Add flags for set command
	setConfigCmd.Flags().StringVar(&configStrategy, "strategy", "", "Deployment strategy (git, rsync, netlify, cloudflare)")
	setConfigCmd.Flags().StringArrayVar(&configValues, "config", []string{}, "Configuration values (can be specified multiple times)")
	setConfigCmd.MarkFlagRequired("strategy")
	
	// Add subcommands
	deployConfigCmd.AddCommand(setConfigCmd)
	deployConfigCmd.AddCommand(getConfigCmd)
	deployConfigCmd.AddCommand(listConfigCmd)
	deployConfigCmd.AddCommand(removeConfigCmd)
	
	rootCmd.AddCommand(deployConfigCmd)
}