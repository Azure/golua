#!/bin/sh

set -ex

go install github.com/Azure/golua/cmd/glua
cd lua-5.3.4
glua -tests all.lua | grep -q OK