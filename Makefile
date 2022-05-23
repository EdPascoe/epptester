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
ifndef VERSION
        VERSION := $(shell git describe --always --long --dirty | sed -e 's/^[A-Za-z]*_//' -e 's/-.*$$//')
endif

GITTAG=$(shell git describe --always --long --dirty )
LDFLAGS=-ldflags="-X 'main.Version=$(VERSION)' -X 'main.GitTag=$(GITTAG)'"

build:
	go build -v ${LDFLAGS} -o ./bin/$(binary) ./cmd/epptester

.PHONY: prod
prod:
	rm -rf bin
	GOOS=windows GOARCH=amd64 go build -v ${LDFLAGS} -o ./bin/$(binary)_windows_amd64.exe ./cmd/epptester
	GOOS=linux GOARCH=amd64 go build -v ${LDFLAGS} -o ./bin/$(binary)_linux_amd64 ./cmd/epptester
	GOOS=darwin GOARCH=amd64 go build -v ${LDFLAGS} -o ./bin/$(binary)_darwin_amd64 ./cmd/epptester

.PHONY: changelog
changelog:
	gitchangelog > CHANGELOG.md

clean:
	rm -rf bin

release: prod CHANGELOG.md
	@ # Tags and pushes the version but only if tagged. Also need a CHANGELOG.md file. (see make changelog) and must be master branch
ifeq "$(shell git branch --show-current | grep master | wc -l)" "0"
	@echo "Master branch only!!"
	false
endif
ifneq ("$(VERSION)", "$(shell git tags | grep "^$(VERSION)$$")")
	@echo "Please tag the repo with this version number before pushing"
	@echo "Eg: git tag -f -a $(VERSION) -m 'Release $(MSG)'"
	@echo "git push --follow-tags"
	@echo " "
	@false
endif
	false
	@# See: https://github.com/cli/cli/releases for the github gh  tool
	gh release create v$(VERSON) bin/* --target $(shell git branch --show-current ) -d -F CHANGELOG.md
