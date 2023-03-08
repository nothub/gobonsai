#!/usr/bin/env sh

set -o errexit
set -o nounset

version="${1}"
arch="${2}"

bin="out/gobonsai_${version}_linux-${arch}"
deb="out/gobonsai_${version}_${arch}.deb"
target="/usr/local/games/gobonsai"

if ! test -r "${bin}"; then
    echo >&2 "can not read ${bin}"
    exit 1
fi

work=$(mktemp -d)
echo >&2 "assembling ${arch} deb in ${work}"

mkdir -p "${work}/DEBIAN"
cat >"${work}/DEBIAN/control" <<EOF
Package: gobonsai
Version: ${version}
Architecture: ${arch}
Depends: libc6 (>= 2.24)
Priority: optional
Section: games
Description: A bonsai tree generator
Homepage: https://github.com/nothub/gobonsai
Maintainer: Florian Hübner <code@hub.lol>
EOF

mkdir -p "$(dirname "${work}${target}")"
cp "${bin}" "${work}${target}"

# generate checksums
(cd "${work}" && find . -type f -not -path "./DEBIAN/*" -exec md5sum {} \; >"${work}/DEBIAN/md5sums")

# bundle deb package
dpkg-deb --root-owner-group --verbose --build "${work}" "${deb}"
