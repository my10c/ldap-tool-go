module main

go 1.20

require (
	my10c.ldap/cmds/common v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/create/group v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/create/menu v0.0.0-00010101000000-000000000000
	my10c.ldap/cmds/create/sudo v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/create/user v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/delete/group v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/delete/menu v0.0.0-00010101000000-000000000000
	my10c.ldap/cmds/delete/sudo v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/delete/user v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/limit v0.0.0-00010101000000-000000000000
	my10c.ldap/cmds/modify/group v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/modify/menu v0.0.0-00010101000000-000000000000
	my10c.ldap/cmds/modify/sudo v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/modify/user v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/search/group v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/search/menu v0.0.0-00010101000000-000000000000
	my10c.ldap/cmds/search/sudo v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/cmds/search/user v0.0.0-00010101000000-000000000000 // indirect
	my10c.ldap/configurator v0.0.0-00010101000000-000000000000
	my10c.ldap/initializer v0.0.0-00010101000000-000000000000
	my10c.ldap/ldap v0.0.0-00010101000000-000000000000
	my10c.ldap/logs v0.0.0-00010101000000-000000000000
	my10c.ldap/vars v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/akamensky/argparse v1.4.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect

	github.com/my10c/packages-go/epoch v0.0.0-20230629045125-5efc1e3334c4 // indirect
	github.com/my10c/packages-go/exit v0.0.0-20230629045125-5efc1e3334c4 // indirect
	github.com/my10c/packages-go/is v0.0.0-20230629045125-5efc1e3334c4 // indirect
	github.com/my10c/packages-go/lock v0.0.0-20230629045125-5efc1e3334c4
	github.com/my10c/packages-go/print v0.0.0-20230629045125-5efc1e3334c4 // indirect
	github.com/my10c/packages-go/random v0.0.0-20230629045125-5efc1e3334c4 // indirect
	github.com/my10c/packages-go/readinput v0.0.0-20230629045125-5efc1e3334c4 // indirect
	github.com/my10c/packages-go/spinner v0.0.0-20230629045125-5efc1e3334c4

	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace my10c.ldap/configurator => ./mod/configurator

replace my10c.ldap/initializer => ./mod/initializer

replace my10c.ldap/ldap => ./mod/ldap

replace my10c.ldap/logs => ./mod/logs

replace my10c.ldap/vars => ./mod/vars

replace my10c.ldap/cmds/create/menu => ./mod/cmds/create/menu

replace my10c.ldap/cmds/create/group => ./mod/cmds/create/group

replace my10c.ldap/cmds/create/sudo => ./mod/cmds/create/sudo

replace my10c.ldap/cmds/create/user => ./mod/cmds/create/user

replace my10c.ldap/cmds/delete/menu => ./mod/cmds/delete/menu

replace my10c.ldap/cmds/delete/group => ./mod/cmds/delete/group

replace my10c.ldap/cmds/delete/sudo => ./mod/cmds/delete/sudo

replace my10c.ldap/cmds/delete/user => ./mod/cmds/delete/user

replace my10c.ldap/cmds/modify/menu => ./mod/cmds/modify/menu

replace my10c.ldap/cmds/modify/group => ./mod/cmds/modify/group

replace my10c.ldap/cmds/modify/sudo => ./mod/cmds/modify/sudo

replace my10c.ldap/cmds/modify/user => ./mod/cmds/modify/user

replace my10c.ldap/cmds/search/menu => ./mod/cmds/search/menu

replace my10c.ldap/cmds/search/group => ./mod/cmds/search/group

replace my10c.ldap/cmds/search/sudo => ./mod/cmds/search/sudo

replace my10c.ldap/cmds/search/user => ./mod/cmds/search/user

replace my10c.ldap/cmds/limit => ./mod/cmds/limit

replace my10c.ldap/cmds/common => ./mod/cmds/common
