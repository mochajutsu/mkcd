#!/bin/bash
# Generate Arch Linux PKGBUILD for mkcd
# Usage: ./generate-pkgbuild.sh <version>

set -e

VERSION=${1:-"1.0.0"}
REPO_URL="https://github.com/mochajutsu/mkcd"
ARCHIVE_URL="${REPO_URL}/archive/v${VERSION}.tar.gz"

# Calculate SHA256 for the source archive
# Note: In practice, you'd download the actual archive and calculate its hash
# For now, we'll use a placeholder that needs to be updated manually
SHA256="# TODO: Replace with actual SHA256 of v${VERSION}.tar.gz"

cat << EOF
# Maintainer: mochajutsu <https://github.com/mochajutsu>
pkgname=mkcd
pkgver=${VERSION}
pkgrel=1
pkgdesc="Enterprise directory creation and workspace initialization tool"
arch=('x86_64' 'aarch64')
url="${REPO_URL}"
license=('MIT')
depends=('glibc')
makedepends=('go')
source=("\${pkgname}-\${pkgver}.tar.gz::${ARCHIVE_URL}")
sha256sums=('${SHA256}')

build() {
    cd "\${pkgname}-\${pkgver}"
    export CGO_CPPFLAGS="\${CPPFLAGS}"
    export CGO_CFLAGS="\${CFLAGS}"
    export CGO_CXXFLAGS="\${CXXFLAGS}"
    export CGO_LDFLAGS="\${LDFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
    
    go build -ldflags="-s -w -X main.Version=\${pkgver}" -o \${pkgname} .
}

check() {
    cd "\${pkgname}-\${pkgver}"
    go test ./...
}

package() {
    cd "\${pkgname}-\${pkgver}"
    
    # Install binary
    install -Dm755 \${pkgname} "\${pkgdir}/usr/bin/\${pkgname}"
    
    # Install license
    install -Dm644 LICENSE "\${pkgdir}/usr/share/licenses/\${pkgname}/LICENSE"
    
    # Install documentation
    install -Dm644 README.md "\${pkgdir}/usr/share/doc/\${pkgname}/README.md"
    install -Dm644 CHANGELOG.md "\${pkgdir}/usr/share/doc/\${pkgname}/CHANGELOG.md"
    
    # Install shell completions
    install -dm755 "\${pkgdir}/usr/share/bash-completion/completions"
    install -dm755 "\${pkgdir}/usr/share/zsh/site-functions"
    install -dm755 "\${pkgdir}/usr/share/fish/vendor_completions.d"
    
    "\${pkgdir}/usr/bin/\${pkgname}" completion bash > "\${pkgdir}/usr/share/bash-completion/completions/\${pkgname}"
    "\${pkgdir}/usr/bin/\${pkgname}" completion zsh > "\${pkgdir}/usr/share/zsh/site-functions/_\${pkgname}"
    "\${pkgdir}/usr/bin/\${pkgname}" completion fish > "\${pkgdir}/usr/share/fish/vendor_completions.d/\${pkgname}.fish"
}
EOF
