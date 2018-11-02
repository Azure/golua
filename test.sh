#!/bin/sh

set -ex

go install ./cmd/glua/...
cd ./test/lua-5.3.4
glua -tests all.lua | grep -q OK