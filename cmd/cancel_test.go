package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/opd-ai/asset-generator/pkg/client"
)

func TestCancelCommand(t *testing.T) {
	// This is a basic test to ensure the command structure is correct
	// Full integration tests require a running SwarmUI server

	// Save original client
	origClient := assetClient
	defer func() { assetClient = origClient }()

	// Create a mock client
	clientConfig := &client.Config{
		BaseURL: "http://localhost:7801",
		Verbose: false,
	}

	mockClient, err := client.NewAssetClient(clientConfig)
	if err != nil {
		t.Fatalf("Failed to create mock client: %v", err)
	}
	assetClient = mockClient

	// Test that the command exists and has proper structure
	if cancelCmd == nil {
		t.Fatal("cancelCmd is nil")
	}

	if cancelCmd.Use != "cancel" {
		t.Errorf("Expected Use='cancel', got '%s'", cancelCmd.Use)
	}

	if cancelCmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cancelCmd.Long == "" {
		t.Error("Long description is empty")
	}

	if cancelCmd.RunE == nil {
		t.Error("RunE is nil")
	}

	// Test that --all flag exists
	allFlag := cancelCmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("--all flag not found")
	}

	if allFlag.DefValue != "false" {
		t.Errorf("Expected --all default to be false, got %s", allFlag.DefValue)
	}
}

func TestCancelCommandHelp(t *testing.T) {
	// Test that help output can be generated without panic
	var buf bytes.Buffer
	cancelCmd.SetOut(&buf)
	cancelCmd.SetErr(&buf)

	// Reset flags to defaults
	cancelCmd.Flags().Set("all", "false")

	// Get help
	cancelCmd.SetArgs([]string{"--help"})

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Help generation panicked: %v", r)
		}
	}()

	// Execute help (will exit with code, but we catch it)
	_ = cancelCmd.Execute()
}

func TestCancelFunctionSignatures(t *testing.T) {
	// Verify that the client methods exist with correct signatures
	// This ensures compile-time compatibility

	cfg := &client.Config{
		BaseURL: "http://localhost:7801",
	}

	c, err := client.NewAssetClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// These calls will fail without a real server, but we're testing
	// that the methods exist with the right signatures
	t.Run("Interrupt method exists", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Interrupt method panicked: %v", r)
			}
		}()
		// Method should exist and accept context
		_ = c.Interrupt(ctx)
	})

	t.Run("InterruptAll method exists", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("InterruptAll method panicked: %v", r)
			}
		}()
		// Method should exist and accept context
		_ = c.InterruptAll(ctx)
	})
}

func TestCancelCommandIntegration(t *testing.T) {
	// Skip this test in CI or without SwarmUI
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would be a full integration test with a real SwarmUI instance
	// For now, we just verify the command structure
	t.Run("Command registration", func(t *testing.T) {
		// Check that cancel command is registered in root
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == "cancel" {
				found = true
				break
			}
		}
		if !found {
			t.Error("cancel command not registered in root command")
		}
	})
}
