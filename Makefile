NAME := sms
DESCRIPTION := Send SMS through RDCom API
COPYRIGHT := 2025 © Andrea Funtò
LICENSE := MIT
LICENSE_URL := https://opensource.org/license/mit/
VERSION_MAJOR := 0
VERSION_MINOR := 0
VERSION_PATCH := 1
VERSION=$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH)
MAINTAINER=dihedron.dev@gmail.com
VENDOR=dihedron.dev@gmail.com
PRODUCER_URL=https://github.com/dihedron/
DOWNLOAD_URL=$(PRODUCER_URL)snoop
METADATA_PACKAGE=$$(grep "module .*" go.mod | sed 's/module //gi')/version

_RULES_MK_MINIMUM_VERSION=202502220945
_RULES_MK_ENABLE_CGO=0
_RULES_MK_ENABLE_GOGEN=1
_RULES_MK_ENABLE_RACE=0
#_RULES_MK_STATIC_LINK=1
#_RULES_MK_ENABLE_NETGO=1
#_RULES_MK_STRIP_SYMBOLS=1
#_RULES_MK_STRIP_DBG_INFO=1
#_RULES_MK_FORCE_DEP_REBUILD=1

include rules.mk

.PHONY: clean-cache ## remove all cached build entries
clean-cache:
	@go clean -x -cache



MY_STRING := my-example-string-with-dashes

# Convert to uppercase and replace dashes with underscores
UPPER_UNDERSCORE_STRING := $(shell echo $(MY_STRING) | tr '[:lower:]' '[:upper:]' | tr '-' '_')_

# Print the result
print-stuff:
	@echo $(UPPER_UNDERSCORE_STRING)
