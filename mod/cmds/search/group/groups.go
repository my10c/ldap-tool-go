//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package group

import (
	"fmt"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printGroups(records *ldapv3.SearchResult, recordCount int, funcs *vars.Funcs) {
	var memberCount = 0
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 55))
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

func Groups(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search Groups", 20, true))
	conn.SearchInfo.SearchBase = vars.GroupSearchBase
	conn.SearchInfo.SearchAttribute = []string{}
	records, recordCount := conn.Search()
	printGroups(records, recordCount, funcs)
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
