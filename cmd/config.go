/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mochajutsu/mkcd/internal/config"
	"github.com/mochajutsu/mkcd/internal/utils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage mkcd configuration",
	Long: `Manage mkcd configuration settings.

The config command allows you to initialize, view, edit, validate, and reset
your mkcd configuration. The configuration file is stored in TOML format
and contains settings for profiles, git integration, templates, and more.

Examples:
  mkcd config init                     # Initialize config file with defaults
  mkcd config show                     # Show current configuration
  mkcd config edit                     # Edit config in $EDITOR
  mkcd config validate                 # Validate configuration
  mkcd config reset                    # Reset to defaults`,
}

// configInitCmd represents the config init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Initialize the mkcd configuration file with default settings.`,
	RunE:  runConfigInit,
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current mkcd configuration settings.`,
	RunE:  runConfigShow,
}

// configEditCmd represents the config edit command
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration in editor",
	Long:  `Open the configuration file in your default editor.`,
	RunE:  runConfigEdit,
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long:  `Validate the current configuration for errors and inconsistencies.`,
	RunE:  runConfigValidate,
}

// configResetCmd represents the config reset command
var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Reset the configuration file to default settings.`,
	RunE:  runConfigReset,
}

func init() {
	rootCmd.AddCommand(configCmd)
	
	// Add subcommands
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configResetCmd)
}

// runConfigInit initializes the configuration file
func runConfigInit(cmd *cobra.Command, args []string) error {
	outputMgr := utils.NewOutputManager(true, true, true, quiet, verbose, debug)

	// Get config path
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		if !force {
			return fmt.Errorf("configuration file already exists at %s (use --force to overwrite)", configPath)
		}
		outputMgr.Warning("Overwriting existing configuration file")
	}

	// Create default configuration
	cfg := config.DefaultConfig()

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	outputMgr.Success(fmt.Sprintf("Configuration initialized at %s", configPath))
	outputMgr.Info("You can now edit the configuration with: mkcd config edit")

	return nil
}

// runConfigShow displays the current configuration
func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	outputMgr := utils.NewOutputManager(
		cfg.Output.Colors,
		cfg.Output.Icons,
		cfg.Output.ProgressBars,
		quiet,
		verbose,
		debug,
	)

	outputMgr.Header("mkcd Configuration")

	// Core settings
	outputMgr.Section("Core Settings")
	coreSettings := []string{
		fmt.Sprintf("Default Profile: %s", cfg.Core.DefaultProfile),
		fmt.Sprintf("Editor: %s", cfg.Core.Editor),
		fmt.Sprintf("Shell Integration: %t", cfg.Core.ShellIntegration),
		fmt.Sprintf("History Limit: %d", cfg.Core.HistoryLimit),
		fmt.Sprintf("Backup Enabled: %t", cfg.Core.BackupEnabled),
		fmt.Sprintf("Temp Directory: %s", cfg.Core.TempDir),
	}
	outputMgr.List(coreSettings)

	// Git settings
	outputMgr.Section("Git Settings")
	gitSettings := []string{
		fmt.Sprintf("Auto Init: %t", cfg.Git.AutoInit),
		fmt.Sprintf("Default Branch: %s", cfg.Git.DefaultBranch),
		fmt.Sprintf("User Name: %s", cfg.Git.UserName),
		fmt.Sprintf("User Email: %s", cfg.Git.UserEmail),
		fmt.Sprintf("Default Remote Name: %s", cfg.Git.DefaultRemoteName),
	}
	outputMgr.List(gitSettings)

	// Template settings
	outputMgr.Section("Template Settings")
	templateSettings := []string{
		fmt.Sprintf("Directory: %s", cfg.Templates.Directory),
		fmt.Sprintf("Auto Update: %t", cfg.Templates.AutoUpdate),
	}
	outputMgr.List(templateSettings)

	// Safety settings
	outputMgr.Section("Safety Settings")
	safetySettings := []string{
		fmt.Sprintf("Confirm Overwrites: %t", cfg.Safety.ConfirmOverwrites),
		fmt.Sprintf("Confirm Deletes: %t", cfg.Safety.ConfirmDeletes),
		fmt.Sprintf("Max Depth: %d", cfg.Safety.MaxDepth),
		fmt.Sprintf("Forbidden Paths: %v", cfg.Safety.ForbiddenPaths),
	}
	outputMgr.List(safetySettings)

	// Output settings
	outputMgr.Section("Output Settings")
	outputSettings := []string{
		fmt.Sprintf("Colors: %t", cfg.Output.Colors),
		fmt.Sprintf("Icons: %t", cfg.Output.Icons),
		fmt.Sprintf("Progress Bars: %t", cfg.Output.ProgressBars),
	}
	outputMgr.List(outputSettings)

	// Profiles
	outputMgr.Section("Profiles")
	if len(cfg.Profiles) == 0 {
		outputMgr.Info("No profiles configured")
	} else {
		profileList := []string{}
		for name := range cfg.Profiles {
			if name == cfg.Core.DefaultProfile {
				profileList = append(profileList, fmt.Sprintf("%s (default)", name))
			} else {
				profileList = append(profileList, name)
			}
		}
		outputMgr.List(profileList)
	}

	// Show config file location
	configPath := cfgFile
	if configPath == "" {
		configPath, _ = config.GetConfigPath()
	}
	outputMgr.Info(fmt.Sprintf("Configuration file: %s", configPath))

	return nil
}

