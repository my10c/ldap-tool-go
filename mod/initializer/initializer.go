//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package initializer

import (
	"fmt"
	"strconv"

	"my10c.ldap/configurator"
	"my10c.ldap/vars"
	"github.com/my10c/packages-go/epoch"
	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"
	"github.com/my10c/packages-go/random"
)

var (
	msg string
)

// initialize the system/variable/template
func Init(config *configurator.Config) *vars.Funcs {
	vars.QuitTool = false

	// ldap fields that will be used
	vars.UserFields = []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowLastChange", "shadowExpire", "shadowWarning", "shadowMax",
		"sshPublicKey"}

	vars.DisplayUserFields = []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowWarning", "shadowMax", "sshPublicKey"}

	vars.GroupFields = []string{"cn", "groupName", "groupType", "gidNumber", "memberUid", "member"}

	vars.SudoFields = []string{"cn", "sudoCommand", "sudoHost", "sudoOption",
		"sudoOrder", "sudoRunAsUser"}

	vars.UserObjectClass = []string{"top", "person",
		"organizationalPerson", "inetOrgPerson",
		"posixAccount", "shadowAccount", "ldapPublicKey"}

	vars.GroupObjectClass = []string{"posix", "groupOfNames"}

	vars.SudoObjectClass = []string{"top", "sudoRole"}

	// set to default if not set in the configuration file
	if config.LogValues.LogsDir == "" {
		vars.Logs.LogsDir = vars.LogsDir
	}
	if config.LogValues.LogFile == "" {
		vars.Logs.LogFile = vars.LogFile
	}
	if config.LogValues.LogMaxSize == 0 {
		vars.Logs.LogMaxSize = vars.LogMaxSize
	}
	if config.LogValues.LogMaxBackups == 0 {
		vars.Logs.LogMaxBackups = vars.LogMaxBackups
	}
	if config.LogValues.LogMaxAge == 0 {
		vars.Logs.LogMaxAge = vars.LogMaxAge
	}

	// initialize the maps
	vars.WorkRecord.Fields = make(map[string]string)
	vars.WorkRecord.Group = make(map[string]string)
	vars.WorkRecord.SudoAddList = make(map[string][]string)
	vars.WorkRecord.SudoDelList = make(map[string][]string)
	vars.Template = make(map[string]vars.Record)

	// set to expire by default as today + ShadowMax
	epochPtr := epoch.New()
	printPtr := print.New()
	isPtr := is.New()
	currExpired := strconv.FormatInt(epochPtr.Days()+int64(config.DefaultValues.ShadowMax), 10)

	// user
	vars.Template["uid"] =
		vars.Record{
			Prompt:   "Enter userid (login name) to be use",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	vars.Template["givenName"] =
		vars.Record{
			Prompt:   "Enter First name",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	vars.Template["sn"] =
		vars.Record{
			Prompt:   "Enter Last name",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	vars.Template["mail"] =
		vars.Record{
			Prompt:   "Enter email",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["uidNumber"] =
		vars.Record{
			Prompt:   "Enter user's UID",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["departmentNumber"] =
		vars.Record{
			Prompt:   "Enter department",
			Value:    config.DefaultValues.GroupName,
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["loginShell"] =
		vars.Record{
			Prompt:   "Enter shell",
			Value:    config.DefaultValues.Shell,
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["userPassword"] =
		vars.Record{
			Prompt:   "Enter password",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["shadowMax"] =
		vars.Record{
			Prompt:   "Enter the max password age",
			Value:    strconv.Itoa(config.DefaultValues.ShadowAge),
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["shadowWarning"] =
		vars.Record{
			Prompt:   "Enter the days notification before the password expires",
			Value:    strconv.Itoa(config.DefaultValues.ShadowWarning),
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["sshPublicKey"] =
		vars.Record{
			Prompt:   "Enter SSH the Public Key",
			Value:    "is missing",
			NoEmpty:  false,
			UseValue: false,
		}

	vars.Template["shadowExpire"] =
		vars.Record{
			Prompt:   fmt.Sprintf("Reset password expired, Y/N"),
			Value:    currExpired,
			NoEmpty:  false,
			UseValue: false,
		}

	vars.Template["shadowLastChange"] =
		vars.Record{
			Prompt:   "Date the password was last changed",
			Value:    strconv.FormatInt(epochPtr.Days(), 10),
			NoEmpty:  false,
			UseValue: false,
		}

	// share in group and sudo rule
	vars.Template["cn"] =
		vars.Record{
			Prompt:   "Short name, will be use to filled the dn value",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	// group
	vars.Template["groupName"] =
		vars.Record{
			Prompt:   "Enter the group name",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	vars.Template["groupType"] =
		vars.Record{
			Prompt:   "Group type (p)osix or (g)roupOfNames (default to posix)",
			Value:    "posix",
			NoEmpty:  false,
			UseValue: true,
		}

	// only use for posix group
	vars.Template["gidNumber"] =
		vars.Record{
			Prompt:   "Group ID/number of the posix group",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	// these are automatically filled
	vars.Template["objectClass"] =
		vars.Record{
			Prompt:   "Auto filled based on group type, posix or groupOfNames (default to posix)",
			Value:    "",
			NoEmpty:  true,
			UseValue: true,
		}

	// the default is always used
	vars.Template["member"] =
		vars.Record{
			Prompt:   "Auto filled based on the groupDN value",
			Value:    fmt.Sprintf("uid=initial-member,%s", config.ServerValues.UserDN),
			NoEmpty:  true,
			UseValue: false,
		}

	vars.Template["sudoCommand"] =
		vars.Record{
			Prompt:   "fully qualified path of the commands allow with this rule",
			Value:    "ALL",
			NoEmpty:  false,
			UseValue: false,
		}

	msg = printPtr.MessageYellow("default to ALL")
	vars.Template["sudoHost"] =
		vars.Record{
			Prompt:   "The host the command is allowed",
			Value:    "ALL",
			NoEmpty:  false,
			UseValue: true,
		}

	msg = fmt.Sprintf("%sExmple%s %s!authenticate%s", vars.Purple, vars.Off, vars.Cyan, vars.Off)
	msg = msg + printPtr.MessageYellow(" for no password required")
	vars.Template["sudoOption"] =
		vars.Record{
			Prompt:   fmt.Sprintf("%s\n\tSudo option with the command", msg),
			Value:    "",
			NoEmpty:  false,
			UseValue: false,
		}

	vars.Template["sudoOrder"] =
		vars.Record{
			Prompt:   "The order of the rule (between 3 and 10)",
			Value:    "4",
			NoEmpty:  false,
			UseValue: true,
		}

	vars.Template["sudoRunAsUser"] =
		vars.Record{
			Prompt:   "Run the command as the user",
			Value:    "root",
			NoEmpty:  false,
			UseValue: true,
		}

	return &vars.Funcs{
		E: epochPtr,
		I: isPtr,
		P: printPtr,
		R: random.New(config.DefaultValues.PassComplex, config.DefaultValues.PassLenght),
	}
}
