#!/bin/sh

set -exv

test "$(git name-rev --name-only --tags HEAD)" == "undefined" && (echo "This is an untagged commit so no release can be made"; exit 0)

GORELEASER_URL="https://github.com/goreleaser/goreleaser/releases"

test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

last_version() {
  curl -sL -o /dev/null -w %{url_effective} "$GORELEASER_URL/latest" | 
    rev | 
    cut -f1 -d'/'| 
    rev
}

test -z "$VERSION" && VERSION="$(last_version)"
test -z "$VERSION" && {
echo "Unable to get goreleaser version." >&2
exit 1
}
rm -f "/tmp/goreleaser.tgz"
curl -L -o "/tmp/goreleaser.tgz" "$GORELEASER_URL/download/$VERSION/goreleaser_$(uname -s)_$(uname -m).tar.gz"

tar -xf "/tmp/goreleaser.tgz" -C "$TMPDIR"
GITHUB_TOKEN=$APP_SRE_BOT_PUSH_TOKEN "${TMPDIR}/goreleaser" "$@"