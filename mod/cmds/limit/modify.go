//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package limit

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func modifyUserPassword(conn *ldap.Connection, records *ldapv3.SearchResult, funcs *vars.Funcs) {
	fieldName := "userPassword"
	reader := bufio.NewReader(os.Stdin)

	userName := records.Entries[0].GetAttributeValue("uid")
	vars.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", userName, conn.Config.ServerValues.UserDN)
	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	funcs.P.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", userName))
	funcs.P.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))

	passWord := funcs.R.Generate()
	funcs.P.PrintCyan(fmt.Sprintf("\n\tCurrent value (encrypted!): %s%s%s\n",
		vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))
	funcs.P.PrintYellow(fmt.Sprintf("\t\tsuggested password: %s\n", passWord))

	funcs.P.PrintPurple(fmt.Sprintf("\t%s: ", vars.Template[fieldName].Prompt))
	valueEntered, _ := reader.ReadString('\n')
	valueEntered = strings.TrimSuffix(valueEntered, "\n")

	if len(valueEntered) != 0 {
		vars.WorkRecord.Fields[fieldName] = valueEntered
		vars.WorkRecord.Fields["shadowLastChange"] = vars.Template["shadowLastChange"].Value
	} else {
		fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
		funcs.P.PrintBlue(fmt.Sprintf("\n\tNo password was entered, the user %s's password was not changed\n",
			userName))
		return
	}

	if !conn.ModifyUser() {
		funcs.P.PrintRed(fmt.Sprintf("\n\tFailed modify the %s's password, check the log file\n", userName))
		funcs.P.PrintYellow("Make sure you are not a replica server!\n")
		return
	}
	funcs.P.PrintGreen(fmt.Sprintf("\n\tUser %s modified successfully\n", userName))
}

func ModifyUserPassword(conn *ldap.Connection, funcs *vars.Funcs) {
	reg, _ := regexp.Compile("^uid=|,ou=users,.*")
	userID := reg.ReplaceAllString(conn.Config.ServerValues.Admin, "")
	conn.SearchInfo.SearchBase = strings.ReplaceAll(vars.UserWildCardSearchBase, "VALUE", userID)
	conn.SearchInfo.SearchAttribute = []string{}
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Modify Own Password", 14, true))
	record, recordCount := conn.Search()
	if recordCount != 1 {
		funcs.P.PrintRed(fmt.Sprintf("\n\tUser %s was not found, aborting...\n", userID))
		return
	} else {
		modifyUserPassword(conn, record, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
