package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage asset-generator CLI configuration",
	Long: `View and modify the asset-generator CLI configuration.

Configuration is stored in ~/.asset-generator/config.yaml by default.

Examples:
  # View current configuration
  asset-generator config view
  
  # Set API URL
  asset-generator config set api-url https://api.swarm.example.com
  
  # Set API key
  asset-generator config set api-key your-api-key-here`,
}

// configViewCmd displays current configuration
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Long:  `Display the current asset-generator CLI configuration.`,
	RunE:  runConfigView,
}

// configSetCmd sets a configuration value
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value in the config file.

Examples:
  asset-generator config set api-url https://api.swarm.example.com
  asset-generator config set api-key your-key-here`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

// configGetCmd gets a configuration value
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long: `Get a specific configuration value.

Examples:
  asset-generator config get api-url
  asset-generator config get api-key`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

// configInitCmd initializes a new config file
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new config file",
	Long:  `Create a new configuration file with default values.`,
	RunE:  runConfigInit,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configInitCmd)
}

func runConfigView(cmd *cobra.Command, args []string) error {
	settings := viper.AllSettings()

	if len(settings) == 0 {
		fmt.Println("No configuration found. Run 'asset-generator config init' to create a default config.")
		return nil
	}

	fmt.Println("Current configuration:")
	fmt.Println("=====================")
	for key, value := range settings {
		// Mask sensitive values
		if key == "api-key" || key == "api_key" {
			if value != "" && value != nil {
				fmt.Printf("%s: %s\n", key, "********")
			} else {
				fmt.Printf("%s: (not set)\n", key)
			}
		} else {
			fmt.Printf("%s: %v\n", key, value)
		}
	}

	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("\nConfig file: %s\n", configFile)
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	viper.Set(key, value)

	// Ensure config file exists
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := home + "/.asset-generator"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := configDir + "/config.yaml"
	if viper.ConfigFileUsed() == "" {
		viper.SetConfigFile(configFile)
	}

	if err := viper.WriteConfig(); err != nil {
		// If config doesn't exist, create it
		if err := viper.SafeWriteConfig(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)
	fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := viper.Get(key)

	if value == nil {
		return fmt.Errorf("configuration key '%s' not found", key)
	}

	// Mask sensitive values
	if key == "api-key" || key == "api_key" {
		fmt.Println("********")
	} else {
		fmt.Println(value)
	}

	return nil
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := home + "/.asset-generator"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := configDir + "/config.yaml"
	viper.SetConfigFile(configFile)

	// Set default values
	viper.SetDefault("api-url", "http://localhost:7801")
	viper.SetDefault("format", "table")
	viper.SetDefault("quiet", false)
	viper.SetDefault("verbose", false)

	if err := viper.SafeWriteConfig(); err != nil {
		// Config already exists
		if os.IsExist(err) {
			return fmt.Errorf("config file already exists at %s", configFile)
		}
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("✓ Configuration file created: %s\n", configFile)
	fmt.Println("\nDefault values:")
	fmt.Println("  api-url: http://localhost:7801")
	fmt.Println("  format: table")
	fmt.Println("\nSet your API key with:")
	fmt.Println("  asset-generator config set api-key YOUR_API_KEY")

	return nil
}
