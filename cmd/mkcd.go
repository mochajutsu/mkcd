/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mochajutsu/mkcd/internal/config"
	"github.com/mochajutsu/mkcd/internal/editor"
	"github.com/mochajutsu/mkcd/internal/files"
	"github.com/mochajutsu/mkcd/internal/git"
	"github.com/mochajutsu/mkcd/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// Command-specific flags for mkcd
var (
	// Workspace setup flags
	gitInit      bool
	gitRemote    string
	template     string
	editorName   string
	editorFlag   bool

	// File creation flags
	touchFiles  []string
	readme      bool
	gitignore   string
	license     string

	// Advanced options
	mode       string
	parentMode string
	symlink    string
	temp       bool
	expire     string
)

// mkcdCmd represents the mkcd command
var mkcdCmd = &cobra.Command{
	Use:   "mkcd <directory>",
	Short: "Create directory and prepare workspace",
	Long: `Create a directory and prepare it for immediate use with optional workspace initialization.

The mkcd command creates a directory (with parent directories as needed) and can optionally:
• Initialize a Git repository with remote setup
• Apply project templates for different languages/frameworks  
• Open the directory in your preferred editor
• Generate common files (README, .gitignore, LICENSE)
• Set up symbolic links or temporary directories

Examples:
  mkcd myproject                           # Basic directory creation
  mkcd myproject --git                     # Create with Git repository
  mkcd myproject --git --remote origin    # Create with Git and remote
  mkcd myproject --template nodejs        # Create using Node.js template
  mkcd myproject --profile dev             # Create using 'dev' profile
  mkcd myproject --editor                  # Create and open in editor
  mkcd myproject --readme --gitignore go   # Create with README and Go .gitignore`,
	Args: cobra.ExactArgs(1),
	RunE: runMkcd,
}

func init() {
	rootCmd.AddCommand(mkcdCmd)

	// Workspace setup flags
	mkcdCmd.Flags().BoolVar(&gitInit, "git", false, "initialize git repository")
	mkcdCmd.Flags().StringVar(&gitRemote, "git-remote", "", "add remote origin URL")
	mkcdCmd.Flags().StringVarP(&template, "template", "t", "", "apply project template")
	mkcdCmd.Flags().StringVarP(&editorName, "editor", "e", "", "open in editor (specify editor or leave empty for auto-detect)")
	mkcdCmd.Flags().BoolVar(&editorFlag, "open-editor", false, "open in editor (auto-detect)")

	// File creation flags
	mkcdCmd.Flags().StringSliceVar(&touchFiles, "touch", []string{}, "create file(s) in directory")
	mkcdCmd.Flags().BoolVar(&readme, "readme", false, "generate README.md")
	mkcdCmd.Flags().StringVar(&gitignore, "gitignore", "", "generate .gitignore for language/framework")
	mkcdCmd.Flags().StringVar(&license, "license", "", "generate LICENSE file")

	// Advanced options
	mkcdCmd.Flags().StringVar(&mode, "mode", "", "set directory permissions (e.g., 755)")
	mkcdCmd.Flags().StringVar(&parentMode, "parent-mode", "", "set parent directory permissions")
	mkcdCmd.Flags().StringVarP(&symlink, "symlink", "s", "", "create as symlink to target")
	mkcdCmd.Flags().BoolVar(&temp, "temp", false, "create in temporary directory")
	mkcdCmd.Flags().StringVar(&expire, "expire", "", "auto-delete after duration (1h, 30m, etc.)")

	// Mark some flags as mutually exclusive
	mkcdCmd.MarkFlagsMutuallyExclusive("symlink", "temp")
	mkcdCmd.MarkFlagsMutuallyExclusive("git-remote", "symlink")
}

// runMkcd executes the main mkcd functionality
func runMkcd(cmd *cobra.Command, args []string) error {
	dirName := args[0]

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get profile configuration if specified
	var profileConfig config.ProfileConfig
	if profile != "" {
		profileConfig, err = cfg.GetProfile(profile)
		if err != nil {
			return fmt.Errorf("failed to get profile: %w", err)
		}
	} else {
		// Use default profile
		profileConfig, err = cfg.GetProfile(cfg.Core.DefaultProfile)
		if err != nil {
			pterm.Debug.Printf("No default profile found, using empty profile")
			profileConfig = config.ProfileConfig{}
		}
	}

	// Create output manager
	outputMgr := utils.NewOutputManager(
		cfg.Output.Colors,
		cfg.Output.Icons,
		cfg.Output.ProgressBars,
		quiet,
		verbose,
		debug,
	)

	// Create filesystem operations manager
	fsOps := utils.NewFileSystemOperations(dryRun, backup || cfg.Core.BackupEnabled)

	// Create path validator
	pathValidator := utils.NewPathValidator(cfg.Safety.ForbiddenPaths, cfg.Safety.MaxDepth)

	// Merge command flags with profile settings
	mergedConfig := mergeConfigWithFlags(profileConfig)

	// Execute the mkcd operation
	return executeMkcd(dirName, cfg, mergedConfig, outputMgr, fsOps, pathValidator)
}

