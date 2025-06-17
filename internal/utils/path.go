/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// PathValidator provides path validation functionality
type PathValidator struct {
	ForbiddenPaths []string
	MaxDepth       int
}

// NewPathValidator creates a new PathValidator instance
func NewPathValidator(forbiddenPaths []string, maxDepth int) *PathValidator {
	return &PathValidator{
		ForbiddenPaths: forbiddenPaths,
		MaxDepth:       maxDepth,
	}
}

// ValidatePath validates a path for safety and correctness
func (pv *PathValidator) ValidatePath(path string) error {
	// Sanitize the path first
	cleanPath, err := SanitizePath(path)
	if err != nil {
		return fmt.Errorf("path sanitization failed: %w", err)
	}

	// Get absolute path for validation
	absPath, err := GetAbsolutePath(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check against forbidden paths
	if err := pv.checkForbiddenPaths(absPath); err != nil {
		return err
	}

	// Check path depth
	if err := pv.checkPathDepth(cleanPath); err != nil {
		return err
	}

	// Check for dangerous characters
	if err := pv.checkDangerousCharacters(cleanPath); err != nil {
		return err
	}

	return nil
}

// checkForbiddenPaths checks if the path is in the forbidden paths list
func (pv *PathValidator) checkForbiddenPaths(absPath string) error {
	for _, forbidden := range pv.ForbiddenPaths {
		// Check if the path is exactly a forbidden path
		if absPath == forbidden {
			return fmt.Errorf("path is forbidden: %s", absPath)
		}

		// Check if the path is under a forbidden directory
		if strings.HasPrefix(absPath, forbidden+string(filepath.Separator)) {
			return fmt.Errorf("path is under forbidden directory %s: %s", forbidden, absPath)
		}
	}
	return nil
}

// checkPathDepth checks if the path depth exceeds the maximum allowed
func (pv *PathValidator) checkPathDepth(path string) error {
	// Count path separators to determine depth
	depth := strings.Count(path, string(filepath.Separator))
	
	// Adjust for relative vs absolute paths
	if filepath.IsAbs(path) {
		depth-- // Don't count the root separator
	}

	if depth > pv.MaxDepth {
		return fmt.Errorf("path depth %d exceeds maximum allowed depth %d: %s", depth, pv.MaxDepth, path)
	}

	return nil
}

// checkDangerousCharacters checks for potentially dangerous characters in the path
func (pv *PathValidator) checkDangerousCharacters(path string) error {
	// Define dangerous patterns
	dangerousPatterns := []struct {
		pattern string
		message string
	}{
		{`\x00`, "null byte"},
		{`[<>:"|?*]`, "invalid filename characters"},
		{`^\s+|\s+$`, "leading or trailing whitespace"},
		{`\.{3,}`, "excessive dots"},
	}

	for _, dp := range dangerousPatterns {
		if matched, _ := regexp.MatchString(dp.pattern, path); matched {
			return fmt.Errorf("path contains %s: %s", dp.message, path)
		}
	}

	return nil
}

// GenerateUniquePath generates a unique path by appending a number if the path already exists
func GenerateUniquePath(basePath string) string {
	if !PathExists(basePath) {
		return basePath
	}

	dir := filepath.Dir(basePath)
	base := filepath.Base(basePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	for i := 1; i < 1000; i++ {
		newName := fmt.Sprintf("%s-%d%s", name, i, ext)
		newPath := filepath.Join(dir, newName)
		if !PathExists(newPath) {
			return newPath
		}
	}

	// Fallback with timestamp if we can't find a unique name
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	newName := fmt.Sprintf("%s-%s%s", name, timestamp, ext)
	return filepath.Join(dir, newName)
}

// ExpandPath expands environment variables and ~ in a path
func ExpandPath(path string) (string, error) {
	// Expand environment variables
	expanded := os.ExpandEnv(path)
	
	// Expand ~ to home directory
	if strings.HasPrefix(expanded, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		expanded = filepath.Join(homeDir, expanded[2:])
	}

	return expanded, nil
}

// RelativePath returns the relative path from base to target
func RelativePath(base, target string) (string, error) {
	// Get absolute paths
	absBase, err := GetAbsolutePath(base)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute base path: %w", err)
	}

	absTarget, err := GetAbsolutePath(target)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute target path: %w", err)
	}

	// Calculate relative path
	relPath, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return "", fmt.Errorf("failed to calculate relative path: %w", err)
	}

	return relPath, nil
}

// JoinPaths safely joins multiple path components
func JoinPaths(paths ...string) string {
	if len(paths) == 0 {
		return ""
	}

	result := paths[0]
	for _, path := range paths[1:] {
		result = filepath.Join(result, path)
	}

	return filepath.Clean(result)
}

// SplitPath splits a path into its directory and filename components
func SplitPath(path string) (dir, filename string) {
	return filepath.Split(path)
}

// GetFileExtension returns the file extension (including the dot)
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// GetBaseName returns the base name of the path without extension
func GetBaseName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// IsSubPath checks if child is a subdirectory of parent
func IsSubPath(parent, child string) (bool, error) {
	// Get absolute paths
	absParent, err := GetAbsolutePath(parent)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute parent path: %w", err)
	}

	absChild, err := GetAbsolutePath(child)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute child path: %w", err)
	}

	// Ensure parent ends with separator for proper comparison
	if !strings.HasSuffix(absParent, string(filepath.Separator)) {
		absParent += string(filepath.Separator)
	}

	return strings.HasPrefix(absChild, absParent), nil
}

// NormalizePath normalizes a path by cleaning it and resolving any symbolic links
func NormalizePath(path string) (string, error) {
	// Clean the path
	cleaned := filepath.Clean(path)

	// Resolve symbolic links
	resolved, err := filepath.EvalSymlinks(cleaned)
	if err != nil {
		// If we can't resolve symlinks (e.g., path doesn't exist), return cleaned path
		return cleaned, nil
	}

	return resolved, nil
}

// ValidateDirectoryName validates a directory name for common issues
func ValidateDirectoryName(name string) error {
	if name == "" {
		return fmt.Errorf("directory name cannot be empty")
	}

	if name == "." || name == ".." {
		return fmt.Errorf("directory name cannot be '.' or '..'")
	}

	// Check for reserved names on Windows (even though we're primarily targeting Unix)
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upperName == reserved {
			return fmt.Errorf("directory name '%s' is reserved", name)
		}
	}

	// Check for invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("directory name contains invalid character '%s'", char)
		}
	}

	// Check length (most filesystems have a 255 character limit)
	if len(name) > 255 {
		return fmt.Errorf("directory name too long (max 255 characters)")
	}

	return nil
}
