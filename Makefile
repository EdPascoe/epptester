# Copyright 2022 Ed Pascoe.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The binary to build (just the basename).
binary := epptester

# This repo's root import path (under GOPATH). From https://github.com/EdPascoe/epptester.git
PKG := github.com/EdPascoe/epptester
VERSION ?= main
TAG_LATEST ?= false

release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/$(binary)_windows_amd64.exe ./cmd/epptester
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(binary)_linux_amd64 ./cmd/epptester
	GOOS=darwin GOARCH=amd64 go build -o ./bin/$(binary)_darwin_amd64 ./cmd/epptester

local:
	go build -o ./bin/$(binary) ./cmd/epptester

clean:
	rm -rf bin
push:
	gsutil cp bin/* gs://$(binary)-release


