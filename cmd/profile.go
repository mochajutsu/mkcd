/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mochajutsu/mkcd/internal/config"
	"github.com/mochajutsu/mkcd/internal/utils"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long: `Manage configuration profiles for mkcd.

Profiles allow you to save and reuse common configurations for different
types of projects. Each profile can specify default settings for git
initialization, editor preferences, file generation, and more.

Examples:
  mkcd profile list                    # List all profiles
  mkcd profile show dev                # Show 'dev' profile details
  mkcd profile create myprofile        # Create new profile interactively
  mkcd profile edit dev                # Edit 'dev' profile in $EDITOR
  mkcd profile delete myprofile        # Delete 'myprofile'
  mkcd profile copy dev mydev          # Copy 'dev' profile to 'mydev'`,
}

// profileListCmd represents the profile list command
var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available profiles",
	Long:  `List all available configuration profiles with their descriptions.`,
	RunE:  runProfileList,
}

// profileShowCmd represents the profile show command
var profileShowCmd = &cobra.Command{
	Use:   "show <profile-name>",
	Short: "Show profile configuration",
	Long:  `Show the detailed configuration of a specific profile.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileShow,
}

// profileCreateCmd represents the profile create command
var profileCreateCmd = &cobra.Command{
	Use:   "create <profile-name>",
	Short: "Create a new profile",
	Long:  `Create a new configuration profile interactively.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileCreate,
}

// profileEditCmd represents the profile edit command
var profileEditCmd = &cobra.Command{
	Use:   "edit <profile-name>",
	Short: "Edit an existing profile",
	Long:  `Edit an existing profile in your default editor.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileEdit,
}

// profileDeleteCmd represents the profile delete command
var profileDeleteCmd = &cobra.Command{
	Use:   "delete <profile-name>",
	Short: "Delete a profile",
	Long:  `Delete an existing configuration profile.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileDelete,
}

