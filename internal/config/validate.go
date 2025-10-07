package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

// ValidateConfig validates the configuration
func ValidateConfig() error {
	// Validate API URL
	apiURL := viper.GetString("api-url")
	if apiURL != "" {
		if err := validateURL(apiURL); err != nil {
			return fmt.Errorf("invalid api-url: %w", err)
		}
	}

	// Validate output format
	format := viper.GetString("format")
	if !isValidFormat(format) {
		return fmt.Errorf("invalid format: %s (must be table, json, or yaml)", format)
	}

	return nil
}

// validateURL validates a URL
func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" && parsedURL.Scheme != "ws" && parsedURL.Scheme != "wss" {
		return fmt.Errorf("URL scheme must be http, https, ws, or wss")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include host")
	}

	return nil
}

// isValidFormat checks if the format is valid
func isValidFormat(format string) bool {
	format = strings.ToLower(format)
	return format == "table" || format == "json" || format == "yaml"
}
