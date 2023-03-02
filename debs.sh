#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

build() {
    local workdir
    workdir=$(mktemp -d)
    echo >&2 "assembling files for gobonsai_${1}_${2}.deb in ${workdir}"
    mkdir -p "${workdir}/DEBIAN"
    cat >"${workdir}/DEBIAN/control" <<EOF
Package: gobonsai
Version: ${1}
Architecture: ${2}
Depends: libc6 (>= 2.24)
Priority: optional
Section: games
Description: A bonsai tree generator
Homepage: https://github.com/nothub/gobonsai
Maintainer: Florian Hübner <code@hub.lol>
EOF
    mkdir -p "${workdir}/usr/local/games"
    cp "out/gobonsai-linux" "${workdir}/usr/local/games/gobonsai"
    (cd "${workdir}" && find . -type f -not -path "./DEBIAN/*" -exec md5sum {} \; >"${workdir}/DEBIAN/md5sums")
    dpkg-deb --root-owner-group --verbose --build "${workdir}" "out/gobonsai_${1}_${2}.deb"
}

build "${1}" "amd64"
build "${1}" "arm64"
