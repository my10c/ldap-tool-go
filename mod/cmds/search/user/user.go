//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package user

import (
	"fmt"
	"strconv"

	"my10c.ldap/cmds/common"
	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printUser(conn *ldap.Connection, records *ldapv3.SearchResult, funcs *vars.Funcs) {
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

	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	funcs.P.PrintPurple(fmt.Sprintf("\tUser %s password information\n", userName))
	funcs.P.PrintCyan(fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	funcs.P.PrintRed(fmt.Sprintf("\tPassword will expired on %s\n\n", passExpired))
}

func User(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search User", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.UserWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.UserWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.UserDisplayFieldID
	if common.GetObjectRecord(conn, true, "user", funcs) {
		printUser(conn, vars.SearchResultData.SearchResult, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
