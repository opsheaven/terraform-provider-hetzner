COMMIT_TYPE   		:=$(subst #,,$(shell git log -1 --pretty=%B | grep -E -i -o "#(minor|major|patch)"))
CURRENT_TAG         :=$(or $(shell git tag | sort -r --version-sort | head -n1 ), v0.0.0)
VERSION 			= $(subst v,,$(CURRENT_TAG))
VERSION_PARTS      	= $(subst ., ,$(VERSION))
MAJOR              	= $(word 1,$(VERSION_PARTS))
MINOR              	= $(word 2,$(VERSION_PARTS))
PATCH              	= $(word 3,$(VERSION_PARTS))

ifeq ($(COMMIT_TYPE),major)
NEXT_VERSION= v$(shell echo $$(($(MAJOR)+1))).0.0
else ifeq ($(COMMIT_TYPE),minor)
NEXT_VERSION= v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0
else ifeq ($(COMMIT_TYPE),patch)
NEXT_VERSION= v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))
endif

.PHONY: tag
tag:
ifneq ($(NEXT_VERSION),)
	@echo "Creating tag $(NEXT_VERSION)"
	@git tag $(NEXT_VERSION)
	@git push origin $(NEXT_VERSION)
	@gh release create $(NEXT_VERSION) --generate-notes
else
	@echo "Last commit does not contain #major,#minor or #patch. Skipping tag and release!"
endif

.PHONY: release
release: tag
ifneq ($(NEXT_VERSION),)
	echo -e "$(GPG_PRIVATE_KEY)" | gpg --import
	go install github.com/goreleaser/goreleaser@latest
	goreleaser --release --clean
endif
