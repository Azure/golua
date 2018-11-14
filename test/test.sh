#!/usr/bin/env bash
set -euxo pipefail

root_dir="$(cd "${BASH_SOURCE[0]%/*}/.." && pwd)"
cd "$root_dir"

go install github.com/Azure/golua/cmd/glua

cd "test/lua-5.3.4"
glua -tests all.lua | grep -q OK
