
# Maintainer: Quentin Bouvet <qbouvet@outlook.com>
pkgname="fanctl-git"
pkgver=0.1
pkgrel=1
pkgdesc="Fan control"
arch=('any')
url="https://github.com/qbouvet/fanctl"
license=('GPL')
groups=()
depends=()
makedepends=('go') 
provides=("${pkgname}")
conflicts=("${pkgname}")
replaces=()
backup=()
options=()
install=
source=("$pkgname-v$pkgver::https://github.com/qbouvet/fanctl/archive/v$pkgver.tar.gz")
noextract=()
md5sums=('SKIP')

# makepkg defines two variables that you should use as part of the build 
# and install process
#
# srcdir
#   This points to the directory where makepkg extracts or symlinks all 
#   files in the source array
#
# pkgdir
#   This points to the directory where makepkg bundles the installed 
#   package, which becomes the root directory of your built package.


#   The primary goal is to generate version numbers that will increase according to
#   pacman's version comparisons with later commits to the repo. The format
#   VERSION='VER_NUM.rREV_NUM.HASH', or a relevant subset in case VER_NUM or HASH
#   are not available, is recommended.
#
pkgver() { 
  cd "$srcdir/${pkgname}"
  # Git, tags available
  printf "%s" "$(git describe --long | sed 's/\([^-]*-\)g/r\1/;s/-/./g')"
  # Git, no tags available
  #printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

#
#   Runs commands that are used to prepare sources for building, such as patching. 
#   This function runs right after package extraction, before pkgver() and the 
#   build function. If extraction is skipped (makepkg --noextract), then prepare() 
#   is not run.
#
#prepare() {
#}


#   uses common shell commands in Bash syntax to automatically compile 
#   software and create a directory called pkg to install the software 
#   to.
#   The build() function in essence automates everything you did by hand 
#   and compiles the software in the fakeroot build environment.
build() {
  cd "$srcdir/${pkgname}"
  #./autogen.sh
  #./configure --prefix=/usr
  make
}

#   Place for calls to make check and similar testing routines
#
#check() {
#  cd "$srcdir/${pkgname%}"
#  make -k check
#}

#   Put the compiled files in a directory where makepkg can retrieve 
#   them to create a package. This by default is the pkg directory—a 
#   simple fakeroot environment. The pkg directory replicates the 
#   hierarchy of the root file system of the software's installation 
#   paths. If you have to manually place files under the root of your 
#   filesystem, you should install them in the pkg directory under the 
#   same directory structure.
package() {
  cd "$srcdir/${pkgname}"
  make DESTDIR="$pkgdir/" install
}
