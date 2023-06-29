
OS_ID = generic
MACHINE = generic

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
	OS_ID = Linux_$(UNAME_M)
endif
ifeq ($(UNAME_S),Darwin)
	OS_ID = Darwin_$(UNAME_M)
endif

BUILT_SOURCES = ldap-tool.go
TOOL_VERSION := $(shell cat mod/vars/vars.go | grep MyVersion | egrep -v MyProgname | awk '{print $$3}')

all:	release/ldap-tool_$(OS_ID) \
		release/ldap-tool_$(OS_ID).tar.gz \
		release/ldap-tool_$(OS_ID).sha256

release/ldap-tool_$(OS_ID): $(BUILT_SOURCES)
	@echo "build the ldap-tool_$(OS_ID) binary..."
	@go build -o release/ldap-tool_$(OS_ID) ldap-tool.go

release/ldap-tool_$(OS_ID).tar.gz: release/ldap-tool_$(OS_ID)
	@echo "create the ldap-tool_$(OS_ID).tar.gz archive..."
	@(cd release ; tar zcf ldap-tool_$(OS_ID).tar.gz ldap-tool_$(OS_ID))

release/ldap-tool_$(OS_ID).sha256: release/ldap-tool_$(OS_ID).tar.gz
	@echo "create the sha256 information file..."
	@sha256sum release/ldap-tool_$(OS_ID).tar.gz | awk '{print $$1}' > release/ldap-tool_$(OS_ID).sha256
	@echo "SHA256: $$(cat release/ldap-tool_$(OS_ID).sha256)"

install: release/ldap-tool_$(OS_ID)
	@echo "Installing the new ldap-tool binary..."
	@sudo cp release/ldap-tool_$(OS_ID) /usr/local/sbin/ldap-tool
	@sudo chmod 0755 /usr/local/sbin/ldap-tool
	@sudo chown 0:0 /usr/local/sbin/ldap-tool

clean:
	@rm -f release/*

changelog:
	@echo "version built $(TOOL_VERSION)"
