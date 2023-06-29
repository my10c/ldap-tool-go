//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package user

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	r "github.com/my10c/packages-go/readinput"
	ldapv3 "gopkg.in/ldap.v2"
)

func printUsers(records *ldapv3.SearchResult, recordCount int, funcs *vars.Funcs) {
	baseInfo := false
	fmt.Printf("\tPrint full name and department (default to N)? [y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	valueEntered, _ := reader.ReadString('\n')
	valueEntered = strings.TrimSuffix(valueEntered, "\n")
	if r.ReadYN(valueEntered, false) == true {
		baseInfo = true
	}

	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 55))
	for idx, entry := range records.Entries {
		funcs.P.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		if baseInfo {
			userBaseInfo := fmt.Sprintf("\tFull namae: %s %s\t\tdepartmentNumber %s\n\n",
				records.Entries[idx].GetAttributeValue("givenName"),
				records.Entries[idx].GetAttributeValue("sn"),
				records.Entries[idx].GetAttributeValue("departmentNumber"))
			funcs.P.PrintCyan(userBaseInfo)
		}
	}
	funcs.P.PrintYellow(fmt.Sprintf("\n\tTotal records: %d \n", recordCount))
}

func Users(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search Users", 20, true))
	conn.SearchInfo.SearchBase = vars.UserSearchBase
	conn.SearchInfo.SearchAttribute = []string{}
	records, recordsCount := conn.Search()
	printUsers(records, recordsCount, funcs)
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
