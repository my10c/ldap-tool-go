//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package limit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

var (
	displayFields = []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowWarning", "shadowMax", "sshPublicKey"}
)

func printUserRecord(conn *ldap.Connection, records *ldapv3.SearchResult, funcs *vars.Funcs) {
	// the values are in days so we need to multiple by 86400
	value, _ := strconv.ParseInt(records.Entries[0].GetAttributeValue("shadowLastChange"), 10, 64)
	passChanged := funcs.E.ReadableEpoch(value * 86400)

	value, _ = strconv.ParseInt(records.Entries[0].GetAttributeValue("shadowExpire"), 10, 64)
	passExpired := funcs.E.ReadableEpoch(value * 86400)

	userName := records.Entries[0].GetAttributeValue("uid")
	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	funcs.P.PrintBlue(fmt.Sprintf("\tdn: %s\n", records.Entries[0].DN))
	userDN := records.Entries[0].DN
	for _, fieldName := range vars.DisplayUserFields {
		funcs.P.PrintCyan(fmt.Sprintf("\t%s: %s\n", fieldName, records.Entries[0].GetAttributeValue(fieldName)))
	}

	conn.SearchInfo.SearchBase =
		fmt.Sprintf("(|(&(objectClass=posixGroup)(memberUid=%s))(&(objectClass=groupOfNames)(member=%s)))",
			userName, userDN)
	conn.SearchInfo.SearchAttribute = []string{"dn"}
	groupRecords, _ := conn.Search()
	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	funcs.P.PrintPurple(fmt.Sprintf("\tUser %s groups:\n", userName))
	for _, entry := range groupRecords.Entries {
		funcs.P.PrintCyan(fmt.Sprintf("\tdn: %s\n", entry.DN))
	}

	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	funcs.P.PrintPurple(fmt.Sprintf("\tUser %s password information\n", userName))
	funcs.P.PrintCyan(fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	funcs.P.PrintRed(fmt.Sprintf("\tPassword will expired on %s\n\n", passExpired))
}

func UserRecord(conn *ldap.Connection, funcs *vars.Funcs) {
	reg, _ := regexp.Compile("^uid=|,ou=users,.*")
	userID := reg.ReplaceAllString(conn.Config.ServerValues.Admin, "")
	conn.SearchInfo.SearchBase = strings.ReplaceAll(vars.UserWildCardSearchBase, "VALUE", userID)
	conn.SearchInfo.SearchAttribute = []string{}
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search User", 18, true))
	record, recordCount := conn.Search()
	if recordCount != 1 {
		funcs.P.PrintRed(fmt.Sprintf("\n\tUser %s was not found, aborting...\n", userID))
		return
	} else {
		printUserRecord(conn, record, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
