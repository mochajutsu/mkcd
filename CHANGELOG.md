# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Placeholder for future features

### Changed
- Placeholder for future changes

### Deprecated
- Placeholder for future deprecations

### Removed
- Placeholder for future removals

### Fixed
- Placeholder for future fixes

### Security
- Placeholder for future security updates

## [1.0.0] - 2025-01-17

### Added

#### Core Functionality
- **Primary mkcd command** - Create directories with comprehensive workspace initialization
  - Support for nested directory creation with `mkdir -p` behavior
  - Cross-platform compatibility (Linux, macOS, Windows)
  - Shell script generation for seamless `cd` integration
  - Dry-run mode (`--dry-run`) to preview operations without execution
  - Verbose output (`--verbose`) with detailed operation logging
  - Interactive mode (`--interactive`) with user confirmations

#### Git Integration
- **Automatic Git repository initialization** with `--git` flag
  - Configurable default branch (defaults to "main")
  - Support for custom user name and email configuration
  - Initial commit creation with all generated files
- **Remote repository setup** with `--git-remote <url>` option
  - Automatic remote origin configuration
  - Support for HTTPS and SSH Git URLs
  - Integration with repository initialization workflow

#### Editor Integration
- **Intelligent editor detection and launching** with `--editor` and `--open-editor` flags
  - Auto-detection of 15+ popular editors including:
    - Visual Studio Code (`code`, `code-insiders`)
    - Cursor AI Editor (`cursor`)
    - JetBrains IDEs (WebStorm, IntelliJ IDEA, GoLand, PyCharm)
    - Terminal editors (Neovim, Vim, Emacs, Nano)
    - Text editors (Sublime Text, Atom)
  - Respect for `$EDITOR` and `$VISUAL` environment variables
  - Platform-specific editor support (TextEdit on macOS, Notepad on Windows)
  - Background launching for GUI editors, foreground for terminal editors

#### File Generation System
- **README.md generation** with `--readme` flag
  - Project-specific content with customizable templates
  - Author and description integration from configuration
  - Standard sections: Installation, Usage, Features, Contributing, License
- **.gitignore generation** with `--gitignore <type>` option
  - Language-specific templates: `go`, `node`, `python`, `general`
  - Comprehensive ignore patterns for each language ecosystem
  - IDE and OS-specific ignore patterns included
- **LICENSE file generation** with `--license <type>` option
  - Support for MIT and Apache-2.0 licenses
  - Automatic year and author substitution
  - Full license text with proper formatting
- **Custom file creation** with `--touch <files>` option
  - Support for multiple files (comma-separated)
  - Automatic parent directory creation
  - Integration with profile-based file templates

#### Configuration System
- **TOML-based configuration file** at `~/.config/mkcd/mkcd.conf`
  - Hierarchical configuration with core, git, templates, safety, and output sections
  - Automatic configuration validation and error reporting
  - Migration support for future configuration format changes
- **Configuration management commands**:
  - `mkcd config init` - Initialize configuration with sensible defaults
  - `mkcd config show` - Display current configuration in organized sections
  - `mkcd config edit` - Open configuration in user's preferred editor
  - `mkcd config validate` - Validate configuration file for errors
  - `mkcd config reset` - Reset configuration to factory defaults

#### Profile Management System
- **Named configuration profiles** for different project types
  - Built-in profiles: `default`, `dev`, `nodejs`, `python`
  - Profile inheritance and customization capabilities
  - Default profile selection with `--profile <name>` flag
- **Profile management commands**:
  - `mkcd profile list` - Display all available profiles in table format
  - `mkcd profile show <name>` - Show detailed profile configuration
  - `mkcd profile create <name>` - Interactive profile creation wizard
  - `mkcd profile edit <name>` - Edit profile in configuration file
  - `mkcd profile delete <name>` - Remove profile with confirmation
  - `mkcd profile copy <src> <dst>` - Duplicate existing profiles
- **Pre-configured profiles**:
  - `dev` - Git + Editor + README + basic template
  - `nodejs` - Node.js project with package.json, index.js, and Node.js .gitignore
  - `python` - Python project with main.py, requirements.txt, and Python .gitignore

