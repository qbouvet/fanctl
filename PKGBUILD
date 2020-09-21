
# Maintainer: Quentin Bouvet <qbouvet@outlook.com>
pkgname="fanctl-git"
pkgver=0.1
pkgrel=4
pkgdesc="Fan control"
arch=('any')
url="https://github.com/qbouvet/fanctl"
license=('GPL')
depends=()
makedepends=('go') 
optdepends=('nvidia-settings: nvidia support'
            'nvidia-utils:    nvidia support'
            's-tui:           CPU power consumption')
provides=("${pkgname%-git}")
conflicts=("${pkgname%-git}")
source=("$pkgname-v$pkgver.r$pkgrel::https://github.com/qbouvet/fanctl/archive/v${pkgver}.r${pkgrel}.tar.gz")
md5sums=('SKIP')


build() {
  cd "$srcdir/${pkgname%-git}-${pkgver}.r${pkgrel}"
  make
}

package() {
  cd "$srcdir/${pkgname%-git}-${pkgver}.r${pkgrel}"
  make DESTDIR="$pkgdir/" install
}