// profileCopyCmd represents the profile copy command
var profileCopyCmd = &cobra.Command{
	Use:   "copy <source-profile> <destination-profile>",
	Short: "Copy a profile",
	Long:  `Copy an existing profile to a new profile name.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runProfileCopy,
}

func init() {
	rootCmd.AddCommand(profileCmd)
	
	// Add subcommands
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileEditCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	profileCmd.AddCommand(profileCopyCmd)
}

// runProfileList lists all available profiles
func runProfileList(cmd *cobra.Command, args []string) error {
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

	if len(cfg.Profiles) == 0 {
		outputMgr.Info("No profiles found")
		return nil
	}

	outputMgr.Header("Available Profiles")

	// Prepare table data
	headers := []string{"Name", "Git", "Editor", "Template", "Description"}
	rows := [][]string{}

	for name, profile := range cfg.Profiles {
		gitStatus := "No"
		if profile.Git {
			gitStatus = "Yes"
		}

		editorStatus := "No"
		if profile.Editor {
			editorStatus = "Yes"
		}

		template := profile.Template
		if template == "" {
			template = "-"
		}

		description := generateProfileDescription(profile)

		// Mark default profile
		displayName := name
		if name == cfg.Core.DefaultProfile {
			displayName = name + " (default)"
		}

		rows = append(rows, []string{displayName, gitStatus, editorStatus, template, description})
	}

	outputMgr.Table(headers, rows)
	return nil
}

// runProfileShow shows details of a specific profile
func runProfileShow(cmd *cobra.Command, args []string) error {
	profileName := args[0]

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

	profile, exists := cfg.Profiles[profileName]
	if !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	outputMgr.Header(fmt.Sprintf("Profile: %s", profileName))

	// Show profile details
	details := []string{
		fmt.Sprintf("Git initialization: %t", profile.Git),
		fmt.Sprintf("Editor integration: %t", profile.Editor),
		fmt.Sprintf("Generate README: %t", profile.Readme),
	}

	if profile.Template != "" {
		details = append(details, fmt.Sprintf("Template: %s", profile.Template))
	}

	if profile.Gitignore != "" {
		details = append(details, fmt.Sprintf("Gitignore type: %s", profile.Gitignore))
	}

	if profile.License != "" {
		details = append(details, fmt.Sprintf("License: %s", profile.License))
	}

	if len(profile.Touch) > 0 {
		details = append(details, fmt.Sprintf("Touch files: %s", strings.Join(profile.Touch, ", ")))
	}

	outputMgr.List(details)

	// Show if this is the default profile
	if profileName == cfg.Core.DefaultProfile {
		outputMgr.Info("This is the default profile")
	}

	return nil
}

// runProfileCreate creates a new profile interactively
func runProfileCreate(cmd *cobra.Command, args []string) error {
	profileName := args[0]

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

	// Check if profile already exists
	if _, exists := cfg.Profiles[profileName]; exists {
		if !force {
			return fmt.Errorf("profile '%s' already exists (use --force to overwrite)", profileName)
		}
		outputMgr.Warning(fmt.Sprintf("Overwriting existing profile '%s'", profileName))
	}

	outputMgr.Header(fmt.Sprintf("Creating Profile: %s", profileName))

	// Interactive profile creation
	profile := config.ProfileConfig{}

	// Git initialization
	gitInit, err := outputMgr.Confirm("Initialize Git repository by default?", false)
	if err != nil {
		return fmt.Errorf("failed to get Git preference: %w", err)
	}
	profile.Git = gitInit

	// Editor integration
	editorOpen, err := outputMgr.Confirm("Open in editor by default?", false)
	if err != nil {
		return fmt.Errorf("failed to get editor preference: %w", err)
	}
	profile.Editor = editorOpen

	// README generation
	readmeGen, err := outputMgr.Confirm("Generate README.md by default?", false)
	if err != nil {
		return fmt.Errorf("failed to get README preference: %w", err)
	}
	profile.Readme = readmeGen

	// Template selection
	templateOptions := []string{"", "basic-dev", "nodejs", "python", "go", "web"}
	template, err := outputMgr.Select("Select default template (or empty for none):", templateOptions)
	if err != nil {
		return fmt.Errorf("failed to get template preference: %w", err)
	}
	profile.Template = template

	// Gitignore type
	gitignoreOptions := []string{"", "general", "go", "node", "python"}
	gitignoreType, err := outputMgr.Select("Select default .gitignore type (or empty for none):", gitignoreOptions)
	if err != nil {
		return fmt.Errorf("failed to get gitignore preference: %w", err)
	}
	profile.Gitignore = gitignoreType

	// License type
	licenseOptions := []string{"", "mit", "apache-2.0"}
	licenseType, err := outputMgr.Select("Select default license (or empty for none):", licenseOptions)
	if err != nil {
		return fmt.Errorf("failed to get license preference: %w", err)
	}
	profile.License = licenseType

	// Touch files
	touchFiles, err := outputMgr.Input("Enter files to create by default (comma-separated, or empty):", "")
	if err != nil {
		return fmt.Errorf("failed to get touch files: %w", err)
	}
	if touchFiles != "" {
		profile.Touch = strings.Split(strings.ReplaceAll(touchFiles, " ", ""), ",")
	}

	// Save profile
	cfg.SetProfile(profileName, profile)

	// Ask if this should be the default profile
	if cfg.Core.DefaultProfile == "" || cfg.Core.DefaultProfile == "default" {
		makeDefault, err := outputMgr.Confirm("Make this the default profile?", false)
		if err != nil {
			return fmt.Errorf("failed to get default preference: %w", err)
		}
		if makeDefault {
			cfg.Core.DefaultProfile = profileName
		}
	}

	// Save configuration
	if err := cfg.Save(cfgFile); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	outputMgr.Success(fmt.Sprintf("Profile '%s' created successfully", profileName))
	return nil
}

// runProfileEdit edits an existing profile in the user's editor
func runProfileEdit(cmd *cobra.Command, args []string) error {
	profileName := args[0]

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

	// Check if profile exists
	if _, exists := cfg.Profiles[profileName]; !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Get config file path
	configPath := cfgFile
	if configPath == "" {
		configPath, err = config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
	}

	// Get editor
	editorCmd := os.Getenv("EDITOR")
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
	outputMgr.Info("Note: Changes will take effect on next mkcd command")

	return nil
}

// runProfileDelete deletes an existing profile
func runProfileDelete(cmd *cobra.Command, args []string) error {
	profileName := args[0]

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

	// Check if profile exists
	if _, exists := cfg.Profiles[profileName]; !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Confirm deletion unless force is used
	if !force {
		confirmed, err := outputMgr.Confirm(fmt.Sprintf("Delete profile '%s'?", profileName), false)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			outputMgr.Info("Deletion cancelled")
			return nil
		}
	}

	// Delete profile
	if err := cfg.DeleteProfile(profileName); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	// Save configuration
	if err := cfg.Save(cfgFile); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	outputMgr.Success(fmt.Sprintf("Profile '%s' deleted successfully", profileName))
	return nil
}

// runProfileCopy copies an existing profile to a new name
func runProfileCopy(cmd *cobra.Command, args []string) error {
	sourceProfile := args[0]
	destProfile := args[1]

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

	// Check if source profile exists
	sourceConfig, exists := cfg.Profiles[sourceProfile]
	if !exists {
		return fmt.Errorf("source profile '%s' not found", sourceProfile)
	}

	// Check if destination profile already exists
	if _, exists := cfg.Profiles[destProfile]; exists {
		if !force {
			return fmt.Errorf("destination profile '%s' already exists (use --force to overwrite)", destProfile)
		}
		outputMgr.Warning(fmt.Sprintf("Overwriting existing profile '%s'", destProfile))
	}

	// Copy profile
	cfg.SetProfile(destProfile, sourceConfig)

	// Save configuration
	if err := cfg.Save(cfgFile); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	outputMgr.Success(fmt.Sprintf("Profile '%s' copied to '%s'", sourceProfile, destProfile))
	return nil
}

// generateProfileDescription generates a brief description of a profile
func generateProfileDescription(profile config.ProfileConfig) string {
	features := []string{}

	if profile.Git {
		features = append(features, "Git")
	}
	if profile.Editor {
		features = append(features, "Editor")
	}
	if profile.Readme {
		features = append(features, "README")
	}
	if profile.Template != "" {
		features = append(features, fmt.Sprintf("Template:%s", profile.Template))
	}

	if len(features) == 0 {
		return "Basic profile"
	}

	return strings.Join(features, ", ")
}
