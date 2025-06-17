/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// EditorLauncher provides high-level editor launching functionality
type EditorLauncher struct {
	detector *EditorDetector
	DryRun   bool
	Verbose  bool
}

// NewEditorLauncher creates a new EditorLauncher instance
func NewEditorLauncher(dryRun, verbose bool) *EditorLauncher {
	return &EditorLauncher{
		detector: NewEditorDetector(dryRun, verbose),
		DryRun:   dryRun,
		Verbose:  verbose,
	}
}

// LaunchOptions contains options for launching an editor
type LaunchOptions struct {
	EditorName    string        // Specific editor to use (empty for auto-detect)
	Path          string        // Path to open
	Wait          bool          // Wait for editor to close
	Timeout       time.Duration // Timeout for waiting
	CreateMissing bool          // Create path if it doesn't exist
	OpenFiles     []string      // Specific files to open within the path
}

// Launch launches an editor with the specified options
func (el *EditorLauncher) Launch(options LaunchOptions) error {
	// Validate and prepare path
	targetPath, err := el.preparePath(options.Path, options.CreateMissing)
	if err != nil {
		return fmt.Errorf("failed to prepare path: %w", err)
	}

	// Determine which editor to use
	var editor *EditorInfo
	if options.EditorName != "" {
		editor, err = el.getSpecificEditor(options.EditorName)
		if err != nil {
			return fmt.Errorf("failed to get specific editor: %w", err)
		}
	} else {
		editor, err = el.detector.DetectEditor()
		if err != nil {
			return fmt.Errorf("failed to detect editor: %w", err)
		}
	}

	// Launch the editor
	return el.launchWithOptions(editor, targetPath, options)
}

// preparePath validates and prepares the target path
func (el *EditorLauncher) preparePath(path string, createMissing bool) (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		if createMissing {
			if el.DryRun {
				pterm.Info.Printf("[DRY RUN] Would create missing path: %s", absPath)
				return absPath, nil
			}

			// Create directory
			if err := os.MkdirAll(absPath, 0755); err != nil {
				return "", fmt.Errorf("failed to create directory: %w", err)
			}
			
			if el.Verbose {
				pterm.Success.Printf("Created directory: %s", absPath)
			}
		} else {
			return "", fmt.Errorf("path does not exist: %s", absPath)
		}
	}

	return absPath, nil
}

// getSpecificEditor gets a specific editor by name or command
func (el *EditorLauncher) getSpecificEditor(editorName string) (*EditorInfo, error) {
	// Get available editors
	editors := el.detector.GetAvailableEditors()
	
	// Search by name or command
	for _, editor := range editors {
		if strings.EqualFold(editor.Name, editorName) ||
		   strings.EqualFold(editor.Command, editorName) ||
		   strings.Contains(strings.ToLower(editor.Name), strings.ToLower(editorName)) {
			return &editor, nil
		}
	}

	// Try as direct command
	if _, err := exec.LookPath(editorName); err == nil {
		return &EditorInfo{
			Name:        fmt.Sprintf("Custom (%s)", editorName),
			Command:     editorName,
			Args:        []string{},
			Description: "Custom editor command",
			Priority:    0,
		}, nil
	}

	return nil, fmt.Errorf("editor '%s' not found", editorName)
}

// launchWithOptions launches the editor with specific options
func (el *EditorLauncher) launchWithOptions(editor *EditorInfo, path string, options LaunchOptions) error {
	if el.DryRun {
		pterm.Info.Printf("[DRY RUN] Would launch %s with path: %s", editor.Name, path)
		if len(options.OpenFiles) > 0 {
			pterm.Info.Printf("[DRY RUN] Would open files: %s", strings.Join(options.OpenFiles, ", "))
		}
		return nil
	}

	// Prepare command arguments
	args := make([]string, len(editor.Args))
	copy(args, editor.Args)

	// Add specific files if provided
	if len(options.OpenFiles) > 0 {
		for _, file := range options.OpenFiles {
			filePath := filepath.Join(path, file)
			args = append(args, filePath)
		}
	} else {
		args = append(args, path)
	}

	if el.Verbose {
		pterm.Debug.Printf("Launching: %s %s", editor.Command, strings.Join(args, " "))
	}

	// Create command
	cmd := exec.Command(editor.Command, args...)
	
	// Set working directory
	cmd.Dir = path

	// Handle different launch modes
	if options.Wait {
		return el.launchAndWait(cmd, editor, options.Timeout)
	} else {
		return el.launchInBackground(cmd, editor)
	}
}

