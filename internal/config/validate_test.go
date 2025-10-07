package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "valid http URL",
			setup: func() {
				viper.Reset()
				viper.Set("api-url", "http://localhost:7801")
				viper.Set("format", "json")
			},
			wantErr: false,
		},
		{
			name: "valid https URL",
			setup: func() {
				viper.Reset()
				viper.Set("api-url", "https://api.example.com")
				viper.Set("format", "table")
			},
			wantErr: false,
		},
		{
			name: "invalid URL - no scheme",
			setup: func() {
				viper.Reset()
				viper.Set("api-url", "localhost:7801")
				viper.Set("format", "json")
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			setup: func() {
				viper.Reset()
				viper.Set("api-url", "http://localhost:7801")
				viper.Set("format", "invalid")
			},
			wantErr: true,
		},
		{
			name: "valid format - table",
			setup: func() {
				viper.Reset()
				viper.Set("api-url", "http://localhost:7801")
				viper.Set("format", "table")
			},
			wantErr: false,
		},
		{
			name: "valid format - yaml",
			setup: func() {
				viper.Reset()
				viper.Set("api-url", "http://localhost:7801")
				viper.Set("format", "yaml")
			},
			wantErr: false,
		},
		{
			name: "empty config",
			setup: func() {
				viper.Reset()
				viper.Set("format", "table")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := ValidateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http", "http://localhost:7801", false},
		{"valid https", "https://api.example.com", false},
		{"valid ws", "ws://localhost:7801", false},
		{"valid wss", "wss://api.example.com", false},
		{"no scheme", "localhost:7801", true},
		{"invalid scheme", "ftp://localhost:7801", true},
		{"no host", "http://", true},
		{"malformed", "://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		format string
		want   bool
	}{
		{"table", true},
		{"json", true},
		{"yaml", true},
		{"TABLE", true},
		{"JSON", true},
		{"YAML", true},
		{"xml", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			if got := isValidFormat(tt.format); got != tt.want {
				t.Errorf("isValidFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
