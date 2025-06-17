/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

// Package utils provides common utility functions for filesystem operations,
// path manipulation, and other shared functionality across the mkcd application.
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// FileSystemOperations provides filesystem utility functions
type FileSystemOperations struct {
	DryRun bool
	Backup bool
}

// NewFileSystemOperations creates a new FileSystemOperations instance
func NewFileSystemOperations(dryRun, backup bool) *FileSystemOperations {
	return &FileSystemOperations{
		DryRun: dryRun,
		Backup: backup,
	}
}

// CreateDirectory creates a directory with the specified permissions
// If the directory already exists, it returns nil (no error)
func (fs *FileSystemOperations) CreateDirectory(path string, mode os.FileMode) error {
	if fs.DryRun {
		pterm.Info.Printf("[DRY RUN] Would create directory: %s (mode: %o)", path, mode)
		return nil
	}

	// Check if directory already exists
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			pterm.Debug.Printf("Directory already exists: %s", path)
			return nil
		}
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	// Create directory with parents
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	pterm.Success.Printf("Created directory: %s", path)
	return nil
}

// CreateFile creates a file with the specified content
func (fs *FileSystemOperations) CreateFile(path, content string, mode os.FileMode) error {
	if fs.DryRun {
		pterm.Info.Printf("[DRY RUN] Would create file: %s (size: %d bytes)", path, len(content))
		return nil
	}

	// Check if file already exists and backup if needed
	if fs.Backup {
		if _, err := os.Stat(path); err == nil {
			if err := fs.BackupFile(path); err != nil {
				return fmt.Errorf("failed to backup existing file: %w", err)
			}
		}
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory %s: %w", dir, err)
	}

	// Create and write file
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content to file %s: %w", path, err)
	}

	pterm.Success.Printf("Created file: %s", path)
	return nil
}

// BackupFile creates a backup of the specified file
func (fs *FileSystemOperations) BackupFile(path string) error {
	if fs.DryRun {
		pterm.Info.Printf("[DRY RUN] Would backup file: %s", path)
		return nil
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", path, timestamp)

	// Copy file to backup location
	if err := CopyFile(path, backupPath); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupPath, err)
	}

	pterm.Info.Printf("Created backup: %s", backupPath)
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Copy file contents
	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := sourceFile.Read(buffer)
		if n > 0 {
			if _, writeErr := destFile.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to destination file: %w", writeErr)
			}
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read from source file: %w", err)
		}
	}

	// Copy file permissions
	if info, err := sourceFile.Stat(); err == nil {
		if err := destFile.Chmod(info.Mode()); err != nil {
			pterm.Warning.Printf("Failed to copy file permissions: %v", err)
		}
	}

	return nil
}

// PathExists checks if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory checks if a path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile checks if a path is a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// GetAbsolutePath returns the absolute path, expanding ~ to home directory
func GetAbsolutePath(path string) (string, error) {
	// Handle ~ expansion
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	return absPath, nil
}

// SanitizePath cleans and validates a path
func SanitizePath(path string) (string, error) {
	// Clean the path
	cleaned := filepath.Clean(path)

	// Check for dangerous patterns
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path contains dangerous '..' components: %s", path)
	}

	// Ensure path is not empty
	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("invalid empty path")
	}

	return cleaned, nil
}

// CreateSymlink creates a symbolic link
func (fs *FileSystemOperations) CreateSymlink(target, linkPath string) error {
	if fs.DryRun {
		pterm.Info.Printf("[DRY RUN] Would create symlink: %s -> %s", linkPath, target)
		return nil
	}

	// Check if target exists
	if !PathExists(target) {
		return fmt.Errorf("symlink target does not exist: %s", target)
	}

	// Remove existing link if it exists
	if PathExists(linkPath) {
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to remove existing symlink %s: %w", linkPath, err)
		}
	}

	// Create symlink
	if err := os.Symlink(target, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink %s -> %s: %w", linkPath, target, err)
	}

	pterm.Success.Printf("Created symlink: %s -> %s", linkPath, target)
	return nil
}

// GetDirectorySize calculates the total size of a directory
func GetDirectorySize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// ListDirectory returns a list of files and directories in the specified path
func ListDirectory(path string) ([]os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory %s: %w", path, err)
	}
	defer file.Close()

	entries, err := file.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	return entries, nil
}

// EnsureDirectoryExists creates a directory if it doesn't exist
func EnsureDirectoryExists(path string) error {
	if !PathExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}
