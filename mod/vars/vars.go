//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package vars

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/my10c/packages-go/epoch"
	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"
	"github.com/my10c/packages-go/random"
	ldapv3 "gopkg.in/ldap.v2"
)

var (
	Off    = "\x1b[0m"    // Text Reset
	Black  = "\x1b[1;30m" // Black
	Red    = "\x1b[1;31m" // Red
	Green  = "\x1b[1;32m" // Green
	Yellow = "\x1b[1;33m" // Yellow
	Blue   = "\x1b[1;34m" // Blue
	Purple = "\x1b[1;35m" // Purple
	Cyan   = "\x1b[1;36m" // Cyan
	White  = "\x1b[1;37m" // White

	RedUnderline = "\x1b[4;31m" // Red underline
	OneLineUP    = "\x1b[A"
	DangerZone   = fmt.Sprintf("%sDanger Zone%s, be sure you understand the implication!",
		RedUnderline, Off)
	ReplicaAlert = fmt.Sprintf("The server is a %sreplica%s!, access is set to read-only",
		RedUnderline, Off)
)

type Funcs struct {
	E *epoch.Epoch
	I *is.Is
	P *print.Print
	R *random.Random
}

// for ldap search
type SearchInfo struct {
	SearchBase      string
	SearchAttribute []string
}

// for input template
type Record struct {
	Value    string // default value from the configuration
	Prompt   string
	NoEmpty  bool
	UseValue bool
}

// ldap record with all possible use fields
type LdapRecord struct {
	ID              string
	DN              string
	CN              string
	GroupType       string
	MemberType      string
	Fields          map[string]string
	Group           map[string]string
	Groups          []string
	GroupAddList    []string
	GroupDelList    []string
	UserGroups      []string
	AvailableGroups []string
	SudoAddList     map[string][]string
	SudoDelList     map[string][]string
}

// for logging
type Log struct {
	LogsDir       string
	LogFile       string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
}

// for output of a ldap search
type SearchResult struct {
	RecordCount        int
	SearchResult       *ldapv3.SearchResult
	WildCardSearchBase string
	RecordSearchbase   string
	DisplayFieldID     string
}

var (
	MyVersion   = "0.2.1c"
	now         = time.Now()
	MyProgname  = path.Base(os.Args[0])
	myAuthor    = "Luc Suryo"
	myCopyright = "Copyright 2019 - " + strconv.Itoa(now.Year()) + " ©Badassops LLC"
	myLicense   = "License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥"
	myEmail     = "<luc@badassops.com>"
	MyInfo      = fmt.Sprintf("%s (version %s)\n%s\n%s\nWritten by %s %s\n",
		MyProgname, MyVersion, myCopyright, myLicense, myAuthor, myEmail)
	MyDescription = "Simple script to manage LDAP users, groups and SUDO rules"

	// use in menu
	QuitTool bool

	// ldap logs
	Logs Log

	// we sets these under variable
	// default values
	LogsDir       = "/var/log/ldap-go"
	LogFile       = fmt.Sprintf("%s.log", MyProgname)
	LogMaxSize    = 128 // megabytes
	LogMaxBackups = 14  // 14 files
	LogMaxAge     = 14  // 14 days

	// working variable and record
	DisplayUserFields []string
	UserFields        []string
	GroupFields       []string
	SudoFields        []string
	UserObjectClass   []string
	GroupObjectClass  []string
	SudoObjectClass   []string
	WorkRecord        LdapRecord
	Template          map[string]Record

	SearchResultData SearchResult

	// use by the common function
	ObjectID      string
	ProtectedList []string

	// variables use for ldap search
	UserSearchBase     = "(objectClass=inetOrgPerson)"
	GroupSearchBase    = "(|(objectClass=posixGroup)(objectClass=groupOfNames))"
	SudoRuleSearchBase = "(objectClass=sudoRole)"

	UserWildCardSearchBase = "(&(objectClass=inetOrgPerson)(uid=VALUE))"
	UserDisplayFieldID     = "uid"

	GroupWildCardSearchBase = "(&(|(objectClass=posixGroup)(objectClass=groupOfNames))(cn=VALUE))"
	GroupDisplayFieldID     = "cn"

	SudoWildCardSearchBase = "(&(objectClass=sudoRole)(cn=VALUE))"
	SudoDisplayFieldID     = "cn"

	// query to check whatever the server is a replica
	ReplicatorSearchBase   = "(&(objectClass=simpleSecurityObject)(objectClass=organizationalRole)(cn=VALUE))"
	ReplicatorDisplayField = "cn"
)
