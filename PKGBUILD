
# Maintainer: Quentin Bouvet <qbouvet@outlook.com>
pkgname="fanctl-git"
pkgver=0.1
pkgrel=3
pkgdesc="Fan control"
arch=('any')
url="https://github.com/qbouvet/fanctl"
license=('GPL')
depends=('s-tui')
makedepends=('go') 
optdepends=('')
provides=("${pkgname%-git}")
conflicts=("${pkgname%-git}")
source=("$pkgname-v$pkgver::https://github.com/qbouvet/fanctl/archive/v${pkgver}r${pkgrel}.tar.gz")
md5sums=('SKIP')


build() {
  cd "$srcdir/${pkgname%-git}-${pkgver}r${pkgrel}"
  make
}

package() {
  cd "$srcdir/${pkgname%-git}-${pkgver}r${pkgrel}"
  make DESTDIR="$pkgdir/" install
}
