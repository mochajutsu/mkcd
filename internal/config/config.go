/*
Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

Licensed under the MIT License. See LICENSE file for details.
*/

// Package config provides configuration management for mkcd.
// It handles loading, validation, and management of TOML configuration files,
// profiles, and application settings.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/pterm/pterm"
)

// Config represents the main configuration structure for mkcd
type Config struct {
	Core      CoreConfig              `toml:"core"`
	Git       GitConfig               `toml:"git"`
	Templates TemplatesConfig         `toml:"templates"`
	Safety    SafetyConfig            `toml:"safety"`
	Output    OutputConfig            `toml:"output"`
	Profiles  map[string]ProfileConfig `toml:"profiles"`
}

// CoreConfig contains core application settings
type CoreConfig struct {
	DefaultProfile    string `toml:"default_profile"`
	Editor            string `toml:"editor"`
	ShellIntegration  bool   `toml:"shell_integration"`
	HistoryLimit      int    `toml:"history_limit"`
	BackupEnabled     bool   `toml:"backup_enabled"`
	TempDir           string `toml:"temp_dir"`
}

// GitConfig contains git-related configuration
type GitConfig struct {
	AutoInit           bool   `toml:"auto_init"`
	DefaultBranch      string `toml:"default_branch"`
	UserName           string `toml:"user_name"`
	UserEmail          string `toml:"user_email"`
	DefaultRemoteName  string `toml:"default_remote_name"`
}

// TemplatesConfig contains template system configuration
type TemplatesConfig struct {
	Directory  string `toml:"directory"`
	AutoUpdate bool   `toml:"auto_update"`
}

// SafetyConfig contains safety and validation settings
type SafetyConfig struct {
	ConfirmOverwrites bool     `toml:"confirm_overwrites"`
	ConfirmDeletes    bool     `toml:"confirm_deletes"`
	MaxDepth          int      `toml:"max_depth"`
	ForbiddenPaths    []string `toml:"forbidden_paths"`
}

// OutputConfig contains output formatting settings
type OutputConfig struct {
	Colors       bool `toml:"colors"`
	Icons        bool `toml:"icons"`
	ProgressBars bool `toml:"progress_bars"`
}

// ProfileConfig represents a named configuration profile
type ProfileConfig struct {
	Git       bool     `toml:"git"`
	Editor    bool     `toml:"editor"`
	Readme    bool     `toml:"readme"`
	Gitignore string   `toml:"gitignore"`
	Template  string   `toml:"template"`
	Touch     []string `toml:"touch"`
	License   string   `toml:"license"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := homedir.Dir()
	
	return &Config{
		Core: CoreConfig{
			DefaultProfile:   "default",
			Editor:           "",
			ShellIntegration: true,
			HistoryLimit:     100,
			BackupEnabled:    false,
			TempDir:          "/tmp/mkcd",
		},
		Git: GitConfig{
			AutoInit:          false,
			DefaultBranch:     "main",
			UserName:          "",
			UserEmail:         "",
			DefaultRemoteName: "origin",
		},
		Templates: TemplatesConfig{
			Directory:  filepath.Join(homeDir, ".config", "mkcd", "templates"),
			AutoUpdate: false,
		},
		Safety: SafetyConfig{
			ConfirmOverwrites: true,
			ConfirmDeletes:    true,
			MaxDepth:          10,
			ForbiddenPaths:    []string{"/", "/usr", "/etc", "/var", "/bin", "/sbin"},
		},
		Output: OutputConfig{
			Colors:       true,
			Icons:        true,
			ProgressBars: true,
		},
		Profiles: map[string]ProfileConfig{
			"default": {
				Git:    false,
				Editor: false,
				Readme: false,
			},
			"dev": {
				Git:       true,
				Editor:    true,
				Readme:    true,
				Gitignore: "general",
				Template:  "basic-dev",
			},
			"nodejs": {
				Git:       true,
				Editor:    true,
				Template:  "nodejs",
				Gitignore: "node",
				Touch:     []string{"package.json", "index.js"},
			},
			"python": {
				Git:       true,
				Editor:    true,
				Template:  "python",
				Gitignore: "python",
				Touch:     []string{"main.py", "requirements.txt"},
			},
		},
	}
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".config", "mkcd")
	configFile := filepath.Join(configDir, "mkcd.conf")
	
	return configFile, nil
}

// Load loads configuration from the specified file path
// If the file doesn't exist, it returns the default configuration
func Load(configPath string) (*Config, error) {
	// If no config path specified, use default
	if configPath == "" {
		var err error
		configPath, err = GetConfigPath()
		if err != nil {
			return nil, fmt.Errorf("failed to determine config path: %w", err)
		}
	}
	
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		pterm.Debug.Printf("Config file not found at %s, using defaults", configPath)
		return DefaultConfig(), nil
	}
	
	// Load and parse config file
	config := DefaultConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}
	
	// Validate the loaded configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	pterm.Debug.Printf("Loaded configuration from %s", configPath)
	return config, nil
}

// Save saves the configuration to the specified file path
func (c *Config) Save(configPath string) error {
	// If no config path specified, use default
	if configPath == "" {
		var err error
		configPath, err = GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to determine config path: %w", err)
		}
	}
	
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}
	
	// Create config file
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file %s: %w", configPath, err)
	}
	defer file.Close()
	
	// Encode configuration to TOML
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config to TOML: %w", err)
	}
	
	pterm.Success.Printf("Configuration saved to %s", configPath)
	return nil
}

// Validate validates the configuration for consistency and correctness
func (c *Config) Validate() error {
	// Validate core settings
	if c.Core.HistoryLimit < 0 {
		return fmt.Errorf("history_limit must be non-negative")
	}
	
	if c.Safety.MaxDepth < 1 {
		return fmt.Errorf("max_depth must be at least 1")
	}
	
	// Validate default profile exists
	if c.Core.DefaultProfile != "" {
		if _, exists := c.Profiles[c.Core.DefaultProfile]; !exists {
			return fmt.Errorf("default profile '%s' does not exist", c.Core.DefaultProfile)
		}
	}
	
	// Validate forbidden paths are absolute
	for _, path := range c.Safety.ForbiddenPaths {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("forbidden path '%s' must be absolute", path)
		}
	}
	
	return nil
}

// GetProfile returns the specified profile or the default profile if name is empty
func (c *Config) GetProfile(name string) (ProfileConfig, error) {
	if name == "" {
		name = c.Core.DefaultProfile
	}
	
	profile, exists := c.Profiles[name]
	if !exists {
		return ProfileConfig{}, fmt.Errorf("profile '%s' not found", name)
	}
	
	return profile, nil
}

// SetProfile sets or updates a profile in the configuration
func (c *Config) SetProfile(name string, profile ProfileConfig) {
	if c.Profiles == nil {
		c.Profiles = make(map[string]ProfileConfig)
	}
	c.Profiles[name] = profile
}

// DeleteProfile removes a profile from the configuration
func (c *Config) DeleteProfile(name string) error {
	if name == c.Core.DefaultProfile {
		return fmt.Errorf("cannot delete default profile '%s'", name)
	}
	
	if _, exists := c.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' does not exist", name)
	}
	
	delete(c.Profiles, name)
	return nil
}