// runConfigEdit opens the configuration file in an editor
func runConfigEdit(cmd *cobra.Command, args []string) error {
	outputMgr := utils.NewOutputManager(true, true, true, quiet, verbose, debug)

	// Get config path
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		outputMgr.Warning("Configuration file does not exist, creating with defaults...")
		cfg := config.DefaultConfig()
		if err := cfg.Save(configPath); err != nil {
			return fmt.Errorf("failed to create configuration file: %w", err)
		}
	}

	// Get editor
	editorCmd := os.Getenv("EDITOR")
	if editorCmd == "" {
		editorCmd = os.Getenv("VISUAL")
	}
	if editorCmd == "" {
		editorCmd = "vi" // fallback
	}

	outputMgr.Info(fmt.Sprintf("Opening configuration file in %s...", editorCmd))

	// Launch editor
	execCmd := exec.Command(editorCmd, configPath)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	outputMgr.Success("Configuration file edited")
	outputMgr.Info("Validating configuration...")

	// Validate the edited configuration
	if _, err := config.Load(configPath); err != nil {
		outputMgr.Error(fmt.Sprintf("Configuration validation failed: %v", err))
		outputMgr.Info("Please fix the configuration file and try again")
		return fmt.Errorf("invalid configuration")
	}

	outputMgr.Success("Configuration is valid")
	return nil
}

// runConfigValidate validates the current configuration
func runConfigValidate(cmd *cobra.Command, args []string) error {
	outputMgr := utils.NewOutputManager(true, true, true, quiet, verbose, debug)

	// Get config path
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
	}

	outputMgr.Info(fmt.Sprintf("Validating configuration: %s", configPath))

	// Load and validate configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		outputMgr.Error(fmt.Sprintf("Configuration validation failed: %v", err))
		return fmt.Errorf("invalid configuration")
	}

	// Additional validation checks
	validationErrors := []string{}

	// Check if default profile exists
	if cfg.Core.DefaultProfile != "" {
		if _, exists := cfg.Profiles[cfg.Core.DefaultProfile]; !exists {
			validationErrors = append(validationErrors, fmt.Sprintf("Default profile '%s' does not exist", cfg.Core.DefaultProfile))
		}
	}

	// Check template directory
	if cfg.Templates.Directory != "" {
		if _, err := os.Stat(cfg.Templates.Directory); os.IsNotExist(err) {
			validationErrors = append(validationErrors, fmt.Sprintf("Template directory does not exist: %s", cfg.Templates.Directory))
		}
	}

	// Check temp directory
	if cfg.Core.TempDir != "" {
		if _, err := os.Stat(cfg.Core.TempDir); os.IsNotExist(err) {
			validationErrors = append(validationErrors, fmt.Sprintf("Temp directory does not exist: %s", cfg.Core.TempDir))
		}
	}

	if len(validationErrors) > 0 {
		outputMgr.Warning("Configuration has warnings:")
		outputMgr.List(validationErrors)
	} else {
		outputMgr.Success("Configuration is valid")
	}

	return nil
}

// runConfigReset resets the configuration to defaults
func runConfigReset(cmd *cobra.Command, args []string) error {
	outputMgr := utils.NewOutputManager(true, true, true, quiet, verbose, debug)

	// Get config path
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
	}

	// Confirm reset unless force is used
	if !force {
		confirmed, err := outputMgr.Confirm("Reset configuration to defaults? This will overwrite your current settings.", false)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			outputMgr.Info("Reset cancelled")
			return nil
		}
	}

	// Create default configuration
	cfg := config.DefaultConfig()

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	outputMgr.Success(fmt.Sprintf("Configuration reset to defaults: %s", configPath))
	return nil
}
