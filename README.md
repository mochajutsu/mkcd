# mkcd - Enterprise Directory Creation Tool

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/mochajutsu/mkcd)](https://github.com/mochajutsu/mkcd/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mochajutsu/mkcd/ci.yml?branch=main)](https://github.com/mochajutsu/mkcd/actions)

A powerful, extensible command-line utility that revolutionizes directory creation and navigation for developers. Built with Go and Cobra, it transforms the simple concept of "make directory and change into it" into a comprehensive workspace initialization tool.

## 🚀 Features

### Core Functionality
- **Cross-platform directory creation** with shell integration
- **Git repository initialization** with remote setup and initial commits
- **Project templates** for different languages and frameworks
- **Editor integration** with auto-detection and launching
- **Configuration profiles** for different project types
- **File generation** (README, .gitignore, LICENSE, custom files)

### Advanced Features
- **Safety checks** and path validation
- **Dry-run mode** for testing operations
- **Interactive confirmations** for destructive operations
- **Verbose output** with progress indicators
- **Backup functionality** for existing directories
- **Shell script generation** for seamless cd integration

### Enterprise-Level Quality
- **Comprehensive error handling** and recovery
- **Extensive logging** and debugging capabilities
- **Configuration validation** and migration
- **Professional documentation** and examples
- **Senior developer-level code quality** with detailed comments

## 📦 Installation

### Binary Downloads (Recommended)
Download the latest release for your platform:

**Linux (x86_64)**
```bash
curl -L https://github.com/mochajutsu/mkcd/releases/latest/download/mkcd-v1.0.0-linux-amd64.tar.gz | tar xz
sudo mv mkcd-v1.0.0-linux-amd64/mkcd /usr/local/bin/
```

**Linux (ARM64)**
```bash
curl -L https://github.com/mochajutsu/mkcd/releases/latest/download/mkcd-v1.0.0-linux-arm64.tar.gz | tar xz
sudo mv mkcd-v1.0.0-linux-arm64/mkcd /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/mochajutsu/mkcd/releases/latest/download/mkcd-v1.0.0-darwin-amd64.tar.gz | tar xz
sudo mv mkcd-v1.0.0-darwin-amd64/mkcd /usr/local/bin/
```

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/mochajutsu/mkcd/releases/latest/download/mkcd-v1.0.0-darwin-arm64.tar.gz | tar xz
sudo mv mkcd-v1.0.0-darwin-arm64/mkcd /usr/local/bin/
```

**Windows**
1. Download `mkcd-v1.0.0-windows-amd64.tar.gz` from [releases](https://github.com/mochajutsu/mkcd/releases)
2. Extract and add `mkcd.exe` to your PATH

### Verify Installation

After installation, verify mkcd is working:

```bash
# Check version
mkcd --version

# View help
mkcd --help

# Test with dry-run
mkcd test-project --dry-run --verbose
```

### Using Go Install
```bash
go install github.com/mochajutsu/mkcd@latest
```

### Package Managers

**Homebrew (macOS/Linux)**
```bash
brew tap mochajutsu/tap
brew install mkcd
```

**Arch Linux (AUR)**
```bash
yay -S mkcd-bin
# or
paru -S mkcd-bin
```

**Docker**
```bash
docker run --rm -v $(pwd):/workspace ghcr.io/mochajutsu/mkcd:latest --help
```

### Manual Installation
```bash
# Clone the repository
git clone https://github.com/mochajutsu/mkcd.git
cd mkcd

# Build and install
make install

# Or build manually
go build -o mkcd .
sudo mv mkcd /usr/local/bin/
```

## 🎯 Quick Start

### Basic Usage
```bash
# Create a simple directory
mkcd myproject

# Create with Git repository
mkcd myproject --git

# Create with Git and open in editor
mkcd myproject --git --open-editor

# Create with profile
mkcd myproject --profile nodejs
```

### Using Profiles
```bash
# List available profiles
mkcd profile list

# Create a Node.js project
mkcd my-app --profile nodejs
# This automatically:
# - Creates the directory
# - Initializes Git repository
# - Generates package.json and index.js
# - Creates Node.js .gitignore
# - Opens in your preferred editor

# Create a Python project
mkcd my-script --profile python
# This automatically:
# - Creates the directory
# - Initializes Git repository
# - Generates main.py and requirements.txt
# - Creates Python .gitignore
# - Opens in your preferred editor
```

### Advanced Usage
```bash
# Create with custom files and settings
mkcd myproject \
  --git \
  --git-remote https://github.com/user/myproject.git \
  --readme \
  --gitignore go \
  --license mit \
  --touch "main.go,go.mod" \
  --editor code

# Use dry-run to see what would happen
mkcd myproject --profile dev --dry-run --verbose

# Interactive mode with confirmations
mkcd myproject --interactive --git --readme
```

## ⚙️ Configuration

### Initialize Configuration
```bash
# Create default configuration file
mkcd config init

# View current configuration
mkcd config show

# Edit configuration in your editor
mkcd config edit
```

### Configuration File Location

- **Linux/macOS**: `~/.config/mkcd/mkcd.conf`
- **Windows**: `%APPDATA%\mkcd\mkcd.conf`

### Sample Configuration

```toml
[core]
default_profile = "dev"
editor = "code"
shell_integration = true
history_limit = 100
backup_enabled = false

[git]
auto_init = false
default_branch = "main"
user_name = "Your Name"
user_email = "your.email@example.com"

[profiles.dev]
git = true
editor = true
readme = true
gitignore = "general"
template = "basic-dev"

[profiles.nodejs]
git = true
editor = true
template = "nodejs"
gitignore = "node"
touch = ["package.json", "index.js"]
```

## 🔧 Commands

### Main Command

```bash
mkcd <directory> [flags]
```

**Flags:**

- `--git` - Initialize Git repository
- `--git-remote <url>` - Add remote origin
- `--template <name>` - Apply project template
- `--editor <editor>` - Open in specific editor
- `--open-editor` - Open in auto-detected editor
- `--readme` - Generate README.md
- `--gitignore <type>` - Generate .gitignore (go, node, python, general)
- `--license <type>` - Generate LICENSE (mit, apache-2.0)
- `--touch <files>` - Create specified files
- `--profile <name>` - Use configuration profile
- `--dry-run` - Show what would be done
- `--verbose` - Detailed output
- `--interactive` - Interactive confirmations

### Profile Management

```bash
mkcd profile list                    # List all profiles
mkcd profile show <name>             # Show profile details
mkcd profile create <name>           # Create new profile
mkcd profile edit <name>             # Edit profile
mkcd profile delete <name>           # Delete profile
mkcd profile copy <src> <dst>        # Copy profile
```

### Configuration Management

```bash
mkcd config init                     # Initialize config file
mkcd config show                     # Show current config
mkcd config edit                     # Edit config in editor
mkcd config validate                 # Validate configuration
mkcd config reset                    # Reset to defaults
```

## 🎨 Examples

### Quick Start Examples

```bash
# Create a simple directory
mkcd my-project

# Create with Git repository
mkcd my-project --git

# Create with Git and open in editor
mkcd my-project --git --open-editor

# Use a profile for instant setup
mkcd my-app --profile nodejs
```

### Real-World Scenarios

**Web Development Project**
```bash
mkcd my-website \
  --profile nodejs \
  --git \
  --git-remote https://github.com/user/my-website.git \
  --license mit
```

**Go CLI Tool**
```bash
mkcd my-tool \
  --git \
  --readme \
  --gitignore go \
  --license mit \
  --touch "main.go,go.mod,Makefile" \
  --editor code
```

**Python Data Science Project**
```bash
mkcd data-analysis \
  --profile python \
  --git \
  --touch "notebook.ipynb,requirements.txt" \
  --git-remote https://github.com/user/data-analysis.git
```

**Quick Prototype**
```bash
mkcd prototype \
  --git \
  --readme \
  --gitignore general \
  --open-editor
```

**Enterprise Project Setup**
```bash
mkcd enterprise-app \
  --profile dev \
  --git \
  --git-remote https://github.com/company/enterprise-app.git \
  --license apache-2.0 \
  --touch "docker-compose.yml,.env.example" \
  --interactive
```

## 🏗️ Architecture

mkcd is built with a modular, enterprise-level architecture:

```
mkcd/
├── cmd/                    # CLI commands and interfaces
│   ├── root.go            # Root command and global flags
│   ├── mkcd.go            # Main mkcd command
│   ├── profile.go         # Profile management
│   └── config.go          # Configuration management
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── editor/            # Editor detection and launching
│   ├── files/             # File generation (README, .gitignore, etc.)
│   ├── git/               # Git operations
│   └── utils/             # Shared utilities
└── templates/             # Built-in project templates
```

### Key Components

- **Configuration System**: TOML-based configuration with profiles and validation
- **Editor Integration**: Auto-detection and launching of 15+ editors
- **Git Integration**: Repository initialization, remote setup, and commit creation
- **File Generation**: Smart generation of project files with templates
- **Safety System**: Path validation, forbidden directory protection
- **Output Management**: Rich terminal output with colors, icons, and progress bars

## 🧪 Development

### Prerequisites

- Go 1.24.4 or later
- Make (optional, for build automation)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/mochajutsu/mkcd.git
cd mkcd

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Lint code
make lint
```

### Development Workflow

```bash
# Build for development
make dev

# Watch for changes (requires fswatch)
make watch

# Build for all platforms
make build-all

# Create release packages
make package
```

### Testing

```bash
# Run all tests
go test ./...

# Test with coverage
go test -cover ./...

# Test specific package
go test ./internal/config

# Benchmark tests
go test -bench=. ./...
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes with tests
4. Run the test suite: `make test`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Standards
- Follow Go best practices and idioms
- Write comprehensive tests for new features
- Add detailed comments for complex logic
- Use meaningful variable and function names
- Ensure all tests pass and coverage remains high

## 📚 Documentation

- [Release Guide](docs/RELEASE.md) - How to create releases and distribute packages
- [Changelog](CHANGELOG.md) - Version history and release notes
- Configuration Reference (Coming Soon)
- Template System (Coming Soon)
- Shell Integration (Coming Soon)
- API Documentation (Coming Soon)

## 🔒 Security

If you discover a security vulnerability, please send an e-mail to [<mochajutsu@gmail.com>. All security vulnerabilities will be promptly addressed.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [pterm](https://github.com/pterm/pterm) - Terminal output library
- [go-git](https://github.com/go-git/go-git) - Git operations
- [TOML](https://github.com/BurntSushi/toml) - Configuration parsing

## 📊 Project Status

**Current Version**: v1.0.0 🎉

mkcd is actively maintained and under continuous development. We follow semantic versioning and maintain backward compatibility.

### ✅ What's Ready (v1.0.0)

- ✅ **Core mkcd command** with comprehensive options
- ✅ **Profile management** system with built-in profiles
- ✅ **Configuration management** with TOML format
- ✅ **Git integration** with repository initialization and remotes
- ✅ **Editor integration** with 15+ supported editors
- ✅ **File generation** (README, .gitignore, LICENSE)
- ✅ **Safety features** and path validation
- ✅ **Cross-platform builds** (Linux, macOS, Windows)
- ✅ **Professional build system** with Makefile
- ✅ **Automated releases** with GitHub Actions

### 🚀 Roadmap (Future Versions)

- [ ] **Template marketplace** and sharing
- [ ] **Plugin system** for extensibility
- [ ] **History and undo** functionality
- [ ] **Batch operations** with pattern support
- [ ] **Shell integration** wrapper functions
- [ ] **Cloud integration** for template sync
- [ ] **IDE extensions** (VSCode, JetBrains)
- [ ] **Advanced workflow** automation
- [ ] **Team collaboration** features

---

**Made with ❤️ by [mochajutsu](https://github.com/mochajutsu)**