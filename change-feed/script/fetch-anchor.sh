#!/usr/bin/env bash
set -euo pipefail
BASE_SHA="1d6267568cc63a85b07f84d1d0afa9b374fd782f"
HEAD_SHA="3e08e45a052d4b8fe77fb5838c28652af9e89ecf"
DIR="$(dirname "$0")/../testdata/anchor"; mkdir -p "$DIR"
P="descriptions/api.github.com/api.github.com.json"
gh api "repos/github/rest-api-description/contents/$P?ref=$BASE_SHA" -H 'Accept: application/vnd.github.raw' > "$DIR/base.json"
gh api "repos/github/rest-api-description/contents/$P?ref=$HEAD_SHA" -H 'Accept: application/vnd.github.raw' > "$DIR/head.json"
echo "fetched $(wc -c < "$DIR/base.json") + $(wc -c < "$DIR/head.json") bytes"