// mergeConfigWithFlags merges profile configuration with command-line flags
func mergeConfigWithFlags(profileConfig config.ProfileConfig) MkcdConfig {
	merged := MkcdConfig{
		Git:       gitInit || profileConfig.Git,
		GitRemote: gitRemote,
		Template:  template,
		Editor:    editorFlag || profileConfig.Editor || (editorName != ""),
		Readme:    readme || profileConfig.Readme,
		Gitignore: gitignore,
		License:   license,
		Touch:     touchFiles,
		Mode:      mode,
		ParentMode: parentMode,
		Symlink:   symlink,
		Temp:      temp,
		Expire:    expire,
	}

	// Use profile values if command flags are empty
	if merged.Template == "" {
		merged.Template = profileConfig.Template
	}
	if merged.Gitignore == "" {
		merged.Gitignore = profileConfig.Gitignore
	}
	if merged.License == "" {
		merged.License = profileConfig.License
	}
	if len(merged.Touch) == 0 {
		merged.Touch = profileConfig.Touch
	}

	return merged
}

// MkcdConfig represents the merged configuration for mkcd operation
type MkcdConfig struct {
	Git        bool
	GitRemote  string
	Template   string
	Editor     bool
	Readme     bool
	Gitignore  string
	License    string
	Touch      []string
	Mode       string
	ParentMode string
	Symlink    string
	Temp       bool
	Expire     string
}

// executeMkcd performs the actual mkcd operation
func executeMkcd(dirName string, cfg *config.Config, mkcdConfig MkcdConfig, outputMgr *utils.OutputManager, fsOps *utils.FileSystemOperations, pathValidator *utils.PathValidator) error {
	// Determine target path
	targetPath, err := determineTargetPath(dirName, mkcdConfig, cfg)
	if err != nil {
		return fmt.Errorf("failed to determine target path: %w", err)
	}

	// Validate path
	if err := pathValidator.ValidatePath(targetPath); err != nil {
		if !force {
			return fmt.Errorf("path validation failed: %w", err)
		}
		outputMgr.Warning(fmt.Sprintf("Path validation failed but continuing due to --force: %v", err))
	}

	// Check for interactive confirmation if needed
	if interactive && !dryRun {
		confirmed, err := outputMgr.Confirm(fmt.Sprintf("Create directory %s?", targetPath), true)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			outputMgr.Info("Operation cancelled by user")
			return nil
		}
	}

	// Create directory structure
	if err := createDirectoryStructure(targetPath, mkcdConfig, fsOps, outputMgr); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate files if requested
	if err := generateProjectFiles(targetPath, mkcdConfig, cfg, fsOps, outputMgr); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	// Initialize Git repository if requested
	if mkcdConfig.Git {
		gitMgr := git.NewGitManager(dryRun, verbose, cfg.Git.UserName, cfg.Git.UserEmail)
		if err := gitMgr.InitRepository(targetPath, cfg.Git.DefaultBranch); err != nil {
			return fmt.Errorf("failed to initialize Git repository: %w", err)
		}

		// Add remote if specified
		if mkcdConfig.GitRemote != "" {
			if err := gitMgr.AddRemote(targetPath, cfg.Git.DefaultRemoteName, mkcdConfig.GitRemote); err != nil {
				return fmt.Errorf("failed to add Git remote: %w", err)
			}
		}

		// Create initial commit if there are files
		if err := gitMgr.CreateInitialCommit(targetPath, "Initial commit"); err != nil {
			outputMgr.Warning(fmt.Sprintf("Failed to create initial commit: %v", err))
		}
	}

	// Open in editor if requested
	if mkcdConfig.Editor {
		if err := openInEditor(targetPath, mkcdConfig, outputMgr); err != nil {
			outputMgr.Warning(fmt.Sprintf("Failed to open in editor: %v", err))
		}
	}

	// Generate shell script for cd operation
	if err := generateShellScript(targetPath, outputMgr); err != nil {
		return fmt.Errorf("failed to generate shell script: %w", err)
	}

	return nil
}

