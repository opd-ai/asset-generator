package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/opd-ai/asset-generator/internal/config"
	"github.com/opd-ai/asset-generator/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	apiURL     string
	apiKey     string
	outputFmt  string
	outputFile string
	quiet      bool
	verbose    bool

	assetClient *client.AssetClient
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "asset-generator",
	Short: "CLI client for asset generation APIs",
	Long: `asset-generator is a command-line interface for interacting with asset generation APIs.
It provides tools for generating assets, managing models, and configuring
the asset generation service.

Examples:
  asset-generator generate image --prompt "a beautiful landscape"
  asset-generator models list
  asset-generator config set api-url https://api.example.com`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if this is a config command - they need special handling
		// Config commands must work even when current config is invalid
		isConfigCommand := cmd.Name() == "config" || (cmd.Parent() != nil && cmd.Parent().Name() == "config")

		// Initialize configuration (skip validation for config commands)
		if err := initConfigWithValidation(!isConfigCommand); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		// Skip client initialization for config commands
		if isConfigCommand {
			return nil
		}

		// Initialize client for non-config commands
		clientCfg := &client.Config{
			BaseURL: viper.GetString("api-url"),
			APIKey:  viper.GetString("api-key"),
			Verbose: verbose,
		}

		var err error
		assetClient, err = client.NewAssetClient(clientCfg)
		if err != nil {
			return fmt.Errorf("failed to create asset generation client: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.asset-generator/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "Asset generation API base URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Asset generation API key")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "format", "f", "table", "output format (table, json, yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "write output to file instead of stdout")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (errors only)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
// The validate parameter controls whether to validate config values.
// Set to false for config commands to allow fixing invalid configurations.
func initConfigWithValidation(validate bool) error {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Search config in home directory with name ".asset-generator" (without extension)
		configDir := home + "/.asset-generator"

		// NOTE: Viper searches config paths in REVERSE order (last added = first searched)
		// We add paths in reverse of their desired precedence to achieve:
		//   1st search (highest precedence): ./config/config.yaml
		//   2nd search (medium precedence):  ~/.asset-generator/config.yaml
		//   3rd search (lowest precedence):  /etc/asset-generator/config.yaml
		viper.AddConfigPath("/etc/asset-generator") // Searched last (lowest precedence)
		viper.AddConfigPath(configDir)              // Searched second
		viper.AddConfigPath("./config")             // Searched first (highest precedence)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml") // Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Read in environment variables that match
	viper.SetEnvPrefix("ASSET_GENERATOR")
	// Replace dashes with underscores for environment variable compatibility
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("api-url", "http://localhost:7801")
	viper.SetDefault("format", "table")
	viper.SetDefault("quiet", false)
	viper.SetDefault("verbose", false)

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if !quiet && verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Only validate if requested (skip for config commands)
	if validate {
		return config.ValidateConfig()
	}

	return nil
}
