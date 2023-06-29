//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package group

import (
	"fmt"

	"my10c.ldap/cmds/common"
	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printGroup(records *ldapv3.SearchResult, funcs *vars.Funcs) {
	var memberCount = 0
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	for idx, entry := range records.Entries {
		funcs.P.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		funcs.P.PrintBlue(fmt.Sprintf("\tcn: %s\n",
			records.Entries[idx].GetAttributeValue("cn")))
		if len(records.Entries[idx].GetAttributeValue("gidNumber")) != 0 {
			funcs.P.PrintCyan(fmt.Sprintf("\tgidNumber: %s\n",
				records.Entries[idx].GetAttributeValue("gidNumber")))
			memberCount = 0
			for _, member := range entry.GetAttributeValues("memberUid") {
				funcs.P.PrintCyan(fmt.Sprintf("\tmemberUid: %s\n", member))
				memberCount++
			}
			funcs.P.PrintYellow(fmt.Sprintf("\tTotal members: %d : posix group\n\n", memberCount))
		} else {
			memberCount = 0
			for _, member := range entry.GetAttributeValues("member") {
				funcs.P.PrintCyan(fmt.Sprintf("\tmember: %s\n", member))
				memberCount++
			}
			funcs.P.PrintYellow(fmt.Sprintf("\tTotal members: %d : groupOfNames group\n\n", memberCount))
		}
	}
}

func Group(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search Group", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.GroupWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.GroupWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.GroupDisplayFieldID
	if common.GetObjectRecord(conn, true, "group", funcs) {
		printGroup(vars.SearchResultData.SearchResult, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
