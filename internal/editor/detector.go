/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

// Package editor provides editor detection and launching functionality for mkcd.
// It can auto-detect available editors and launch them with the specified directory.
package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
)

// EditorInfo contains information about an editor
type EditorInfo struct {
	Name        string   // Display name
	Command     string   // Executable command
	Args        []string // Default arguments
	Description string   // Description
	Priority    int      // Priority for auto-detection (higher = preferred)
}

// EditorDetector handles editor detection and launching
type EditorDetector struct {
	DryRun  bool
	Verbose bool
}

// NewEditorDetector creates a new EditorDetector instance
func NewEditorDetector(dryRun, verbose bool) *EditorDetector {
	return &EditorDetector{
		DryRun:  dryRun,
		Verbose: verbose,
	}
}

// GetAvailableEditors returns a list of available editors on the system
func (ed *EditorDetector) GetAvailableEditors() []EditorInfo {
	editors := []EditorInfo{
		// IDEs and Advanced Editors (highest priority)
		{
			Name:        "Visual Studio Code",
			Command:     "code",
			Args:        []string{},
			Description: "Microsoft Visual Studio Code",
			Priority:    100,
		},
		{
			Name:        "VSCode Insiders",
			Command:     "code-insiders",
			Args:        []string{},
			Description: "Visual Studio Code Insiders",
			Priority:    95,
		},
		{
			Name:        "Cursor",
			Command:     "cursor",
			Args:        []string{},
			Description: "Cursor AI Editor",
			Priority:    90,
		},
		{
			Name:        "Sublime Text",
			Command:     "subl",
			Args:        []string{},
			Description: "Sublime Text",
			Priority:    85,
		},
		{
			Name:        "Atom",
			Command:     "atom",
			Args:        []string{},
			Description: "GitHub Atom",
			Priority:    80,
		},
		{
			Name:        "WebStorm",
			Command:     "webstorm",
			Args:        []string{},
			Description: "JetBrains WebStorm",
			Priority:    75,
		},
		{
			Name:        "IntelliJ IDEA",
			Command:     "idea",
			Args:        []string{},
			Description: "JetBrains IntelliJ IDEA",
			Priority:    75,
		},
		{
			Name:        "GoLand",
			Command:     "goland",
			Args:        []string{},
			Description: "JetBrains GoLand",
			Priority:    75,
		},
		{
			Name:        "PyCharm",
			Command:     "pycharm",
			Args:        []string{},
			Description: "JetBrains PyCharm",
			Priority:    75,
		},

		// Terminal Editors (medium priority)
		{
			Name:        "Neovim",
			Command:     "nvim",
			Args:        []string{},
			Description: "Neovim",
			Priority:    60,
		},
		{
			Name:        "Vim",
			Command:     "vim",
			Args:        []string{},
			Description: "Vim",
			Priority:    55,
		},
		{
			Name:        "Emacs",
			Command:     "emacs",
			Args:        []string{},
			Description: "GNU Emacs",
			Priority:    50,
		},
		{
			Name:        "Nano",
			Command:     "nano",
			Args:        []string{},
			Description: "GNU Nano",
			Priority:    30,
		},

		// Platform-specific editors
		{
			Name:        "TextEdit",
			Command:     "open",
			Args:        []string{"-a", "TextEdit"},
			Description: "macOS TextEdit",
			Priority:    20,
		},
		{
			Name:        "Notepad",
			Command:     "notepad",
			Args:        []string{},
			Description: "Windows Notepad",
			Priority:    10,
		},
	}

	// Filter editors based on platform
	filteredEditors := []EditorInfo{}
	for _, editor := range editors {
		if ed.isEditorAvailable(editor) {
			filteredEditors = append(filteredEditors, editor)
		}
	}

	// Sort by priority (highest first)
	for i := 0; i < len(filteredEditors)-1; i++ {
		for j := i + 1; j < len(filteredEditors); j++ {
			if filteredEditors[i].Priority < filteredEditors[j].Priority {
				filteredEditors[i], filteredEditors[j] = filteredEditors[j], filteredEditors[i]
			}
		}
	}

	return filteredEditors
}

// isEditorAvailable checks if an editor is available on the system
func (ed *EditorDetector) isEditorAvailable(editor EditorInfo) bool {
	// Platform-specific filtering
	switch runtime.GOOS {
	case "darwin":
		// macOS
		if editor.Command == "notepad" {
			return false
		}
	case "windows":
		// Windows
		if editor.Command == "open" && len(editor.Args) > 0 && editor.Args[0] == "-a" {
			return false
		}
	case "linux":
		// Linux
		if editor.Command == "notepad" || (editor.Command == "open" && len(editor.Args) > 0 && editor.Args[0] == "-a") {
			return false
		}
	}

	// Check if command exists
	_, err := exec.LookPath(editor.Command)
	return err == nil
}

