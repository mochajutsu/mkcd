/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// OutputManager handles formatted output for the mkcd application
type OutputManager struct {
	Colors       bool
	Icons        bool
	ProgressBars bool
	Quiet        bool
	VerboseMode  bool
	DebugMode    bool
}

// NewOutputManager creates a new OutputManager instance
func NewOutputManager(colors, icons, progressBars, quiet, verbose, debug bool) *OutputManager {
	om := &OutputManager{
		Colors:       colors,
		Icons:        icons,
		ProgressBars: progressBars,
		Quiet:        quiet,
		VerboseMode:  verbose,
		DebugMode:    debug,
	}

	// Configure pterm based on settings
	om.configurePterm()
	return om
}

// configurePterm configures pterm based on output settings
func (om *OutputManager) configurePterm() {
	if om.Quiet {
		pterm.DisableOutput()
		return
	}

	if !om.Colors {
		pterm.DisableColor()
	}

	if om.DebugMode {
		pterm.EnableDebugMessages()
	}

	if !om.VerboseMode && !om.DebugMode {
		pterm.DisableStyling()
	}
}

// Success prints a success message
func (om *OutputManager) Success(message string) {
	if om.Quiet {
		return
	}

	if om.Icons {
		pterm.Success.Println(message)
	} else {
		pterm.Println(pterm.Green(message))
	}
}

// Error prints an error message
func (om *OutputManager) Error(message string) {
	if om.Quiet {
		return
	}

	if om.Icons {
		pterm.Error.Println(message)
	} else {
		pterm.Println(pterm.Red(message))
	}
}

// Warning prints a warning message
func (om *OutputManager) Warning(message string) {
	if om.Quiet {
		return
	}

	if om.Icons {
		pterm.Warning.Println(message)
	} else {
		pterm.Println(pterm.Yellow(message))
	}
}

// Info prints an info message
func (om *OutputManager) Info(message string) {
	if om.Quiet {
		return
	}

	if om.Icons {
		pterm.Info.Println(message)
	} else {
		pterm.Println(pterm.Cyan(message))
	}
}

// Debug prints a debug message
func (om *OutputManager) Debug(message string) {
	if om.Quiet || !om.DebugMode {
		return
	}

	pterm.Debug.Println(message)
}

// Verbose prints a verbose message
func (om *OutputManager) Verbose(message string) {
	if om.Quiet || !om.VerboseMode {
		return
	}

	pterm.Println(pterm.Gray(message))
}

// Print prints a regular message
func (om *OutputManager) Print(message string) {
	if om.Quiet {
		return
	}

	pterm.Println(message)
}

// Printf prints a formatted message
func (om *OutputManager) Printf(format string, args ...interface{}) {
	if om.Quiet {
		return
	}

	pterm.Printf(format, args...)
}

// Header prints a styled header
func (om *OutputManager) Header(title string) {
	if om.Quiet {
		return
	}

	if om.Icons && om.Colors {
		pterm.DefaultHeader.WithFullWidth().Println(title)
	} else {
		om.Print(strings.ToUpper(title))
		om.Print(strings.Repeat("=", len(title)))
	}
}

// Section prints a styled section header
func (om *OutputManager) Section(title string) {
	if om.Quiet {
		return
	}

	if om.Icons && om.Colors {
		pterm.DefaultSection.Println(title)
	} else {
		om.Print(title)
		om.Print(strings.Repeat("-", len(title)))
	}
}

// List prints a bulleted list
func (om *OutputManager) List(items []string) {
	if om.Quiet {
		return
	}

	if om.Icons && om.Colors {
		// Convert strings to BulletListItems
		bulletItems := make([]pterm.BulletListItem, len(items))
		for i, item := range items {
			bulletItems[i] = pterm.BulletListItem{Text: item}
		}
		pterm.DefaultBulletList.WithItems(bulletItems).Render()
	} else {
		for _, item := range items {
			om.Printf("• %s\n", item)
		}
	}
}