// determineTargetPath determines the final target path based on configuration
func determineTargetPath(dirName string, mkcdConfig MkcdConfig, cfg *config.Config) (string, error) {
	var targetPath string

	if mkcdConfig.Temp {
		// Create in temporary directory
		tempDir := cfg.Core.TempDir
		if tempDir == "" {
			tempDir = os.TempDir()
		}
		targetPath = filepath.Join(tempDir, dirName)
	} else {
		// Use current directory as base
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		targetPath = filepath.Join(cwd, dirName)
	}

	// Get absolute path
	absPath, err := utils.GetAbsolutePath(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return absPath, nil
}

// createDirectoryStructure creates the directory and any required structure
func createDirectoryStructure(targetPath string, mkcdConfig MkcdConfig, fsOps *utils.FileSystemOperations, outputMgr *utils.OutputManager) error {
	// Determine directory mode
	dirMode := os.FileMode(0755) // Default
	if mkcdConfig.Mode != "" {
		// Parse mode from string (e.g., "755")
		// This is a simplified implementation
		outputMgr.Debug(fmt.Sprintf("Custom mode specified: %s", mkcdConfig.Mode))
	}

	// Handle symlink creation
	if mkcdConfig.Symlink != "" {
		return fsOps.CreateSymlink(mkcdConfig.Symlink, targetPath)
	}

	// Create directory
	if err := fsOps.CreateDirectory(targetPath, dirMode); err != nil {
		return err
	}

	// Create files specified in touch
	for _, fileName := range mkcdConfig.Touch {
		filePath := filepath.Join(targetPath, fileName)
		if err := fsOps.CreateFile(filePath, "", 0644); err != nil {
			outputMgr.Warning(fmt.Sprintf("Failed to create file %s: %v", fileName, err))
		}
	}

	return nil
}

// generateProjectFiles generates project files based on configuration
func generateProjectFiles(targetPath string, mkcdConfig MkcdConfig, cfg *config.Config, fsOps *utils.FileSystemOperations, outputMgr *utils.OutputManager) error {
	// Create file generator
	fileGen := files.NewFileGenerator(fsOps, dryRun, verbose)

	// Create generation context
	ctx := files.NewGenerationContext(targetPath)
	ctx.Author = cfg.Git.UserName
	ctx.Email = cfg.Git.UserEmail

	// Generate README if requested
	if mkcdConfig.Readme {
		if err := fileGen.GenerateReadme(ctx); err != nil {
			return fmt.Errorf("failed to generate README: %w", err)
		}
	}

	// Generate .gitignore if requested
	if mkcdConfig.Gitignore != "" {
		if err := fileGen.GenerateGitignore(ctx, mkcdConfig.Gitignore); err != nil {
			return fmt.Errorf("failed to generate .gitignore: %w", err)
		}
	}

	// Generate LICENSE if requested
	if mkcdConfig.License != "" {
		if err := fileGen.GenerateLicense(ctx, mkcdConfig.License); err != nil {
			return fmt.Errorf("failed to generate LICENSE: %w", err)
		}
	}

	return nil
}

// openInEditor opens the project directory in an editor
func openInEditor(targetPath string, mkcdConfig MkcdConfig, outputMgr *utils.OutputManager) error {
	editorLauncher := editor.NewEditorLauncher(dryRun, verbose)

	options := editor.LaunchOptions{
		EditorName:    editorName,
		Path:          targetPath,
		Wait:          false, // Don't wait for editor to close
		CreateMissing: dryRun, // In dry-run mode, allow "creating" missing paths
	}

	return editorLauncher.Launch(options)
}

// generateShellScript generates the shell script for cd operation
func generateShellScript(targetPath string, outputMgr *utils.OutputManager) error {
	// This is where we output the shell script that the wrapper function will eval
	// The actual shell integration will be implemented in the shell package

	if !quiet {
		outputMgr.Success(fmt.Sprintf("Directory created: %s", targetPath))
		outputMgr.Info("To change to the directory, run: cd " + targetPath)
	}

	// For now, just output the cd command
	// In the full implementation, this would generate proper shell scripts
	fmt.Printf("cd %s\n", targetPath)

	return nil
}