// launchAndWait launches the editor and waits for it to complete
func (el *EditorLauncher) launchAndWait(cmd *exec.Cmd, editor *EditorInfo, timeout time.Duration) error {
	// For terminal editors, connect to current terminal
	if !el.detector.isGUIEditor(editor) {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", editor.Name, err)
	}

	pterm.Info.Printf("Launched %s (PID: %d)", editor.Name, cmd.Process.Pid)

	// Wait with optional timeout
	if timeout > 0 {
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				return fmt.Errorf("%s exited with error: %w", editor.Name, err)
			}
			pterm.Success.Printf("%s completed successfully", editor.Name)
		case <-time.After(timeout):
			if err := cmd.Process.Kill(); err != nil {
				pterm.Warning.Printf("Failed to kill %s after timeout: %v", editor.Name, err)
			}
			return fmt.Errorf("%s timed out after %v", editor.Name, timeout)
		}
	} else {
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("%s exited with error: %w", editor.Name, err)
		}
		pterm.Success.Printf("%s completed successfully", editor.Name)
	}

	return nil
}

// launchInBackground launches the editor in the background
func (el *EditorLauncher) launchInBackground(cmd *exec.Cmd, editor *EditorInfo) error {
	// For GUI editors, detach from terminal
	if el.detector.isGUIEditor(editor) {
		// Redirect outputs to prevent hanging
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
	} else {
		// For terminal editors, connect to current terminal
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", editor.Name, err)
	}

	pterm.Success.Printf("Launched %s in background (PID: %d)", editor.Name, cmd.Process.Pid)

	// For GUI editors, we don't wait
	if el.detector.isGUIEditor(editor) {
		return nil
	}

	// For terminal editors, we still need to wait
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s exited with error: %w", editor.Name, err)
	}

	return nil
}

// GetEditorCommand returns the command that would be executed for an editor
func (el *EditorLauncher) GetEditorCommand(editorName, path string) (string, []string, error) {
	var editor *EditorInfo
	var err error

	if editorName != "" {
		editor, err = el.getSpecificEditor(editorName)
	} else {
		editor, err = el.detector.DetectEditor()
	}

	if err != nil {
		return "", nil, err
	}

	args := make([]string, len(editor.Args))
	copy(args, editor.Args)
	args = append(args, path)

	return editor.Command, args, nil
}

// ValidateEditor checks if an editor is available and working
func (el *EditorLauncher) ValidateEditor(editorName string) error {
	editor, err := el.getSpecificEditor(editorName)
	if err != nil {
		return err
	}

	// Check if command exists
	if _, err := exec.LookPath(editor.Command); err != nil {
		return fmt.Errorf("editor command '%s' not found in PATH", editor.Command)
	}

	// Try to get version or help (non-destructive test)
	versionArgs := []string{"--version"}
	if editor.Command == "vim" || editor.Command == "nvim" {
		versionArgs = []string{"--version"}
	} else if editor.Command == "emacs" {
		versionArgs = []string{"--version"}
	} else if editor.Command == "code" || editor.Command == "code-insiders" {
		versionArgs = []string{"--version"}
	}

	cmd := exec.Command(editor.Command, versionArgs...)
	if err := cmd.Run(); err != nil {
		// Some editors might not support --version, so we just check if they exist
		if el.Verbose {
			pterm.Debug.Printf("Editor %s exists but version check failed (this is often normal)", editor.Name)
		}
	}

	return nil
}

// GetRecommendedEditor returns the recommended editor for a specific project type
func (el *EditorLauncher) GetRecommendedEditor(projectType string) (*EditorInfo, error) {
	editors := el.detector.GetAvailableEditors()
	if len(editors) == 0 {
		return nil, fmt.Errorf("no editors available")
	}

	// Project-specific recommendations
	recommendations := map[string][]string{
		"go":         {"goland", "code", "vim", "nvim"},
		"javascript": {"webstorm", "code", "atom", "subl"},
		"typescript": {"webstorm", "code", "atom", "subl"},
		"python":     {"pycharm", "code", "vim", "nvim"},
		"rust":       {"code", "vim", "nvim", "emacs"},
		"java":       {"idea", "code", "vim", "nvim"},
		"web":        {"webstorm", "code", "atom", "subl"},
		"general":    {"code", "vim", "nvim", "subl"},
	}

	preferredCommands, exists := recommendations[strings.ToLower(projectType)]
	if !exists {
		preferredCommands = recommendations["general"]
	}

	// Find the first available preferred editor
	for _, preferred := range preferredCommands {
		for _, editor := range editors {
			if editor.Command == preferred {
				return &editor, nil
			}
		}
	}

	// Fallback to highest priority available editor
	return &editors[0], nil
}