#### Safety and Validation
- **Path validation system** with configurable safety checks
  - Forbidden path protection (prevents operations in system directories)
  - Maximum directory depth limits
  - Dangerous character detection in paths
  - Path sanitization and normalization
- **Backup functionality** with `--backup` flag
  - Automatic backup of existing files before overwrite
  - Timestamped backup files for recovery
  - Configurable backup behavior per operation
- **Interactive confirmations** for potentially destructive operations
  - Overwrite confirmations for existing directories
  - Deletion confirmations for profile and configuration management
  - Force flag (`--force`) to bypass confirmations when needed

#### Rich Terminal Output
- **Enhanced output formatting** using pterm library
  - Colorized output with semantic color coding
  - Icons and symbols for different message types
  - Progress bars for long-running operations
  - Structured table output for lists and configurations
- **Multiple output modes**:
  - Quiet mode (`--quiet`) for script-friendly output
  - Verbose mode (`--verbose`) for detailed operation logging
  - Debug mode (`--debug`) for troubleshooting and development
- **Interactive prompts and selections**
  - Confirmation dialogs with yes/no prompts
  - Multi-select menus for option selection
  - Text input prompts with default values
  - Validation and error handling for user input

#### Command Line Interface
- **Comprehensive flag system** with short and long options
  - Global flags available across all commands
  - Command-specific flags for targeted functionality
  - Mutually exclusive flag validation
  - Help system with detailed descriptions and examples
- **Shell completion support** (framework ready)
  - Cobra-based completion system
  - Support for bash, zsh, and fish shells
  - Dynamic completion for profiles and configuration options

#### Build and Development System
- **Professional Makefile** with 25+ targets
  - Development workflow: `make dev`, `make build`, `make test`
  - Cross-platform builds: `make build-all` for Linux, macOS, Windows
  - Quality assurance: `make lint`, `make test-coverage`
  - Distribution: `make package`, `make release`
  - Installation: `make install`, `make install-global`
- **Go module configuration** with proper dependency management
  - Modern Go 1.24.4+ compatibility
  - Curated dependencies with security considerations
  - Reproducible builds with go.sum verification

### Changed
- N/A (Initial release)

### Deprecated
- N/A (Initial release)

### Removed
- N/A (Initial release)

### Fixed
- N/A (Initial release)

### Security
- **Path traversal protection** - Prevents directory creation outside intended locations
- **Input validation** - Sanitizes all user inputs to prevent injection attacks
- **Safe file operations** - Atomic file operations with proper error handling
- **Configuration validation** - Prevents malformed configuration from causing security issues

---

## Release Notes

### Version 1.0.0 - "Foundation Release"

This initial release establishes mkcd as a comprehensive, enterprise-level directory creation and workspace initialization tool. The focus has been on building a solid foundation with professional-grade features that can scale and evolve.

#### Key Highlights

- **Complete CLI Framework**: Built on Cobra with rich terminal output via pterm
- **Modular Architecture**: Clean separation of concerns with 8+ internal packages
- **Enterprise Features**: Configuration management, profiles, safety checks, and validation
- **Developer Experience**: Intelligent editor integration, git automation, and file generation
- **Production Ready**: Comprehensive error handling, logging, and user feedback

#### Getting Started

```bash
# Install mkcd
go install github.com/mochajutsu/mkcd@latest

# Initialize configuration
mkcd config init

# Create your first project
mkcd my-project --profile dev

# Create a Node.js project with full setup
mkcd my-app --profile nodejs --git-remote https://github.com/user/my-app.git
```

#### Migration Notes
- This is the initial release, no migration required
- Configuration file will be created automatically on first use
- All features are stable and ready for production use

#### Known Limitations
- Template system is planned for future release
- History/undo functionality is planned for future release
- Shell integration wrapper functions are planned for future release
- Batch operations are planned for future release

For detailed usage instructions, see the [README.md](https://github.com/mochajutsu/mkcd/blob/main/README.md) file.

---