// DetectEditor automatically detects the best available editor
func (ed *EditorDetector) DetectEditor() (*EditorInfo, error) {
	// First, check environment variables
	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		if ed.Verbose {
			pterm.Debug.Printf("Using editor from EDITOR environment variable: %s", envEditor)
		}
		return &EditorInfo{
			Name:        "Environment Editor",
			Command:     envEditor,
			Args:        []string{},
			Description: "Editor from EDITOR environment variable",
			Priority:    1000, // Highest priority
		}, nil
	}

	if envEditor := os.Getenv("VISUAL"); envEditor != "" {
		if ed.Verbose {
			pterm.Debug.Printf("Using editor from VISUAL environment variable: %s", envEditor)
		}
		return &EditorInfo{
			Name:        "Visual Editor",
			Command:     envEditor,
			Args:        []string{},
			Description: "Editor from VISUAL environment variable",
			Priority:    999,
		}, nil
	}

	// Get available editors
	editors := ed.GetAvailableEditors()
	if len(editors) == 0 {
		return nil, fmt.Errorf("no editors found on the system")
	}

	// Return the highest priority editor
	bestEditor := editors[0]
	if ed.Verbose {
		pterm.Debug.Printf("Auto-detected editor: %s (%s)", bestEditor.Name, bestEditor.Command)
	}

	return &bestEditor, nil
}

// LaunchEditor launches the specified editor with the given path
func (ed *EditorDetector) LaunchEditor(editor *EditorInfo, path string) error {
	if ed.DryRun {
		pterm.Info.Printf("[DRY RUN] Would launch %s with path: %s", editor.Name, path)
		return nil
	}

	// Ensure path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Prepare command arguments
	args := append(editor.Args, absPath)

	if ed.Verbose {
		pterm.Debug.Printf("Launching editor: %s %s", editor.Command, strings.Join(args, " "))
	}

	// Execute command
	cmd := exec.Command(editor.Command, args...)
	
	// For GUI editors, we typically want to start them in the background
	if ed.isGUIEditor(editor) {
		// Start in background
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start editor %s: %w", editor.Name, err)
		}
		pterm.Success.Printf("Launched %s with path: %s", editor.Name, absPath)
	} else {
		// For terminal editors, run in foreground
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("editor %s exited with error: %w", editor.Name, err)
		}
	}

	return nil
}

// isGUIEditor determines if an editor is a GUI application
func (ed *EditorDetector) isGUIEditor(editor *EditorInfo) bool {
	guiEditors := []string{
		"code", "code-insiders", "cursor", "subl", "atom",
		"webstorm", "idea", "goland", "pycharm", "open",
	}

	for _, gui := range guiEditors {
		if editor.Command == gui {
			return true
		}
	}

	return false
}

// LaunchWithAutoDetection automatically detects and launches the best available editor
func (ed *EditorDetector) LaunchWithAutoDetection(path string) error {
	editor, err := ed.DetectEditor()
	if err != nil {
		return fmt.Errorf("failed to detect editor: %w", err)
	}

	return ed.LaunchEditor(editor, path)
}

// LaunchSpecificEditor launches a specific editor by name or command
func (ed *EditorDetector) LaunchSpecificEditor(editorName, path string) error {
	// First, try to find by name
	editors := ed.GetAvailableEditors()
	for _, editor := range editors {
		if strings.EqualFold(editor.Name, editorName) || 
		   strings.EqualFold(editor.Command, editorName) {
			return ed.LaunchEditor(&editor, path)
		}
	}

	// If not found, try to use as direct command
	if _, err := exec.LookPath(editorName); err == nil {
		customEditor := &EditorInfo{
			Name:        "Custom Editor",
			Command:     editorName,
			Args:        []string{},
			Description: "Custom editor command",
			Priority:    0,
		}
		return ed.LaunchEditor(customEditor, path)
	}

	return fmt.Errorf("editor '%s' not found", editorName)
}

// ListAvailableEditors returns a formatted list of available editors
func (ed *EditorDetector) ListAvailableEditors() []string {
	editors := ed.GetAvailableEditors()
	result := make([]string, len(editors))
	
	for i, editor := range editors {
		result[i] = fmt.Sprintf("%s (%s) - %s", editor.Name, editor.Command, editor.Description)
	}
	
	return result
}
