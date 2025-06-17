/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

// Package git provides Git repository management functionality for mkcd.
// It handles repository initialization, remote setup, and basic Git operations
// using the go-git library for pure Go implementation.
package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pterm/pterm"
)

// GitManager handles Git operations for mkcd
type GitManager struct {
	DryRun    bool
	Verbose   bool
	UserName  string
	UserEmail string
}

// NewGitManager creates a new GitManager instance
func NewGitManager(dryRun, verbose bool, userName, userEmail string) *GitManager {
	return &GitManager{
		DryRun:    dryRun,
		Verbose:   verbose,
		UserName:  userName,
		UserEmail: userEmail,
	}
}

// InitRepository initializes a new Git repository in the specified directory
func (gm *GitManager) InitRepository(path string, defaultBranch string) error {
	if gm.DryRun {
		pterm.Info.Printf("[DRY RUN] Would initialize Git repository in: %s", path)
		return nil
	}

	// Check if repository already exists
	if gm.isGitRepository(path) {
		pterm.Debug.Printf("Git repository already exists in: %s", path)
		return nil
	}

	// Initialize repository
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	// Set default branch if specified
	if defaultBranch != "" && defaultBranch != "master" {
		if err := gm.setDefaultBranch(repo, defaultBranch); err != nil {
			pterm.Warning.Printf("Failed to set default branch to %s: %v", defaultBranch, err)
		}
	}

	pterm.Success.Printf("Initialized Git repository in: %s", path)
	return nil
}

// isGitRepository checks if a directory is already a Git repository
func (gm *GitManager) isGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if info, err := os.Stat(gitDir); err == nil {
		return info.IsDir()
	}
	return false
}

// setDefaultBranch sets the default branch for the repository
func (gm *GitManager) setDefaultBranch(repo *git.Repository, branchName string) error {
	// Get repository configuration
	cfg, err := repo.Config()
	if err != nil {
		return fmt.Errorf("failed to get repository config: %w", err)
	}

	// Update init.defaultBranch setting
	cfg.Init.DefaultBranch = branchName

	// Save configuration
	if err := repo.Storer.SetConfig(cfg); err != nil {
		return fmt.Errorf("failed to save repository config: %w", err)
	}

	return nil
}

// AddRemote adds a remote repository to the Git repository
func (gm *GitManager) AddRemote(repoPath, remoteName, remoteURL string) error {
	if gm.DryRun {
		pterm.Info.Printf("[DRY RUN] Would add remote %s: %s", remoteName, remoteURL)
		return nil
	}

	// Open repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open Git repository: %w", err)
	}

	// Check if remote already exists
	if _, err := repo.Remote(remoteName); err == nil {
		pterm.Debug.Printf("Remote %s already exists", remoteName)
		return nil
	}

	// Create remote configuration
	remoteConfig := &config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteURL},
	}

	// Add remote
	_, err = repo.CreateRemote(remoteConfig)
	if err != nil {
		return fmt.Errorf("failed to add remote %s: %w", remoteName, err)
	}

	pterm.Success.Printf("Added remote %s: %s", remoteName, remoteURL)
	return nil
}

// CreateInitialCommit creates an initial commit with any existing files
func (gm *GitManager) CreateInitialCommit(repoPath, message string) error {
	if gm.DryRun {
		pterm.Info.Printf("[DRY RUN] Would create initial commit: %s", message)
		return nil
	}

	// Open repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open Git repository: %w", err)
	}

	// Get working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get working tree: %w", err)
	}

	// Add all files
	if err := worktree.AddGlob("."); err != nil {
		return fmt.Errorf("failed to add files to staging: %w", err)
	}

	// Check if there are any changes to commit
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get repository status: %w", err)
	}

	if status.IsClean() {
		pterm.Debug.Println("No changes to commit")
		return nil
	}

	// Create commit
	author := gm.getCommitAuthor()
	commitHash, err := worktree.Commit(message, &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	pterm.Success.Printf("Created initial commit: %s", commitHash.String()[:8])
	return nil
}

// getCommitAuthor returns the commit author information
func (gm *GitManager) getCommitAuthor() *object.Signature {
	name := gm.UserName
	email := gm.UserEmail

	// Try to get from git config if not provided
	if name == "" {
		name = gm.getGitConfig("user.name")
	}
	if email == "" {
		email = gm.getGitConfig("user.email")
	}

	// Use defaults if still empty
	if name == "" {
		name = "mkcd user"
	}
	if email == "" {
		email = "user@example.com"
	}

	return &object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}
}

// getGitConfig retrieves a git configuration value
func (gm *GitManager) getGitConfig(key string) string {
	// This is a simplified implementation
	// In a real scenario, you might want to use git config commands
	// or parse the global git config file
	return ""
}

// GetRepositoryInfo returns information about the Git repository
func (gm *GitManager) GetRepositoryInfo(repoPath string) (*RepositoryInfo, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Git repository: %w", err)
	}

	info := &RepositoryInfo{
		Path: repoPath,
	}

	// Get current branch
	head, err := repo.Head()
	if err == nil {
		info.CurrentBranch = head.Name().Short()
	}

	// Get remotes
	remotes, err := repo.Remotes()
	if err == nil {
		info.Remotes = make(map[string]string)
		for _, remote := range remotes {
			cfg := remote.Config()
			if len(cfg.URLs) > 0 {
				info.Remotes[cfg.Name] = cfg.URLs[0]
			}
		}
	}

	// Get last commit
	if head != nil {
		commit, err := repo.CommitObject(head.Hash())
		if err == nil {
			info.LastCommit = &CommitInfo{
				Hash:    commit.Hash.String(),
				Message: commit.Message,
				Author:  commit.Author.Name,
				Date:    commit.Author.When,
			}
		}
	}

	return info, nil
}

// RepositoryInfo contains information about a Git repository
type RepositoryInfo struct {
	Path          string
	CurrentBranch string
	Remotes       map[string]string
	LastCommit    *CommitInfo
}

// CommitInfo contains information about a Git commit
type CommitInfo struct {
	Hash    string
	Message string
	Author  string
	Date    time.Time
}

// ValidateRemoteURL validates a Git remote URL
func ValidateRemoteURL(url string) error {
	if url == "" {
		return fmt.Errorf("remote URL cannot be empty")
	}

	// Basic validation for common Git URL formats
	validPrefixes := []string{
		"https://",
		"http://",
		"git://",
		"ssh://",
		"git@",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return nil
		}
	}

	return fmt.Errorf("invalid Git remote URL format: %s", url)
}

// CloneRepository clones a repository to the specified path
func (gm *GitManager) CloneRepository(url, path string, shallow bool) error {
	if gm.DryRun {
		pterm.Info.Printf("[DRY RUN] Would clone repository %s to %s", url, path)
		return nil
	}

	// Validate URL
	if err := ValidateRemoteURL(url); err != nil {
		return err
	}

	// Clone options
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	if shallow {
		cloneOptions.Depth = 1
	}

	// Clone repository
	_, err := git.PlainClone(path, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	pterm.Success.Printf("Cloned repository %s to %s", url, path)
	return nil
}

// GetBranches returns a list of branches in the repository
func (gm *GitManager) GetBranches(repoPath string) ([]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Git repository: %w", err)
	}

	branches := []string{}
	
	// Get branch references
	refs, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return branches, nil
}
