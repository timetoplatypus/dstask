#!/bin/sh
set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <tag>"
    echo "Release version required as argument"
    exit 1
fi

VERSION="$1"
GIT_COMMIT=$(git rev-list -1 HEAD)
RELEASE_FILE=RELEASE.md

LDFLAGS="-s -w -X  main.GIT_COMMIT=$GIT_COMMIT -X  main.VERSION=$VERSION"

# get release information

if ! test -f $RELEASE_FILE || head -n 1 $RELEASE_FILE | grep -vq $VERSION; then
    # file doesn't exist or is for old version, replace
    printf "$VERSION\n\n\n" > $RELEASE_FILE
fi

vim "+ normal G $" $RELEASE_FILE


# build

mkdir -p dist

GOOS=linux GOARCH=arm GOARM=5 go build -mod=vendor -ldflags="$LDFLAGS" cmd/dstask.go
upx -q dstask
mv dstask dist/dstask-linux-arm5

GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="$LDFLAGS" cmd/dstask.go
upx -q dstask
mv dstask dist/dstask-linux-amd64

GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags="$LDFLAGS" cmd/dstask.go
# see https://github.com/upx/upx/issues/222 -- UPX produces broken darwin executables.
#upx -q dstask
mv dstask dist/dstask-darwin-amd64

hub release create \
    -a dist/dstask-linux-arm5#"dstask linux-arm5" \
    -a dist/dstask-linux-amd64#"dstask linux-amd64" \
    -a dist/dstask-darwin-amd64#"dstask darwin-amd64" \
    -F $RELEASE_FILE \
    $1

rm -rf dist/tmp
