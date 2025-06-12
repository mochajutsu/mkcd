# mkcd

> 🍵 Create a new directory and immediately jump into it — with extras.

`mkcd` is a minimalist project bootstrapper and directory-jumping CLI built for developers who are tired of typing `mkdir mydir && cd mydir`.

Born from the time-honored shell alias that nearly every developer adds to their config (`alias mkcd='mkdir -p $1 && cd $1'`), this project takes that humble idea and expands it with smart defaults, developer-friendly scaffolding, and zero-runtime dependencies.

## ✨ Features

- **📁 Smart directory creation** — Create and navigate to directories in one command
- **🧠 Intelligent project setup** — Optional git initialization, README generation, and .gitignore creation
- **🗂 Stack bootstrapping** — Built-in templates for Node.js, Python, Go, Rust, and more
- **⚙️ Custom templates** — Define your own project scaffolds
- **🧪 Test scaffolding** — Optional test directory and config setup
- **🎯 Cross-platform** — Native support for Linux and macOS
- **🧊 Zero dependencies** — Distributed as a single native binary, no Node.js required
- **⚡ Lightning fast** — Built in TypeScript, compiled to native code

## 🛠 Installation

### Via Homebrew (macOS & Linux)
```bash
brew install mochajutsu/mkcd/mkcd
```

### Via AUR (Arch Linux)
```bash
yay -S mkcd
```

### Manual Installation
Download the latest binary from [Releases](https://github.com/mochajutsu/mkcd/releases) and add it to your PATH.

## 🚀 Usage

### Basic Usage
```bash
# Create directory and cd into it
mkcd my-new-project

# With git initialization
mkcd my-app --git

# Full project scaffold
mkcd my-startup --git --template=node --readme --license=MIT --open
```

### Common Templates
```bash
# Node.js project
mkcd my-node-app --template=node --git --readme

# Python project with virtual environment
mkcd my-python-app --template=python --git --venv

# Go module
mkcd my-go-app --template=go --git --mod

# Rust project
mkcd my-rust-app --template=rust --git
```

### Advanced Options
```bash
# Custom template from URL
mkcd my-project --template=https://github.com/user/template.git

# Multiple features at once
mkcd full-stack-app \
  --git \
  --template=node \
  --readme \
  --license=MIT \
  --tests \
  --open=code
```

## 📋 Command Reference

### Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--git` | Initialize git repository | `mkcd app --git` |
| `--template=<name>` | Use built-in or custom template | `--template=node` |
| `--readme` | Generate README.md | `mkcd app --readme` |
| `--license=<type>` | Add license file | `--license=MIT` |
| `--tests` | Create test directory structure | `mkcd app --tests` |
| `--open[=editor]` | Open in editor after creation | `--open=code` |
| `--config=<path>` | Use custom config file | `--config=~/.mkcd-work` |

### Built-in Templates

| Template | Description | Includes |
|----------|-------------|----------|
| `node` | Node.js project | package.json, .gitignore, basic structure |
| `python` | Python project | requirements.txt, .gitignore, virtual env setup |
| `go` | Go module | go.mod, main.go, basic structure |
| `rust` | Rust project | Cargo.toml, src/main.rs, .gitignore |
| `web` | Static web project | index.html, style.css, script.js |
| `docs` | Documentation site | index.md, basic structure |

## ⚙️ Configuration

Create `~/.mkcdrc` to customize default behavior:

```toml
[defaults]
git = true
readme = true
editor = "code"
license = "MIT"

[templates]
# Custom template shortcuts
react = "https://github.com/user/react-template.git"
api = "node"

[aliases]
# Command aliases
new = "mkcd"
start = "mkcd"
```

## 🔧 Creating Custom Templates

Templates are directories with optional `.mkcd` metadata:

```
my-template/
├── .mkcd/
│   ├── config.toml     # Template configuration
│   └── hooks.sh        # Pre/post creation hooks
├── src/
├── package.json
└── README.md
```

Example `.mkcd/config.toml`:
```toml
name = "My Custom Template"
description = "A template for my preferred stack"
variables = ["project_name", "author"]
```

## 📦 Roadmap

- [ ] **Remote template registry** — Share and discover community templates
- [ ] **Plugin system** — Extend functionality with custom plugins
- [ ] **Interactive mode** — Guided project setup with prompts
- [ ] **Template versioning** — Lock templates to specific versions
- [ ] **Workspace support** — Multi-project directory management
- [ ] **Cloud integration** — Sync templates and configs across machines

## 🤝 Contributing

We welcome contributions of all kinds! Whether you're fixing bugs, adding features, improving documentation, or creating templates, your help makes `mkcd` better for everyone.

Check out our [Contributing Guide](CONTRIBUTING.md) to get started.

### Quick Start for Contributors
```bash
git clone https://github.com/mochajutsu/mkcd.git
cd mkcd
npm install
npm run dev
```

## 📜 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Inspired by the countless developers who've added `alias mkcd='mkdir -p $1 && cd $1'` to their shell configs. This tool is for all of us who knew there had to be a better way.

---

**Built with ☕ by [@mochajutsu](https://github.com/mochajutsu)**

*Making project setup fast, consistent, and less annoying, one directory at a time.*
