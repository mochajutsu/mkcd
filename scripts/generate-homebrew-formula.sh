#!/bin/bash
# Generate Homebrew formula for mkcd
# Usage: ./generate-homebrew-formula.sh <version>

set -e

VERSION=${1:-"1.0.0"}
REPO_URL="https://github.com/mochajutsu/mkcd"
ARCHIVE_URL="${REPO_URL}/archive/v${VERSION}.tar.gz"

# Calculate SHA256 for the source archive
# Note: In practice, you'd download the actual archive and calculate its hash
# For now, we'll use a placeholder that needs to be updated manually
SHA256="# TODO: Replace with actual SHA256 of v${VERSION}.tar.gz"

cat << EOF
class Mkcd < Formula
  desc "Enterprise directory creation and workspace initialization tool"
  homepage "${REPO_URL}"
  url "${ARCHIVE_URL}"
  sha256 "${SHA256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.Version=#{version}")
    
    # Install shell completions
    generate_completions_from_executable(bin/"mkcd", "completion")
    
    # Install man page if available
    # man1.install "docs/mkcd.1" if File.exist?("docs/mkcd.1")
  end

  test do
    # Test basic functionality
    system "#{bin}/mkcd", "--version"
    system "#{bin}/mkcd", "--help"
    
    # Test config initialization
    system "#{bin}/mkcd", "config", "init", "--dry-run"
    
    # Test profile listing
    system "#{bin}/mkcd", "profile", "list"
  end
end
EOF
