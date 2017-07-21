#!/bin/bash
#
# Copyright 2016 ThoughtWorks, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License..
# You may obtain a copy of the License at
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set +x
set -e

BUILD_DIR=$(pwd)
export GOPATH=$BUILD_DIR
export GOBIN=$GOPATH/bin
#
# Pull dependencies
#
echo "================"
echo "Get dependencies"
echo "================"
echo ""

go get -u github.com/mna/pigeon
go get -u github.com/Pallinder/go-randomdata

echo "================"
echo "Building DSL parser"
echo "================"
echo ""

pigeon -o dsl/dsl.go dsl/dsl.peg

mkdir -p src/github.com/ThoughtWorksStudios/
ln -s $(pwd) src/github.com/ThoughtWorksStudios/

echo "================"
echo "Building project"
echo "================"
echo ""

go build
