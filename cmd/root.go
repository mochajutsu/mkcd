/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/
package cmd

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// Global configuration variables
var (
	cfgFile     string
	profile     string
	dryRun      bool
	verbose     bool
	quiet       bool
	debug       bool
	force       bool
	interactive bool
	backup      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mkcd",
	Short: "A powerful, extensible directory creation and workspace initialization tool",
	Long: `mkcd is an enterprise-level CLI tool that revolutionizes directory creation and navigation.
It transforms the simple concept of "make directory and change into it" into a comprehensive
workspace initialization tool with features like:

• Cross-platform directory creation with shell integration
• Git repository initialization and remote setup
• Project templates and workspace profiles
• Editor integration and file generation
• History tracking with undo functionality
• Batch operations and safety checks

Examples:
  mkcd myproject                    # Create directory and prepare for cd
  mkcd myproject --git --editor     # Create with git repo and open in editor
  mkcd myproject --template nodejs  # Create using Node.js template
  mkcd myproject --profile dev      # Create using 'dev' profile`,
	Version: "1.0.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure pterm based on flags
		if quiet {
			pterm.DisableOutput()
		}
		if debug {
			pterm.EnableDebugMessages()
		}
		if !verbose && !debug {
			pterm.DisableStyling()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		if !quiet {
			pterm.Error.Printf("Command failed: %v\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	// Global persistent flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: ~/.config/mkcd/mkcd.conf)")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "use named profile from config")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "show what would be done without executing")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode with trace information")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "override safety checks")
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "interactive mode for confirmations")
	rootCmd.PersistentFlags().BoolVar(&backup, "backup", false, "backup existing directories before operations")

	// Mark some flags as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "quiet")
}


