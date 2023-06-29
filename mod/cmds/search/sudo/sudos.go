//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package sudo

import (
	"fmt"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printSudos(records []*ldapv3.Entry, recordCount int, protectedSudoRules []string, funcs *vars.Funcs) {
	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	for _, entry := range records {
		funcs.P.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		for _, attributes := range entry.Attributes {
			for _, value := range attributes.Values {
				if attributes.Name != "objectClass" {
					if attributes.Name == "cn" {
						if funcs.I.IsInList(protectedSudoRules, value) {
							funcs.P.PrintYellow("\t\tThis entry can not be modified\n")
						}
					}
					funcs.P.PrintCyan(fmt.Sprintf("\t%s : %s \n", attributes.Name, value))
				}
			}
		}
		fmt.Printf("\n")
	}
}

func Sudos(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search Sudo Rules", 15, true))
	conn.SearchInfo.SearchBase = vars.SudoRuleSearchBase
	conn.SearchInfo.SearchAttribute = []string{}
	records, recordCount := conn.Search()
	printSudos(records.Entries, recordCount, conn.Config.SudoValues.ExcludeSudo, funcs)
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