// Table prints a table with headers and rows
func (om *OutputManager) Table(headers []string, rows [][]string) {
	if om.Quiet {
		return
	}

	if om.Colors {
		tableData := pterm.TableData{headers}
		for _, row := range rows {
			tableData = append(tableData, row)
		}
		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	} else {
		// Simple text table
		om.printSimpleTable(headers, rows)
	}
}

// printSimpleTable prints a simple text-based table
func (om *OutputManager) printSimpleTable(headers []string, rows [][]string) {
	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print header
	for i, header := range headers {
		om.Printf("%-*s", colWidths[i]+2, header)
	}
	om.Print("")

	// Print separator
	for _, width := range colWidths {
		om.Print(strings.Repeat("-", width+2))
	}
	om.Print("")

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				om.Printf("%-*s", colWidths[i]+2, cell)
			}
		}
		om.Print("")
	}
}

// ProgressBar creates and returns a progress bar
func (om *OutputManager) ProgressBar(title string, total int) *pterm.ProgressbarPrinter {
	if om.Quiet || !om.ProgressBars {
		return nil
	}

	return pterm.DefaultProgressbar.WithTitle(title).WithTotal(total)
}

// Spinner creates and returns a spinner
func (om *OutputManager) Spinner(text string) *pterm.SpinnerPrinter {
	if om.Quiet {
		return nil
	}

	return pterm.DefaultSpinner.WithText(text)
}

// Confirm prompts the user for confirmation
func (om *OutputManager) Confirm(message string, defaultValue bool) (bool, error) {
	if om.Quiet {
		return defaultValue, nil
	}

	prompt := message
	if defaultValue {
		prompt += " [Y/n]"
	} else {
		prompt += " [y/N]"
	}

	result, err := pterm.DefaultInteractiveConfirm.WithDefaultValue(defaultValue).Show(prompt)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to get user confirmation: %w", err)
	}

	return result, nil
}

// Select prompts the user to select from a list of options
func (om *OutputManager) Select(message string, options []string) (string, error) {
	if om.Quiet {
		if len(options) > 0 {
			return options[0], nil
		}
		return "", fmt.Errorf("no options available")
	}

	result, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show(message)
	if err != nil {
		return "", fmt.Errorf("failed to get user selection: %w", err)
	}

	return result, nil
}

// Input prompts the user for text input
func (om *OutputManager) Input(message string, defaultValue string) (string, error) {
	if om.Quiet {
		return defaultValue, nil
	}

	result, err := pterm.DefaultInteractiveTextInput.WithDefaultValue(defaultValue).Show(message)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to get user input: %w", err)
	}

	return result, nil
}

// MultiSelect prompts the user to select multiple options
func (om *OutputManager) MultiSelect(message string, options []string) ([]string, error) {
	if om.Quiet {
		return options, nil
	}

	result, err := pterm.DefaultInteractiveMultiselect.WithOptions(options).Show(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get user selection: %w", err)
	}

	return result, nil
}

// TimedOperation executes an operation with timing information
func (om *OutputManager) TimedOperation(name string, operation func() error) error {
	if om.Quiet {
		return operation()
	}

	start := time.Now()
	
	var spinner *pterm.SpinnerPrinter
	if om.ProgressBars {
		spinner = om.Spinner(fmt.Sprintf("Executing %s...", name))
		spinner.Start()
	} else {
		om.Info(fmt.Sprintf("Starting %s...", name))
	}

	err := operation()
	duration := time.Since(start)

	if spinner != nil {
		if err != nil {
			spinner.Fail(fmt.Sprintf("%s failed after %v", name, duration))
		} else {
			spinner.Success(fmt.Sprintf("%s completed in %v", name, duration))
		}
	} else {
		if err != nil {
			om.Error(fmt.Sprintf("%s failed after %v: %v", name, duration, err))
		} else {
			om.Success(fmt.Sprintf("%s completed in %v", name, duration))
		}
	}

	return err
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// FormatBytes formats bytes in a human-readable way
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
