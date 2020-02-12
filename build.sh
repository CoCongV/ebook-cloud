set -eu


VERSION=$(git describe --abbrev=0 --tags)
BUILDTIME=$(date -u +%Y/%m/%d-%H:%M:%S)

go build -o dist/ebook-cloud-"$VERSION" main.go
