
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

SOURCES = ldap-tool.go \
	mod/vars/vars.go \
	mod/initializer/initializer.go \
	mod/ldap/create.go \
	mod/ldap/ldap.go \
	mod/ldap/delete.go \
	mod/ldap/search.go \
	mod/ldap/modify.go \
	mod/ldap/get.go \
	mod/ldap/password.go \
	mod/cmds/delete/group/group.go \
	mod/cmds/delete/user/user.go \
	mod/cmds/delete/menu/menu.go \
	mod/cmds/delete/sudo/sudo.go \
	mod/cmds/limit/search.go \
	mod/cmds/limit/modify.go \
	mod/cmds/modify/group/group.go \
	mod/cmds/modify/user/user.go \
	mod/cmds/modify/menu/menu.go \
	mod/cmds/modify/sudo/sudo.go \
	mod/cmds/search/group/groups.go \
	mod/cmds/search/group/group.go \
	mod/cmds/search/user/users.go \
	mod/cmds/search/user/user.go \
	mod/cmds/search/menu/menu.go \
	mod/cmds/search/sudo/sudo.go \
	mod/cmds/search/sudo/sudos.go \
	mod/cmds/common/delete.go \
	mod/cmds/common/search.go \
	mod/cmds/create/group/group.go \
	mod/cmds/create/user/user.go \
	mod/cmds/create/menu/menu.go \
	mod/cmds/create/sudo/sudo.go \
	mod/logs/logs.go \
	mod/configurator/configurator.go

BUILT_SOURCES = $(SOURCES)
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
	@rm -f release/*$(OS_ID)*

changelog:
	@echo "version built $(TOOL_VERSION)"
